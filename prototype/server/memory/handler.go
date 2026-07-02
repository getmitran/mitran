package memory

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handler provides REST endpoints for the memory system.
type Handler struct {
	store *MemoryStore
}

// NewHandler creates a new memory HTTP handler.
func NewHandler(store *MemoryStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes registers memory routes on a mux. Expects {id} to be extracted by caller.
// Routes:
//   GET/POST/DELETE /api/v1/projects/{id}/memory/facts
//   GET/POST/DELETE /api/v1/projects/{id}/memory/lessons
//   GET/POST        /api/v1/projects/{id}/memory/episodes
//   GET             /api/v1/projects/{id}/memory/search?q=&limit=
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract project ID and sub-path from URL
	// Expected path: /api/v1/projects/{id}/memory/{resource}
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// Find "projects" segment and extract ID
	var projectID, resource string
	for i, p := range parts {
		if p == "projects" && i+1 < len(parts) {
			projectID = parts[i+1]
			// resource is after "memory"
			for j := i + 2; j < len(parts); j++ {
				if parts[j] == "memory" && j+1 < len(parts) {
					resource = parts[j+1]
					break
				}
			}
			break
		}
	}

	if projectID == "" || resource == "" {
		http.Error(w, `{"error":"invalid path"}`, http.StatusBadRequest)
		return
	}

	switch resource {
	case "facts":
		h.handleFacts(w, r, projectID)
	case "lessons":
		h.handleLessons(w, r, projectID)
	case "episodes":
		h.handleEpisodes(w, r, projectID)
	case "search":
		h.handleSearch(w, r, projectID)
	default:
		http.Error(w, `{"error":"unknown resource"}`, http.StatusNotFound)
	}
}

func (h *Handler) handleFacts(w http.ResponseWriter, r *http.Request, projectID string) {
	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("key")
		if key != "" {
			fact, err := h.store.GetFact(projectID, key)
			if err != nil {
				writeError(w, err)
				return
			}
			if fact == nil {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, fact)
		} else {
			facts, err := h.store.ListFacts(projectID)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, facts)
		}

	case http.MethodPost:
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Key == "" {
			http.Error(w, `{"error":"key required"}`, http.StatusBadRequest)
			return
		}
		fact, err := h.store.AddFact(projectID, req.Key, req.Value)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, fact)

	case http.MethodDelete:
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, `{"error":"key query param required"}`, http.StatusBadRequest)
			return
		}
		if err := h.store.DeleteFact(projectID, key); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleLessons(w http.ResponseWriter, r *http.Request, projectID string) {
	switch r.Method {
	case http.MethodGet:
		scope := r.URL.Query().Get("scope")
		lessons, err := h.store.ListLessons(projectID, scope)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, lessons)

	case http.MethodPost:
		var req struct {
			Rule     string `json:"rule"`
			Negative string `json:"negative"`
			Category string `json:"category"`
			Scope    string `json:"scope"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Rule == "" {
			http.Error(w, `{"error":"rule required"}`, http.StatusBadRequest)
			return
		}
		lesson, err := h.store.AddLesson(projectID, req.Rule, req.Negative, req.Category, req.Scope)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, lesson)

	case http.MethodDelete:
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, `{"error":"q query param required"}`, http.StatusBadRequest)
			return
		}
		count, err := h.store.RemoveLesson(projectID, query)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int64{"deleted": count})

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleEpisodes(w http.ResponseWriter, r *http.Request, projectID string) {
	switch r.Method {
	case http.MethodGet:
		// List recent episodes (default limit 20)
		episodes, err := h.store.SearchEpisodes(projectID, "", 20)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, episodes)

	case http.MethodPost:
		var req struct {
			Content   string `json:"content"`
			SessionID string `json:"session_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Content == "" {
			http.Error(w, `{"error":"content required"}`, http.StatusBadRequest)
			return
		}
		ep, err := h.store.AddEpisode(projectID, req.Content, req.SessionID)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, ep)

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, `{"error":"q query param required"}`, http.StatusBadRequest)
		return
	}
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	episodes, err := h.store.SearchEpisodes(projectID, query, limit)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, episodes)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
