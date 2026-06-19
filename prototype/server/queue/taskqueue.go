// TaskQueue is the in-memory bounded queue used by the dispatcher.
// For persistent/durable queuing, see PersistentQueue in persistent.go.
package queue

import (
	"errors"
	"os"
	"strconv"
)

var ErrQueueFull = errors.New("task queue full")

type Task struct {
	ID      string
	Payload interface{}
}

type TaskQueue struct {
	ch chan Task
}

func New() *TaskQueue {
	size := 100
	if v := os.Getenv("TASK_QUEUE_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			size = n
		}
	}
	return &TaskQueue{ch: make(chan Task, size)}
}

func (q *TaskQueue) Enqueue(task Task) error {
	select {
	case q.ch <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *TaskQueue) Dequeue() Task { return <-q.ch }
func (q *TaskQueue) Len() int      { return len(q.ch) }
func (q *TaskQueue) IsFull() bool   { return len(q.ch) == cap(q.ch) }
