package observability

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics provides Prometheus-compatible metrics for Mitran.
// Uses sync/atomic counters for lock-free performance.
var Metrics = &MetricsRegistry{
	counters:   make(map[string]*atomic.Int64),
	gauges:     make(map[string]*atomic.Int64),
	histograms: make(map[string]*Histogram),
}

// MetricsRegistry holds all registered metrics.
type MetricsRegistry struct {
	mu         sync.RWMutex
	counters   map[string]*atomic.Int64
	gauges     map[string]*atomic.Int64
	histograms map[string]*Histogram
}

// Histogram tracks value distributions with predefined buckets.
type Histogram struct {
	mu      sync.Mutex
	buckets []float64
	counts  []int64
	sum     float64
	count   int64
}

var defaultBuckets = []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0, 60.0}

func newHistogram(buckets []float64) *Histogram {
	if buckets == nil {
		buckets = defaultBuckets
	}
	sort.Float64s(buckets)
	return &Histogram{
		buckets: buckets,
		counts:  make([]int64, len(buckets)+1), // +1 for +Inf
	}
}

func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sum += value
	h.count++
	for i, b := range h.buckets {
		if value <= b {
			h.counts[i]++
		}
	}
	h.counts[len(h.buckets)]++ // +Inf always incremented
}

// --- Counter operations ---

// IncrCounter increments a counter by 1.
func IncrCounter(name string, labels map[string]string) {
	key := metricKey(name, labels)
	c := Metrics.getOrCreateCounter(key)
	c.Add(1)
}

// AddCounter adds a value to a counter.
func AddCounter(name string, labels map[string]string, value int64) {
	key := metricKey(name, labels)
	c := Metrics.getOrCreateCounter(key)
	c.Add(value)
}

// --- Gauge operations ---

// SetGauge sets a gauge value.
func SetGauge(name string, labels map[string]string, value int64) {
	key := metricKey(name, labels)
	g := Metrics.getOrCreateGauge(key)
	g.Store(value)
}

// IncrGauge increments a gauge by 1.
func IncrGauge(name string, labels map[string]string) {
	key := metricKey(name, labels)
	g := Metrics.getOrCreateGauge(key)
	g.Add(1)
}

// DecrGauge decrements a gauge by 1.
func DecrGauge(name string, labels map[string]string) {
	key := metricKey(name, labels)
	g := Metrics.getOrCreateGauge(key)
	g.Add(-1)
}

// --- Histogram operations ---

// ObserveHistogram records a value in a histogram.
func ObserveHistogram(name string, labels map[string]string, value float64) {
	key := metricKey(name, labels)
	h := Metrics.getOrCreateHistogram(key)
	h.Observe(value)
}

// --- Convenience methods for Mitran-specific metrics ---

// RecordTaskCompleted increments the task counter.
func RecordTaskCompleted(agent, status string) {
	IncrCounter("mitran_tasks_total", map[string]string{"agent": agent, "status": status})
}

// RecordTaskDuration records task execution duration.
func RecordTaskDuration(agent string, duration time.Duration) {
	ObserveHistogram("mitran_task_duration_seconds", map[string]string{"agent": agent}, duration.Seconds())
}

// SetActiveAgents sets the number of active agents.
func SetActiveAgents(count int64) {
	SetGauge("mitran_agents_active", nil, count)
}

// RecordLLMTokens records LLM token usage.
func RecordLLMTokens(provider, model string, tokens int64) {
	AddCounter("mitran_llm_tokens_total", map[string]string{"provider": provider, "model": model}, tokens)
}

// --- HTTP Handler ---

// RegisterRoutes registers the /metrics endpoint.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", HandleMetrics)
}

// HandleMetrics serves metrics in Prometheus text exposition format.
func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	var sb strings.Builder

	// Emit HELP/TYPE for known metrics
	sb.WriteString("# HELP mitran_tasks_total Total number of tasks processed\n")
	sb.WriteString("# TYPE mitran_tasks_total counter\n")
	sb.WriteString("# HELP mitran_task_duration_seconds Task execution duration in seconds\n")
	sb.WriteString("# TYPE mitran_task_duration_seconds histogram\n")
	sb.WriteString("# HELP mitran_agents_active Number of currently active agents\n")
	sb.WriteString("# TYPE mitran_agents_active gauge\n")
	sb.WriteString("# HELP mitran_llm_tokens_total Total LLM tokens consumed\n")
	sb.WriteString("# TYPE mitran_llm_tokens_total counter\n")

	Metrics.mu.RLock()
	defer Metrics.mu.RUnlock()

	// Counters
	keys := sortedKeys(Metrics.counters)
	for _, key := range keys {
		v := Metrics.counters[key]
		sb.WriteString(fmt.Sprintf("%s %d\n", key, v.Load()))
	}

	// Gauges
	keys = sortedKeys(Metrics.gauges)
	for _, key := range keys {
		v := Metrics.gauges[key]
		sb.WriteString(fmt.Sprintf("%s %d\n", key, v.Load()))
	}

	// Histograms
	histKeys := sortedHistKeys(Metrics.histograms)
	for _, key := range histKeys {
		h := Metrics.histograms[key]
		h.mu.Lock()
		for i, b := range h.buckets {
			sb.WriteString(fmt.Sprintf("%s_bucket{le=\"%s\"} %d\n", key, formatFloat(b), h.counts[i]))
		}
		sb.WriteString(fmt.Sprintf("%s_bucket{le=\"+Inf\"} %d\n", key, h.counts[len(h.buckets)]))
		sb.WriteString(fmt.Sprintf("%s_sum %s\n", key, formatFloat(h.sum)))
		sb.WriteString(fmt.Sprintf("%s_count %d\n", key, h.count))
		h.mu.Unlock()
	}

	w.Write([]byte(sb.String()))
}

// --- Internal helpers ---

func (m *MetricsRegistry) getOrCreateCounter(key string) *atomic.Int64 {
	m.mu.RLock()
	if c, ok := m.counters[key]; ok {
		m.mu.RUnlock()
		return c
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.counters[key]; ok {
		return c
	}
	c := &atomic.Int64{}
	m.counters[key] = c
	return c
}

func (m *MetricsRegistry) getOrCreateGauge(key string) *atomic.Int64 {
	m.mu.RLock()
	if g, ok := m.gauges[key]; ok {
		m.mu.RUnlock()
		return g
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if g, ok := m.gauges[key]; ok {
		return g
	}
	g := &atomic.Int64{}
	m.gauges[key] = g
	return g
}

func (m *MetricsRegistry) getOrCreateHistogram(key string) *Histogram {
	m.mu.RLock()
	if h, ok := m.histograms[key]; ok {
		m.mu.RUnlock()
		return h
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if h, ok := m.histograms[key]; ok {
		return h
	}
	h := newHistogram(nil)
	m.histograms[key] = h
	return h
}

func metricKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	pairs := make([]string, 0, len(labels))
	for k, v := range labels {
		pairs = append(pairs, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	sort.Strings(pairs)
	return fmt.Sprintf("%s{%s}", name, strings.Join(pairs, ","))
}

func formatFloat(f float64) string {
	if f == math.Trunc(f) {
		return fmt.Sprintf("%.1f", f)
	}
	return fmt.Sprintf("%g", f)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedHistKeys(m map[string]*Histogram) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
