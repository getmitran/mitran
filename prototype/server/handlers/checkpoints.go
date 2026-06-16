package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type CheckpointHandler struct {
	Store *db.Store
}

func (h *CheckpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/checkpoints")
	switch {
	case r.Method == "GET" && path == "":
		h.list(w, r)
	case r.Method == "POST" && strings.HasSuffix(path, "/approve"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/approve")
		h.decide(w, r, id, "approved")
	case r.Method == "POST" && strings.HasSuffix(path, "/reject"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/reject")
		h.decide(w, r, id, "rejected")
	default:
		http.NotFound(w, r)
	}
}

func (h *CheckpointHandler) list(w http.ResponseWriter, r *http.Request) {
	cps := h.Store.ListPendingCheckpoints()
	if cps == nil {
		cps = []db.Checkpoint{}
	}
	writeJSON(w, http.StatusOK, cps)
}

func (h *CheckpointHandler) decide(w http.ResponseWriter, r *http.Request, id, action string) {
	var req struct {
		Feedback string `json:"feedback"`
		By       string `json:"by"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	now := time.Now()
	err := h.Store.UpdateCheckpoint(id, func(c *db.Checkpoint) {
		c.DecisionAction = action
		c.DecisionFeedback = req.Feedback
		c.DecisionBy = req.By
		c.DecisionAt = &now
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}

	// Update task status based on decision
	cp := h.Store.GetCheckpoint(id)
	if cp != nil {
		taskStatus := db.StatusApproved
		if action == "rejected" {
			taskStatus = db.StatusRejected
		}
		h.Store.UpdateTask(cp.TaskID, func(t *db.Task) {
			t.Status = taskStatus
			now2 := time.Now()
			t.CompletedAt = &now2
		})
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": action})
}
