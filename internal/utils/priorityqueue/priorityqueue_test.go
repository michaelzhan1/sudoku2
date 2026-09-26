package priorityqueue_test

import (
	"math/rand/v2"
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/utils/priorityqueue"
)

func TestNewPriorityQueue(t *testing.T) {
	pq := priorityqueue.NewPriorityQueue[int](func(a, b int) bool { return a < b })
	if pq == nil {
		t.Errorf("Expected new priority queue to be non-nil")
	}
	if pq.IsEmpty() != true {
		t.Errorf("Expected new priority queue to be empty")
	}
}

func TestOperations(t *testing.T) {
	vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	t.Run("in_order", func(t *testing.T) {
		pq := priorityqueue.NewPriorityQueue(func(a, b int) bool { return a < b })
		for _, v := range vals {
			pq.Push(v)
		}

		for _, v := range vals {
			popped, ok := pq.Pop()
			if !ok {
				t.Errorf("Expected pop to succeed")
			}
			if popped != v {
				t.Errorf("Expected popped value %d, got %d", v, popped)
			}
		}

		if !pq.IsEmpty() {
			t.Errorf("Expected priority queue to be empty after popping all elements")
		}
		if _, ok := pq.Pop(); ok {
			t.Errorf("Expected pop to fail on empty priority queue")
		}
	})

	t.Run("random_order", func(t *testing.T) {
		valCopy := make([]int, len(vals))
		copy(valCopy, vals)
		rand.Shuffle(len(valCopy), func(i, j int) { valCopy[i], valCopy[j] = valCopy[j], valCopy[i] })
		pq := priorityqueue.NewPriorityQueue(func(a, b int) bool { return a < b })
		for _, v := range valCopy {
			pq.Push(v)
		}

		for _, v := range vals {
			popped, ok := pq.Pop()
			if !ok {
				t.Errorf("Expected pop to succeed")
			}
			if popped != v {
				t.Errorf("Expected popped value %d, got %d", v, popped)
			}
		}

		if !pq.IsEmpty() {
			t.Errorf("Expected priority queue to be empty after popping all elements")
		}
		if _, ok := pq.Pop(); ok {
			t.Errorf("Expected pop to fail on empty priority queue")
		}
	})

	t.Run("max_heap", func(t *testing.T) {
		pq := priorityqueue.NewPriorityQueue(func(a, b int) bool { return a > b })
		for _, v := range vals {
			pq.Push(v)
		}

		for i := len(vals) - 1; i >= 0; i-- {
			popped, ok := pq.Pop()
			if !ok {
				t.Errorf("Expected pop to succeed")
			}
			if popped != vals[i] {
				t.Errorf("Expected popped value %d, got %d", vals[i], popped)
			}
		}

		if !pq.IsEmpty() {
			t.Errorf("Expected priority queue to be empty after popping all elements")
		}
		if _, ok := pq.Pop(); ok {
			t.Errorf("Expected pop to fail on empty priority queue")
		}
	})

	t.Run("one_by_one", func(t *testing.T) {
		pq := priorityqueue.NewPriorityQueue(func(a, b int) bool { return a < b })
		for _, v := range vals {
			pq.Push(v)
			popped, ok := pq.Pop()
			if !ok {
				t.Errorf("Expected pop to succeed")
			}
			if popped != v {
				t.Errorf("Expected popped value %d, got %d", v, popped)
			}
		}

		if !pq.IsEmpty() {
			t.Errorf("Expected priority queue to be empty after popping all elements")
		}
		if _, ok := pq.Pop(); ok {
			t.Errorf("Expected pop to fail on empty priority queue")
		}
	})
}
