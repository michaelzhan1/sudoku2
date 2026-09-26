package priorityqueue

import "container/heap"

// PriorityQueue is a minheap priority queue.
type PriorityQueue[T any] struct {
	heap *heapImpl[T]
}

func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: &heapImpl[T]{
			items: []T{},
			less:  less,
		},
	}
}

func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(pq.heap, item)
}

func (pq *PriorityQueue[T]) Pop() (T, bool) {
	if pq.heap.Len() == 0 {
		var zero T
		return zero, false
	}
	item := heap.Pop(pq.heap).(T)
	return item, true
}

func (pq *PriorityQueue[T]) Peek() (T, bool) {
	if pq.heap.Len() == 0 {
		var zero T
		return zero, false
	}
	return pq.heap.items[0], true
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.heap.Len() == 0
}
