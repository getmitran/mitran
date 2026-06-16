package lock

import "sync"

type Manager struct {
	mu    sync.Mutex
	locks map[string]string // path -> held by task ID
}

func NewManager() *Manager {
	return &Manager{locks: make(map[string]string)}
}

func (m *Manager) Acquire(path, taskID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if holder, ok := m.locks[path]; ok && holder != taskID {
		return false
	}
	m.locks[path] = taskID
	return true
}

func (m *Manager) Release(path, taskID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.locks[path] != taskID {
		return false
	}
	delete(m.locks, path)
	return true
}

func (m *Manager) IsLocked(path string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.locks[path]
	return ok
}

func (m *Manager) HeldBy(path string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.locks[path]
}
