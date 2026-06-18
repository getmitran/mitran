package scheduler

import (
	"container/heap"
	"sync"
	"time"
)

type PriorityItem struct {
	TaskID    string
	Priority  int       // 1=critical, 5=lowest
	CreatedAt time.Time
	Deadline  *time.Time
	Agent     string
	index     int
}

func (p *PriorityItem) EffectivePriority() float64 {
	base := float64(p.Priority)
	// Age boost: tasks waiting >1h get priority boost
	age := time.Since(p.CreatedAt).Hours()
	if age > 1 {
		base -= age * 0.1
	}
	// Deadline urgency: tasks near deadline get boosted
	if p.Deadline != nil {
		untilDeadline := time.Until(*p.Deadline).Hours()
		if untilDeadline < 1 {
			base -= 2
		} else if untilDeadline < 4 {
			base -= 1
		}
	}
	if base < 0 {
		base = 0
	}
	return base
}

type PriorityQueue struct {
	mu    sync.Mutex
	items []*PriorityItem
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(pq)
	return pq
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*PriorityItem)
	item.index = len(pq.items)
	pq.items = append(pq.items, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[:n-1]
	item.index = -1
	return item
}

func (pq *PriorityQueue) Len() int { return len(pq.items) }
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].EffectivePriority() < pq.items[j].EffectivePriority()
}
func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Enqueue(item *PriorityItem) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(pq, item)
}

func (pq *PriorityQueue) Dequeue() *PriorityItem {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*PriorityItem)
}

func (pq *PriorityQueue) Peek() *PriorityItem {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.Len() == 0 {
		return nil
	}
	return pq.items[0]
}
