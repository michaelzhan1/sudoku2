package priorityqueue

// heapImpl implements the container/heap interface. It exists to be wrapped for type safety.
type heapImpl[T any] struct {
	items []T
	less  func(a, b T) bool
}

func (h *heapImpl[T]) Len() int {
	return len(h.items)
}

func (h *heapImpl[T]) Less(i, j int) bool {
	return h.less(h.items[i], h.items[j])
}

func (h *heapImpl[T]) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *heapImpl[T]) Push(x any) {
	h.items = append(h.items, x.(T))
}

func (h *heapImpl[T]) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[0 : n-1]
	return x
}
