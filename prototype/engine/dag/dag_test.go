package dag

import (
	"testing"
)

func TestAddNodeAndEdge(t *testing.T) {
	d := New()
	d.AddNode("a")
	d.AddNode("b")
	if err := d.AddEdge("a", "b"); err != nil {
		t.Fatal(err)
	}
	if err := d.AddEdge("a", "x"); err == nil {
		t.Fatal("expected error for missing node")
	}
}

func TestDetectCycle(t *testing.T) {
	d := New()
	d.AddNode("a")
	d.AddNode("b")
	d.AddNode("c")
	d.AddEdge("a", "b")
	d.AddEdge("b", "c")

	if d.DetectCycle() {
		t.Fatal("no cycle expected")
	}

	d.AddEdge("c", "a")
	if !d.DetectCycle() {
		t.Fatal("cycle expected")
	}
}

func TestTopologicalSort(t *testing.T) {
	d := New()
	d.AddNode("a")
	d.AddNode("b")
	d.AddNode("c")
	d.AddEdge("a", "b") // a depends on b
	d.AddEdge("b", "c") // b depends on c

	sorted, err := d.TopologicalSort()
	if err != nil {
		t.Fatal(err)
	}

	idx := make(map[string]int)
	for i, n := range sorted {
		idx[n] = i
	}
	// c before b, b before a
	if idx["c"] > idx["b"] || idx["b"] > idx["a"] {
		t.Fatalf("bad order: %v", sorted)
	}
}

func TestTopologicalSortCycle(t *testing.T) {
	d := New()
	d.AddNode("a")
	d.AddNode("b")
	d.AddEdge("a", "b")
	d.AddEdge("b", "a")

	_, err := d.TopologicalSort()
	if err == nil {
		t.Fatal("expected error for cycle")
	}
}

func TestGetReady(t *testing.T) {
	d := New()
	d.AddNode("a")
	d.AddNode("b")
	d.AddNode("c")
	d.AddEdge("a", "b") // a depends on b
	d.AddEdge("b", "c") // b depends on c

	ready := d.GetReady()
	if len(ready) != 1 || ready[0] != "c" {
		t.Fatalf("expected [c], got %v", ready)
	}

	d.MarkCompleted("c")
	ready = d.GetReady()
	if len(ready) != 1 || ready[0] != "b" {
		t.Fatalf("expected [b], got %v", ready)
	}

	d.MarkCompleted("b")
	ready = d.GetReady()
	if len(ready) != 1 || ready[0] != "a" {
		t.Fatalf("expected [a], got %v", ready)
	}
}

func TestGetReadyParallel(t *testing.T) {
	d := New()
	d.AddNode("a")
	d.AddNode("b")
	d.AddNode("c")
	// a and b have no deps, c depends on both
	d.AddEdge("c", "a")
	d.AddEdge("c", "b")

	ready := d.GetReady()
	if len(ready) != 2 {
		t.Fatalf("expected 2 ready nodes, got %v", ready)
	}
}
