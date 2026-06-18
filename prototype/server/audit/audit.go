package audit

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
	IP        string    `json:"ip"`
}

type AuditFilter struct {
	UserID string
	Action string
	Since  time.Time
	Until  time.Time
}

type AuditLog struct {
	path string
	mu   sync.Mutex
}

func New() *AuditLog {
	path := os.Getenv("AUDIT_LOG_PATH")
	if path == "" {
		path = "./data/audit.jsonl"
	}
	os.MkdirAll("./data", 0755)
	return &AuditLog{path: path}
}

func (a *AuditLog) Log(e AuditEntry) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}

// readAllLocked reads all entries from the audit log file.
// Caller must hold a.mu.
func (a *AuditLog) readAllLocked() []AuditEntry {
	f, err := os.Open(a.path)
	if err != nil {
		return nil
	}
	defer f.Close()
	var entries []AuditEntry
	s := bufio.NewScanner(f)
	for s.Scan() {
		var e AuditEntry
		if json.Unmarshal(s.Bytes(), &e) == nil {
			entries = append(entries, e)
		}
	}
	return entries
}

func (a *AuditLog) Query(filter AuditFilter) []AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	var result []AuditEntry
	for _, e := range a.readAllLocked() {
		if filter.UserID != "" && e.UserID != filter.UserID {
			continue
		}
		if filter.Action != "" && e.Action != filter.Action {
			continue
		}
		if !filter.Since.IsZero() && e.Timestamp.Before(filter.Since) {
			continue
		}
		if !filter.Until.IsZero() && e.Timestamp.After(filter.Until) {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (a *AuditLog) Recent(n int) []AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	all := a.readAllLocked()
	if len(all) <= n {
		return all
	}
	return all[len(all)-n:]
}

func AuditMiddleware(log *AuditLog) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
				log.Log(AuditEntry{
					Timestamp: time.Now().UTC(),
					Action:    r.Method,
					Resource:  r.URL.Path,
					IP:        r.RemoteAddr,
				})
			}
			next.ServeHTTP(w, r)
		})
	}
}
