package queue

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type PersistentQueue interface {
	Enqueue(task Task) error
	Dequeue() (*Task, error)
	Len() int
	Close() error
}

// MemoryQueue is an in-memory implementation of PersistentQueue.
type MemoryQueue struct {
	mu    sync.Mutex
	tasks []Task
}

func (q *MemoryQueue) Enqueue(task Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)
	return nil
}

func (q *MemoryQueue) Dequeue() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tasks) == 0 {
		return nil, nil
	}
	t := q.tasks[0]
	q.tasks = q.tasks[1:]
	return &t, nil
}

func (q *MemoryQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.tasks)
}

func (q *MemoryQueue) Close() error { return nil }

// FileQueue persists tasks as JSON lines in a file.
type FileQueue struct {
	mu    sync.Mutex
	tasks []Task
	path  string
	file  *os.File
}

func NewFileQueue(path string) (*FileQueue, error) {
	if err := os.MkdirAll(filepath(path), 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var t Task
		if err := json.Unmarshal(scanner.Bytes(), &t); err == nil {
			tasks = append(tasks, t)
		}
	}
	return &FileQueue{tasks: tasks, path: path, file: f}, nil
}

func filepath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}

func (q *FileQueue) Enqueue(task Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	if _, err := q.file.Write(append(data, '\n')); err != nil {
		return err
	}
	q.tasks = append(q.tasks, task)
	return nil
}

func (q *FileQueue) Dequeue() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tasks) == 0 {
		return nil, nil
	}
	t := q.tasks[0]
	q.tasks = q.tasks[1:]
	// Rewrite file without dequeued task
	return &t, q.rewrite()
}

func (q *FileQueue) rewrite() error {
	tmp := q.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	for _, t := range q.tasks {
		data, _ := json.Marshal(t)
		f.Write(append(data, '\n'))
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	f.Close()
	if err := os.Rename(tmp, q.path); err != nil {
		return err
	}
	q.file.Close()
	q.file, _ = os.OpenFile(q.path, os.O_RDWR|os.O_APPEND, 0644)
	return nil
}

func (q *FileQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.tasks)
}

func (q *FileQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.file.Close()
}

// NewPersistentQueue creates a queue based on MITRAN_QUEUE_BACKEND env (memory|file).
func NewPersistentQueue() (PersistentQueue, error) {
	backend := os.Getenv("MITRAN_QUEUE_BACKEND")
	if backend == "file" {
		path := os.Getenv("MITRAN_QUEUE_FILE")
		if path == "" {
			path = "./data/queue.jsonl"
		}
		return NewFileQueue(path)
	}
	return &MemoryQueue{}, nil
}
