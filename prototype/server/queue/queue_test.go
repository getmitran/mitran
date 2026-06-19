package queue

import "testing"

func TestNew(t *testing.T) {
	q := New()
	if q == nil {
		t.Fatal("New() returned nil")
	}
	if q.Len() != 0 {
		t.Fatalf("expected empty queue, got len %d", q.Len())
	}
}

func TestEnqueueDequeue(t *testing.T) {
	q := New()
	task := Task{ID: "t1", Payload: "data"}
	if err := q.Enqueue(task); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}
	got := q.Dequeue()
	if got.ID != task.ID || got.Payload != task.Payload {
		t.Fatalf("Dequeue returned %+v, want %+v", got, task)
	}
}

func TestIsFull(t *testing.T) {
	t.Setenv("TASK_QUEUE_SIZE", "2")
	q := New()
	q.Enqueue(Task{ID: "1"})
	q.Enqueue(Task{ID: "2"})
	if !q.IsFull() {
		t.Fatal("expected IsFull true at capacity")
	}
	if err := q.Enqueue(Task{ID: "3"}); err != ErrQueueFull {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}
}

func TestLen(t *testing.T) {
	t.Setenv("TASK_QUEUE_SIZE", "10")
	q := New()
	for i := 0; i < 5; i++ {
		q.Enqueue(Task{ID: string(rune('a' + i))})
	}
	if q.Len() != 5 {
		t.Fatalf("expected Len 5, got %d", q.Len())
	}
}
