package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	body    []byte
	status  int
	headers http.Header
	expiry  time.Time
}

type Cache struct {
	mu    sync.RWMutex
	store map[string]*cacheEntry
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{store: make(map[string]*cacheEntry), ttl: ttl}
	go c.cleanup()
	return c
}

func (c *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.URL.String()
		c.mu.RLock()
		if entry, ok := c.store[key]; ok && time.Now().Before(entry.expiry) {
			c.mu.RUnlock()
			for k, v := range entry.headers {
				w.Header()[k] = v
			}
			w.WriteHeader(entry.status)
			w.Write(entry.body)
			return
		}
		c.mu.RUnlock()
		rec := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}, status: 200}
		next.ServeHTTP(rec, r)
		if rec.status == 200 {
			c.mu.Lock()
			c.store[key] = &cacheEntry{body: rec.body.Bytes(), status: rec.status, headers: rec.Header().Clone(), expiry: time.Now().Add(c.ttl)}
			c.mu.Unlock()
		}
	})
}

func (c *Cache) cleanup() {
	for {
		time.Sleep(c.ttl)
		c.mu.Lock()
		for k, v := range c.store {
			if time.Now().After(v.expiry) {
				delete(c.store, k)
			}
		}
		c.mu.Unlock()
	}
}

type responseRecorder struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) WriteHeader(s int) {
	r.status = s
	r.ResponseWriter.WriteHeader(s)
}
