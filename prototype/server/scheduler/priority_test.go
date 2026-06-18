package scheduler

import (
	"testing"
	"time"
)

func TestEffectivePriority_Base(t *testing.T) {
	item := &PriorityItem{Priority: 3, CreatedAt: time.Now()}
	if item.EffectivePriority() != 3.0 { t.Errorf("expected 3.0, got %f", item.EffectivePriority()) }
}

func TestEffectivePriority_Aging(t *testing.T) {
	item := &PriorityItem{Priority: 3, CreatedAt: time.Now().Add(-2 * time.Hour)}
	ep := item.EffectivePriority()
	if ep >= 3.0 { t.Errorf("aged item should have lower effective priority, got %f", ep) }
}

func TestEffectivePriority_Deadline(t *testing.T) {
	deadline := time.Now().Add(30 * time.Minute)
	item := &PriorityItem{Priority: 3, CreatedAt: time.Now(), Deadline: &deadline}
	ep := item.EffectivePriority()
	if ep >= 2.0 { t.Errorf("deadline-urgent item should be boosted, got %f", ep) }
}

func TestPriorityQueue_EnqueueDequeue(t *testing.T) {
	pq := NewPriorityQueue()
	pq.Enqueue(&PriorityItem{TaskID: "low", Priority: 5, CreatedAt: time.Now()})
	pq.Enqueue(&PriorityItem{TaskID: "high", Priority: 1, CreatedAt: time.Now()})
	pq.Enqueue(&PriorityItem{TaskID: "med", Priority: 3, CreatedAt: time.Now()})
	first := pq.Dequeue()
	if first.TaskID != "high" { t.Errorf("expected high priority first, got %s", first.TaskID) }
}

func TestPriorityQueue_Empty(t *testing.T) {
	pq := NewPriorityQueue()
	if pq.Dequeue() != nil { t.Error("expected nil from empty queue") }
	if pq.Peek() != nil { t.Error("expected nil peek from empty queue") }
}
