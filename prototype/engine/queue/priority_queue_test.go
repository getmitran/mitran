package queue

import "testing"

func TestInsertAndPop(t *testing.T) {
	q := New()
	q.Insert("low", 10)
	q.Insert("high", 1)
	q.Insert("mid", 5)

	id, ok := q.Pop()
	if !ok || id != "high" {
		t.Fatalf("expected high, got %s", id)
	}
	id, _ = q.Pop()
	if id != "mid" {
		t.Fatalf("expected mid, got %s", id)
	}
	id, _ = q.Pop()
	if id != "low" {
		t.Fatalf("expected low, got %s", id)
	}
	_, ok = q.Pop()
	if ok {
		t.Fatal("expected empty")
	}
}

func TestPeek(t *testing.T) {
	q := New()
	_, _, ok := q.Peek()
	if ok {
		t.Fatal("expected empty")
	}
	q.Insert("a", 3)
	q.Insert("b", 1)
	id, pri, ok := q.Peek()
	if !ok || id != "b" || pri != 1 {
		t.Fatalf("expected b/1, got %s/%d", id, pri)
	}
}

func TestReorder(t *testing.T) {
	q := New()
	q.Insert("a", 10)
	q.Insert("b", 5)

	q.Reorder("a", 1)
	id, _ := q.Pop()
	if id != "a" {
		t.Fatalf("expected a after reorder, got %s", id)
	}
}

func TestReorderMissing(t *testing.T) {
	q := New()
	if q.Reorder("x", 1) {
		t.Fatal("expected false for missing item")
	}
}

func TestLen(t *testing.T) {
	q := New()
	if q.Len() != 0 {
		t.Fatal("expected 0")
	}
	q.Insert("a", 1)
	q.Insert("b", 2)
	if q.Len() != 2 {
		t.Fatalf("expected 2, got %d", q.Len())
	}
	q.Pop()
	if q.Len() != 1 {
		t.Fatalf("expected 1, got %d", q.Len())
	}
}
