package session

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func RegisterRoutes(mux *http.ServeMux, mgr *Manager) {
	mux.HandleFunc("/api/v1/sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreate(w, r, mgr)
		case http.MethodGet:
			handleList(w, r, mgr)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/v1/sessions/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/sessions/")
		parts := strings.Split(path, "/")
		id := parts[0]

		if len(parts) == 2 && parts[1] == "messages" && r.Method == http.MethodPost {
			handleAddMessage(w, r, mgr, id)
			return
		}
		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				handleGet(w, r, mgr, id)
			case http.MethodDelete:
				handleEnd(w, r, mgr, id)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		http.NotFound(w, r)
	})
}

func handleCreate(w http.ResponseWriter, r *http.Request, mgr *Manager) {
	var req struct {
		UserID string `json:"user_id"`
		Agent  string `json:"agent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := mgr.Create(req.UserID, req.Agent)
	writeJSON(w, http.StatusCreated, s)
}

func handleList(w http.ResponseWriter, r *http.Request, mgr *Manager) {
	userID := r.URL.Query().Get("user_id")
	sessions := mgr.List(userID)
	if sessions == nil {
		sessions = []Session{}
	}
	writeJSON(w, http.StatusOK, sessions)
}

func handleGet(w http.ResponseWriter, _ *http.Request, mgr *Manager, id string) {
	s := mgr.Get(id)
	if s == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func handleAddMessage(w http.ResponseWriter, r *http.Request, mgr *Manager, sessionID string) {
	s := mgr.Get(sessionID)
	if s == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if msg.ID == "" {
		msg.ID = generateID()
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	mgr.AddMessage(sessionID, msg)
	writeJSON(w, http.StatusCreated, msg)
}

func handleEnd(w http.ResponseWriter, _ *http.Request, mgr *Manager, id string) {
	s := mgr.Get(id)
	if s == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	mgr.End(id)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
