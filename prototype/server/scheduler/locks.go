package scheduler

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrLockTimeout = errors.New("lock acquisition timed out")
	ErrDeadlock    = errors.New("potential deadlock detected")
)

type LockEntry struct {
	Owner    string
	Resource string
	Acquired time.Time
	Timeout  time.Duration
	Waiters  []string
}

type ResourceLocks struct {
	mu    sync.Mutex
	locks map[string]*LockEntry
}

func NewResourceLocks() *ResourceLocks {
	rl := &ResourceLocks{locks: make(map[string]*LockEntry)}
	go rl.expireLoop()
	return rl
}

func (rl *ResourceLocks) Acquire(resource, owner string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		rl.mu.Lock()
		entry, held := rl.locks[resource]
		if !held || time.Since(entry.Acquired) > entry.Timeout {
			rl.locks[resource] = &LockEntry{
				Owner: owner, Resource: resource,
				Acquired: time.Now(), Timeout: timeout,
			}
			rl.mu.Unlock()
			return nil
		}
		if rl.wouldDeadlock(owner, entry.Owner) {
			rl.mu.Unlock()
			return ErrDeadlock
		}
		entry.Waiters = append(entry.Waiters, owner)
		rl.mu.Unlock()
		if time.Now().After(deadline) {
			return ErrLockTimeout
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (rl *ResourceLocks) Release(resource, owner string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if entry, ok := rl.locks[resource]; ok && entry.Owner == owner {
		delete(rl.locks, resource)
	}
}

func (rl *ResourceLocks) IsLocked(resource string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	_, ok := rl.locks[resource]
	return ok
}

func (rl *ResourceLocks) ListLocks() []LockEntry {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	var result []LockEntry
	for _, e := range rl.locks {
		result = append(result, *e)
	}
	return result
}

func (rl *ResourceLocks) wouldDeadlock(requester, holder string) bool {
	// Cycle detection: check if holder is waiting on requester
	for _, entry := range rl.locks {
		if entry.Owner == requester {
			for _, w := range entry.Waiters {
				if w == holder {
					return true
				}
			}
		}
	}
	return false
}

func (rl *ResourceLocks) expireLoop() {
	for {
		time.Sleep(10 * time.Second)
		rl.mu.Lock()
		for resource, entry := range rl.locks {
			if time.Since(entry.Acquired) > entry.Timeout {
				delete(rl.locks, resource)
			}
		}
		rl.mu.Unlock()
	}
}
