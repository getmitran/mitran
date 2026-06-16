package scheduler

import (
	"testing"

	"github.com/getmitran/mitran/types"
)

func TestBasicScheduling(t *testing.T) {
	s := New()
	s.AddTask(&types.Task{ID: "a", Priority: 5, Status: types.Pending})
	s.AddTask(&types.Task{ID: "b", Priority: 1, Status: types.Pending})

	task, err := s.Schedule()
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "b" {
		t.Fatalf("expected b (higher priority), got %s", task.ID)
	}
}

func TestDependencyOrder(t *testing.T) {
	s := New()
	s.AddTask(&types.Task{ID: "base", Priority: 10, Status: types.Pending})
	s.AddTask(&types.Task{ID: "child", Priority: 1, Status: types.Pending, Deps: []string{"base"}})

	task, _ := s.Schedule()
	if task.ID != "base" {
		t.Fatalf("expected base first, got %s", task.ID)
	}

	s.Complete("base")
	task, _ = s.Schedule()
	if task.ID != "child" {
		t.Fatalf("expected child after base complete, got %s", task.ID)
	}
}

func TestResourceConflict(t *testing.T) {
	s := New()
	s.AddTask(&types.Task{ID: "a", Priority: 1, Status: types.Pending, Resources: []string{"/src/main.go"}})
	s.AddTask(&types.Task{ID: "b", Priority: 2, Status: types.Pending, Resources: []string{"/src/main.go"}})

	task, _ := s.Schedule()
	if task.ID != "a" {
		t.Fatalf("expected a, got %s", task.ID)
	}

	// b should be blocked due to resource conflict
	task, _ = s.Schedule()
	if task != nil {
		t.Fatalf("expected nil (conflict), got %s", task.ID)
	}

	s.Complete("a")
	task, _ = s.Schedule()
	if task.ID != "b" {
		t.Fatalf("expected b after a released, got %s", task.ID)
	}
}

func TestHasConflict(t *testing.T) {
	s := New()
	taskA := &types.Task{ID: "a", Priority: 1, Resources: []string{"/x"}}
	taskB := &types.Task{ID: "b", Priority: 2, Resources: []string{"/x"}}
	s.AddTask(taskA)
	s.AddTask(taskB)

	if s.HasConflict(taskA) {
		t.Fatal("no conflict expected before scheduling")
	}

	s.Schedule() // schedules a, locks /x

	if !s.HasConflict(taskB) {
		t.Fatal("conflict expected for b")
	}
}

func TestCompleteUnknownTask(t *testing.T) {
	s := New()
	if err := s.Complete("nonexistent"); err == nil {
		t.Fatal("expected error")
	}
}

func TestNoTasksReady(t *testing.T) {
	s := New()
	task, err := s.Schedule()
	if err != nil {
		t.Fatal(err)
	}
	if task != nil {
		t.Fatal("expected nil")
	}
}
