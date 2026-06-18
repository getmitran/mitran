package memory

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	Store *MemoryStore
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case strings.HasSuffix(path, "/facts") && r.Method == "GET":
		h.listFacts(w, r)
	case strings.HasSuffix(path, "/facts") && r.Method == "POST":
		h.setFact(w, r)
	case strings.HasSuffix(path, "/search") && r.Method == "GET":
		h.search(w, r)
	case strings.HasSuffix(path, "/episodes") && r.Method == "POST":
		h.addEpisode(w, r)
	case strings.Contains(path, "/corrections/") && r.Method == "DELETE":
		parts := strings.Split(path, "/corrections/")
		if len(parts) == 2 && parts[1] != "" {
			h.removeCorrection(w, r, parts[1])
		} else {
			http.NotFound(w, r)
		}
	case strings.HasSuffix(path, "/corrections") && r.Method == "GET":
		h.listCorrections(w, r)
	case strings.HasSuffix(path, "/corrections") && r.Method == "POST":
		h.addCorrection(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) listFacts(w http.ResponseWriter, _ *http.Request) {
	h.Store.mu.RLock()
	defer h.Store.mu.RUnlock()
	writeJSON(w, 200, h.Store.Facts)
}

func (h *Handler) setFact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" {
		writeErr(w, 400, "key and value required")
		return
	}
	h.Store.SetFact(req.Key, req.Value)
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeErr(w, 400, "q parameter required")
		return
	}
	episodes := h.Store.SearchEpisodes(q)
	facts := h.Store.SearchFacts(q)
	if episodes == nil {
		episodes = []EpisodicEntry{}
	}
	writeJSON(w, 200, map[string]interface{}{"episodes": episodes, "facts": facts})
}

func (h *Handler) addEpisode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Text == "" {
		writeErr(w, 400, "text required")
		return
	}
	e := h.Store.AddEpisode(req.Text)
	writeJSON(w, 201, e)
}

func (h *Handler) listCorrections(w http.ResponseWriter, _ *http.Request) {
	c := h.Store.ListCorrections()
	if c == nil {
		c = []Correction{}
	}
	writeJSON(w, 200, c)
}

func (h *Handler) addCorrection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Rule     string `json:"rule"`
		Negative string `json:"negative"`
		Category string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Rule == "" {
		writeErr(w, 400, "rule required")
		return
	}
	c := h.Store.AddCorrection(req.Rule, req.Negative, req.Category)
	writeJSON(w, 201, c)
}

func (h *Handler) removeCorrection(w http.ResponseWriter, _ *http.Request, id string) {
	if h.Store.RemoveCorrection(id) {
		writeJSON(w, 200, map[string]string{"status": "removed"})
	} else {
		writeErr(w, 404, "not found")
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
