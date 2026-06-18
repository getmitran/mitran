package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type TicketHandler struct {
	Store *db.Store
}

func (h *TicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tickets")
	switch {
	case r.Method == "GET" && path == "":
		h.list(w, r)
	case r.Method == "POST" && path == "":
		h.create(w, r)
	case r.Method == "GET" && !strings.Contains(path[1:], "/"):
		h.get(w, r, strings.TrimPrefix(path, "/"))
	case r.Method == "PUT" && !strings.Contains(path[1:], "/"):
		h.update(w, r, strings.TrimPrefix(path, "/"))
	case r.Method == "DELETE" && !strings.Contains(path[1:], "/"):
		h.delete(w, r, strings.TrimPrefix(path, "/"))
	case r.Method == "POST" && strings.HasSuffix(path, "/comments"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/comments")
		h.addComment(w, r, id)
	case r.Method == "GET" && strings.HasSuffix(path, "/comments"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/comments")
		h.listComments(w, r, id)
	default:
		http.NotFound(w, r)
	}
}

func (h *TicketHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	tickets := h.Store.ListTickets(q.Get("status"), q.Get("assignee"), q.Get("sprint_id"), q.Get("labels"))
	if tickets == nil {
		tickets = []db.Ticket{}
	}
	writeJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) create(w http.ResponseWriter, r *http.Request) {
	var t db.Ticket
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	t.ID = genID()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if t.Status == "" {
		t.Status = "open"
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}
	if err := h.Store.AddTicket(t); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *TicketHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	t := h.Store.GetTicket(id)
	if t == nil {
		writeErr(w, http.StatusNotFound, "ticket not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *TicketHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.Store.UpdateTicket(id, func(t *db.Ticket) {
		if v, ok := updates["title"].(string); ok {
			t.Title = v
		}
		if v, ok := updates["description"].(string); ok {
			t.Description = v
		}
		if v, ok := updates["status"].(string); ok {
			t.Status = v
		}
		if v, ok := updates["priority"].(string); ok {
			t.Priority = v
		}
		if v, ok := updates["assignee"].(string); ok {
			t.Assignee = v
		}
		if v, ok := updates["sprint_id"].(string); ok {
			t.SprintID = v
		}
		if v, ok := updates["sla_due_at"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, v); err == nil {
				t.SLADueAt = &parsed
			}
		}
		if v, ok := updates["labels"].([]interface{}); ok {
			labels := make([]string, 0, len(v))
			for _, l := range v {
				if s, ok := l.(string); ok {
					labels = append(labels, s)
				}
			}
			t.Labels = labels
		}
		t.UpdatedAt = time.Now()
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, h.Store.GetTicket(id))
}

func (h *TicketHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Store.DeleteTicket(id); err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *TicketHandler) addComment(w http.ResponseWriter, r *http.Request, ticketID string) {
	var c db.TicketComment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	c.ID = genID()
	c.TicketID = ticketID
	c.CreatedAt = time.Now()
	if err := h.Store.AddTicketComment(c); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *TicketHandler) listComments(w http.ResponseWriter, r *http.Request, ticketID string) {
	comments := h.Store.ListTicketComments(ticketID)
	if comments == nil {
		comments = []db.TicketComment{}
	}
	writeJSON(w, http.StatusOK, comments)
}

// --- Sprint Handler ---

type SprintHandler struct {
	Store *db.Store
}

func (h *SprintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/sprints")
	switch {
	case r.Method == "GET" && path == "":
		h.list(w, r)
	case r.Method == "POST" && path == "":
		h.create(w, r)
	case r.Method == "PUT" && path != "":
		h.update(w, r, strings.TrimPrefix(path, "/"))
	default:
		http.NotFound(w, r)
	}
}

func (h *SprintHandler) list(w http.ResponseWriter, r *http.Request) {
	sprints := h.Store.ListSprints()
	if sprints == nil {
		sprints = []db.Sprint{}
	}
	writeJSON(w, http.StatusOK, sprints)
}

func (h *SprintHandler) create(w http.ResponseWriter, r *http.Request) {
	var s db.Sprint
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	s.ID = genID()
	s.CreatedAt = time.Now()
	if err := h.Store.AddSprint(s); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

func (h *SprintHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.Store.UpdateSprint(id, func(s *db.Sprint) {
		if v, ok := updates["name"].(string); ok {
			s.Name = v
		}
		if v, ok := updates["start_date"].(string); ok {
			s.StartDate = v
		}
		if v, ok := updates["end_date"].(string); ok {
			s.EndDate = v
		}
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
