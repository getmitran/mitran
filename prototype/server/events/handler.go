package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func RegisterRoutes(mux *http.ServeMux, bus *Bus) {
	mux.HandleFunc("/api/v1/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		limit := 50
		if q := r.URL.Query().Get("limit"); q != "" {
			if v, err := strconv.Atoi(q); err == nil && v > 0 {
				limit = v
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bus.History(limit))
	})

	mux.HandleFunc("/api/v1/events/stream", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher.Flush()

		ch := make(chan Event, 64)
		bus.SubscribeAll(func(e Event) {
			select {
			case ch <- e:
			default:
			}
		})

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-ch:
				data, _ := json.Marshal(evt)
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, data)
				flusher.Flush()
			}
		}
	})
}
