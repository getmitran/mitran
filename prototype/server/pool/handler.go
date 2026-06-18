package pool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RegisterHandlers(mux *http.ServeMux, p *Pool) {
	mux.HandleFunc("GET /api/v1/pool/stats", handleStats(p))
	mux.HandleFunc("GET /api/v1/pool/items", handleItems(p))
	mux.HandleFunc("POST /api/v1/pool/submit", handleSubmit(p))
}

func handleStats(p *Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queued, running, done, failed := p.Stats()
		json.NewEncoder(w).Encode(map[string]int{
			"workers": p.workers,
			"queued":  queued,
			"running": running,
			"done":    done,
			"failed":  failed,
		})
	}
}

func handleItems(p *Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(p.List())
	}
}

type submitReq struct {
	ID      string `json:"id"`
	Agent   string `json:"agent"`
	Payload string `json:"payload"`
}

func handleSubmit(p *Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req submitReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.ID == "" {
			req.ID = fmt.Sprintf("task-%d", time.Now().UnixNano())
		}
		// Placeholder task fn — real dispatch wires into agent worker
		p.Submit(req.ID, req.Agent, req.Payload, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"id": req.ID, "status": "queued"})
	}
}
