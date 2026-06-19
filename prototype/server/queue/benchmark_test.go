package queue

import (
	"fmt"
	"testing"
)

func BenchmarkTaskQueueEnqueue(b *testing.B) {
	q := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(Task{ID: fmt.Sprintf("task-%d", i), Payload: nil})
	}
}
