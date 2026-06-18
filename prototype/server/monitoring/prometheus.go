package monitoring

import (
	"fmt"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

// Metrics holds Mitran custom counters.
var (
	TasksTotal       atomic.Int64
	TasksActive      atomic.Int64
	AgentInvocations atomic.Int64
	LLMTokensUsed    atomic.Int64
)

// MetricsHandler returns an http.HandlerFunc that serves Prometheus text exposition format.
func MetricsHandler() http.HandlerFunc {
	startTime := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		// Go runtime metrics
		fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines.\n")
		fmt.Fprintf(w, "# TYPE go_goroutines gauge\ngo_goroutines %d\n", runtime.NumGoroutine())
		fmt.Fprintf(w, "# HELP go_memstats_alloc_bytes Number of bytes allocated and in use.\n")
		fmt.Fprintf(w, "# TYPE go_memstats_alloc_bytes gauge\ngo_memstats_alloc_bytes %d\n", m.Alloc)
		fmt.Fprintf(w, "# HELP go_memstats_sys_bytes Number of bytes obtained from system.\n")
		fmt.Fprintf(w, "# TYPE go_memstats_sys_bytes gauge\ngo_memstats_sys_bytes %d\n", m.Sys)
		fmt.Fprintf(w, "# HELP go_gc_duration_seconds Total GC pause duration.\n")
		fmt.Fprintf(w, "# TYPE go_gc_duration_seconds gauge\ngo_gc_duration_seconds %f\n", float64(m.PauseTotalNs)/1e9)
		fmt.Fprintf(w, "# HELP process_uptime_seconds Time since server start.\n")
		fmt.Fprintf(w, "# TYPE process_uptime_seconds gauge\nprocess_uptime_seconds %f\n", time.Since(startTime).Seconds())

		// Mitran custom metrics
		fmt.Fprintf(w, "# HELP mitran_tasks_total Total tasks created.\n")
		fmt.Fprintf(w, "# TYPE mitran_tasks_total counter\nmitran_tasks_total %d\n", TasksTotal.Load())
		fmt.Fprintf(w, "# HELP mitran_tasks_active Currently active tasks.\n")
		fmt.Fprintf(w, "# TYPE mitran_tasks_active gauge\nmitran_tasks_active %d\n", TasksActive.Load())
		fmt.Fprintf(w, "# HELP mitran_agent_invocations_total Total agent invocations.\n")
		fmt.Fprintf(w, "# TYPE mitran_agent_invocations_total counter\nmitran_agent_invocations_total %d\n", AgentInvocations.Load())
		fmt.Fprintf(w, "# HELP mitran_llm_tokens_used_total Total LLM tokens consumed.\n")
		fmt.Fprintf(w, "# TYPE mitran_llm_tokens_used_total counter\nmitran_llm_tokens_used_total %d\n", LLMTokensUsed.Load())
	}
}
