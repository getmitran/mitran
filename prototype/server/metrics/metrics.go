package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Metrics struct {
	RequestCount       int64
	RequestDurationSum int64
	ErrorCount         int64
}

func (m *Metrics) RecordRequest(d time.Duration) {
	atomic.AddInt64(&m.RequestCount, 1)
	atomic.AddInt64(&m.RequestDurationSum, int64(d))
}

func (m *Metrics) RecordError() { atomic.AddInt64(&m.ErrorCount, 1) }

func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "# HELP mitran_requests_total Total requests\n# TYPE mitran_requests_total counter\nmitran_requests_total %d\n", atomic.LoadInt64(&m.RequestCount))
		fmt.Fprintf(w, "# HELP mitran_request_duration_seconds_sum Total request duration\n# TYPE mitran_request_duration_seconds_sum counter\nmitran_request_duration_seconds_sum %f\n", float64(atomic.LoadInt64(&m.RequestDurationSum))/float64(time.Second))
		fmt.Fprintf(w, "# HELP mitran_errors_total Total errors\n# TYPE mitran_errors_total counter\nmitran_errors_total %d\n", atomic.LoadInt64(&m.ErrorCount))
	}
}
