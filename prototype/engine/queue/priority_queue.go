package queue

import "container/heap"

type Item struct {
	ID       string
	Priority int // lower = higher priority
	index    int
}

type pq []*Item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].Priority < p[j].Priority }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i]; p[i].index = i; p[j].index = j }
func (p *pq) Push(x interface{}) { item := x.(*Item); item.index = len(*p); *p = append(*p, item) }
func (p *pq) Pop() interface{}   { old := *p; n := len(old); item := old[n-1]; old[n-1] = nil; item.index = -1; *p = old[:n-1]; return item }

type PriorityQueue struct {
	items pq
	index map[string]*Item
}

func New() *PriorityQueue {
	pq := &PriorityQueue{
		items: make(pq, 0),
		index: make(map[string]*Item),
	}
	heap.Init(&pq.items)
	return pq
}

func (q *PriorityQueue) Insert(id string, priority int) {
	item := &Item{ID: id, Priority: priority}
	heap.Push(&q.items, item)
	q.index[id] = item
}

func (q *PriorityQueue) Pop() (string, bool) {
	if q.items.Len() == 0 {
		return "", false
	}
	item := heap.Pop(&q.items).(*Item)
	delete(q.index, item.ID)
	return item.ID, true
}

func (q *PriorityQueue) Peek() (string, int, bool) {
	if q.items.Len() == 0 {
		return "", 0, false
	}
	return q.items[0].ID, q.items[0].Priority, true
}

func (q *PriorityQueue) Reorder(id string, newPriority int) bool {
	item, ok := q.index[id]
	if !ok {
		return false
	}
	item.Priority = newPriority
	heap.Fix(&q.items, item.index)
	return true
}

func (q *PriorityQueue) Len() int {
	return q.items.Len()
}
