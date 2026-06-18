package pool

import (
	"log"
	"sync"
	"time"
)

type TaskFunc func() error

type WorkItem struct {
	ID        string     `json:"id"`
	Agent     string     `json:"agent"`
	Payload   string     `json:"payload"`
	Status    string     `json:"status"`
	Error     string     `json:"error,omitempty"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	DoneAt    *time.Time `json:"done_at,omitempty"`
	Fn        TaskFunc   `json:"-"`
}

type Pool struct {
	mu      sync.RWMutex
	workers int
	queue   chan *WorkItem
	items   []*WorkItem
	stopCh  chan struct{}
}

func New(workers int) *Pool {
	if workers < 1 {
		workers = 4
	}
	return &Pool{
		workers: workers,
		queue:   make(chan *WorkItem, 100),
		stopCh:  make(chan struct{}),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		go p.worker(i)
	}
	log.Printf("[pool] started %d workers", p.workers)
}

func (p *Pool) Stop() { close(p.stopCh) }

func (p *Pool) Submit(id, agent, payload string, fn TaskFunc) {
	item := &WorkItem{ID: id, Agent: agent, Payload: payload, Status: "queued", Fn: fn}
	p.mu.Lock()
	p.items = append(p.items, item)
	p.mu.Unlock()
	p.queue <- item
}

func (p *Pool) List() []*WorkItem {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return append([]*WorkItem{}, p.items...)
}

func (p *Pool) Stats() (queued, running, done, failed int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, item := range p.items {
		switch item.Status {
		case "queued":
			queued++
		case "running":
			running++
		case "done":
			done++
		case "failed":
			failed++
		}
	}
	return
}

func (p *Pool) worker(id int) {
	for {
		select {
		case <-p.stopCh:
			return
		case item := <-p.queue:
			now := time.Now()
			item.StartedAt = &now
			item.Status = "running"
			if err := item.Fn(); err != nil {
				item.Status = "failed"
				item.Error = err.Error()
			} else {
				item.Status = "done"
			}
			done := time.Now()
			item.DoneAt = &done
		}
	}
}
