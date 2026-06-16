package scheduler

import (
	"fmt"

	"github.com/getmitran/mitran/dag"
	"github.com/getmitran/mitran/lock"
	"github.com/getmitran/mitran/queue"
	"github.com/getmitran/mitran/types"
)

type Scheduler struct {
	dag   *dag.DAG
	queue *queue.PriorityQueue
	locks *lock.Manager
	tasks map[string]*types.Task
}

func New() *Scheduler {
	return &Scheduler{
		dag:   dag.New(),
		queue: queue.New(),
		locks: lock.NewManager(),
		tasks: make(map[string]*types.Task),
	}
}

func (s *Scheduler) AddTask(task *types.Task) error {
	s.tasks[task.ID] = task
	s.dag.AddNode(task.ID)
	for _, dep := range task.Deps {
		if err := s.dag.AddEdge(task.ID, dep); err != nil {
			return fmt.Errorf("adding dep edge: %w", err)
		}
	}
	return nil
}

func (s *Scheduler) HasConflict(task *types.Task) bool {
	for _, res := range task.Resources {
		if s.locks.IsLocked(res) && s.locks.HeldBy(res) != task.ID {
			return true
		}
	}
	return false
}

// Schedule returns the next executable task (ready deps, no conflicts, highest priority).
func (s *Scheduler) Schedule() (*types.Task, error) {
	if s.dag.DetectCycle() {
		return nil, fmt.Errorf("cycle detected in task graph")
	}

	ready := s.dag.GetReady()

	// rebuild queue with ready tasks
	pq := queue.New()
	for _, id := range ready {
		t := s.tasks[id]
		if t.Status == types.Pending {
			pq.Insert(id, t.Priority)
		}
	}

	for pq.Len() > 0 {
		id, ok := pq.Pop()
		if !ok {
			break
		}
		task := s.tasks[id]
		if !s.HasConflict(task) {
			// acquire locks
			for _, res := range task.Resources {
				s.locks.Acquire(res, task.ID)
			}
			task.Status = types.Running
			return task, nil
		}
	}
	return nil, nil
}

func (s *Scheduler) Complete(taskID string) error {
	task, ok := s.tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found: %s", taskID)
	}
	task.Status = types.Completed
	for _, res := range task.Resources {
		s.locks.Release(res, taskID)
	}
	s.dag.MarkCompleted(taskID)
	return nil
}
