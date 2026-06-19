package supervisor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"
	"github.com/getmitran/mitran/server/apierr"
)

type WorkerState string

const (
	Running  WorkerState = "running"
	Stopped  WorkerState = "stopped"
	Crashed  WorkerState = "crashed"
	Starting WorkerState = "starting"
)

type Worker struct {
	ID           string      `json:"id"`
	Cmd          string      `json:"cmd"`
	Args         []string    `json:"args"`
	PID          int         `json:"pid"`
	State        WorkerState `json:"state"`
	StartedAt    time.Time   `json:"started_at"`
	Uptime       string      `json:"uptime"`
	RestartCount int         `json:"restart_count"`
	LastCrash    *time.Time  `json:"last_crash,omitempty"`

	process *exec.Cmd
	stopCh  chan struct{}
}

type Manager struct {
	mu      sync.RWMutex
	workers map[string]*Worker
	nextID  int
}

func New() *Manager {
	return &Manager{
		workers: make(map[string]*Worker),
	}
}

func (m *Manager) Start(cmd string, args []string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nextID++
	id := fmt.Sprintf("worker-%d", m.nextID)

	w := &Worker{
		ID:     id,
		Cmd:    cmd,
		Args:   args,
		State:  Starting,
		stopCh: make(chan struct{}),
	}

	if err := m.startProcess(w); err != nil {
		return "", err
	}

	m.workers[id] = w
	go m.monitor(w)
	return id, nil
}

func (m *Manager) Stop(id string) error {
	m.mu.Lock()
	w, ok := m.workers[id]
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("worker %s not found", id)
	}

	close(w.stopCh)
	if w.process != nil && w.process.Process != nil {
		_ = w.process.Process.Kill()
	}

	m.mu.Lock()
	w.State = Stopped
	m.mu.Unlock()
	return nil
}

func (m *Manager) RestartAll() {
	m.mu.RLock()
	ids := make([]string, 0, len(m.workers))
	for id := range m.workers {
		ids = append(ids, id)
	}
	m.mu.RUnlock()

	for _, id := range ids {
		m.mu.Lock()
		w := m.workers[id]
		if w.process != nil && w.process.Process != nil {
			_ = w.process.Process.Kill()
		}
		w.stopCh = make(chan struct{})
		_ = m.startProcess(w)
		w.RestartCount++
		m.mu.Unlock()
		go m.monitor(w)
	}
}

func (m *Manager) startProcess(w *Worker) error {
	proc := exec.Command(w.Cmd, w.Args...)
	if err := proc.Start(); err != nil {
		w.State = Crashed
		return err
	}
	w.process = proc
	w.PID = proc.Process.Pid
	w.State = Running
	w.StartedAt = time.Now()
	return nil
}

func (m *Manager) monitor(w *Worker) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	done := make(chan error, 1)
	go func() { done <- w.process.Wait() }()

	for {
		select {
		case <-w.stopCh:
			return
		case <-done:
			m.mu.Lock()
			if w.State == Stopped {
				m.mu.Unlock()
				return
			}
			now := time.Now()
			w.State = Crashed
			w.LastCrash = &now
			w.RestartCount++

			backoff := time.Duration(w.RestartCount) * 2 * time.Second
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
			m.mu.Unlock()

			time.Sleep(backoff)

			m.mu.Lock()
			select {
			case <-w.stopCh:
				m.mu.Unlock()
				return
			default:
			}
			if err := m.startProcess(w); err != nil {
				m.mu.Unlock()
				return
			}
			m.mu.Unlock()

			done = make(chan error, 1)
			go func() { done <- w.process.Wait() }()

		case <-ticker.C:
			m.mu.RLock()
			if w.State == Running && w.process != nil && w.process.ProcessState != nil {
				m.mu.RUnlock()
				m.mu.Lock()
				w.State = Crashed
				m.mu.Unlock()
			} else {
				m.mu.RUnlock()
			}
		}
	}
}

func (m *Manager) List() []Worker {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Worker, 0, len(m.workers))
	for _, w := range m.workers {
		copy := *w
		if copy.State == Running {
			copy.Uptime = time.Since(copy.StartedAt).Truncate(time.Second).String()
		}
		result = append(result, copy)
	}
	return result
}

// RegisterRoutes adds supervisor endpoints to a mux.
func (m *Manager) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/workers", m.handleList)
	mux.HandleFunc("/api/v1/workers/restart", m.handleRestart)
}

func (m *Manager) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.List())
}

func (m *Manager) handleRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	m.RestartAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restarted"})
}
