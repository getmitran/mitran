package hr

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"github.com/getmitran/mitran/server/apierr"
)

type Employee struct {
	ID, Name, Email, Department, Title, ManagerID, Phone, Location, JoinDate, Status string
	Active                                                                            bool
}

type Directory struct {
	mu   sync.RWMutex
	emps map[string]Employee
}

func NewDirectory() *Directory { return &Directory{emps: make(map[string]Employee)} }

func (d *Directory) Add(e Employee) { d.mu.Lock(); d.emps[e.ID] = e; d.mu.Unlock() }

func (d *Directory) Get(id string) *Employee {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if e, ok := d.emps[id]; ok {
		return &e
	}
	return nil
}

// Employees returns a snapshot of all employees (safe for iteration).
func (d *Directory) Employees() map[string]Employee {
	d.mu.RLock()
	defer d.mu.RUnlock()
	cp := make(map[string]Employee, len(d.emps))
	for k, v := range d.emps {
		cp[k] = v
	}
	return cp
}

func (d *Directory) Len() int { d.mu.RLock(); defer d.mu.RUnlock(); return len(d.emps) }

func (d *Directory) Search(query string) []Employee {
	d.mu.RLock()
	defer d.mu.RUnlock()
	q := strings.ToLower(query)
	var res []Employee
	for _, e := range d.emps {
		if strings.Contains(strings.ToLower(e.Name), q) || strings.Contains(strings.ToLower(e.Email), q) || strings.Contains(strings.ToLower(e.Department), q) {
			res = append(res, e)
		}
	}
	return res
}

func (d *Directory) ListByDepartment(dept string) []Employee {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var res []Employee
	for _, e := range d.emps {
		if strings.EqualFold(e.Department, dept) {
			res = append(res, e)
		}
	}
	return res
}

func (d *Directory) ListByManager(managerID string) []Employee {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var res []Employee
	for _, e := range d.emps {
		if e.ManagerID == managerID {
			res = append(res, e)
		}
	}
	return res
}

func (d *Directory) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/hr/employees")
	switch {
	case r.Method == "POST":
		var e Employee
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			apierr.WriteError(w, apierr.BadRequest("invalid request body"))
			return
		}
		d.Add(e)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(e)
	case r.Method == "GET" && path != "" && path != "/":
		if e := d.Get(strings.TrimPrefix(path, "/")); e != nil {
			json.NewEncoder(w).Encode(e)
		} else {
			http.NotFound(w, r)
		}
	case r.Method == "GET":
		if q := r.URL.Query().Get("q"); q != "" {
			json.NewEncoder(w).Encode(d.Search(q))
		} else {
			d.mu.RLock()
			all := make([]Employee, 0, len(d.emps))
			for _, e := range d.emps {
				all = append(all, e)
			}
			d.mu.RUnlock()
			json.NewEncoder(w).Encode(all)
		}
	}
}
