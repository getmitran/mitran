package health

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

var startTime = time.Now()

type Status struct {
	Status     string  `json:"status"`
	Uptime     string  `json:"uptime"`
	Version    string  `json:"version"`
	GoVersion  string  `json:"go_version"`
	Goroutines int     `json:"goroutines"`
	MemoryMB   float64 `json:"memory_mb"`
}

func Handler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		status := Status{
			Status: "healthy", Uptime: time.Since(startTime).String(),
			Version: version, GoVersion: runtime.Version(),
			Goroutines: runtime.NumGoroutine(), MemoryMB: float64(m.Alloc) / 1024 / 1024,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
