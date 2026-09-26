package fastset

type FastSet[T comparable] struct {
	arr    []T
	idxMap map[T]int
}

func NewFastSet[T comparable]() *FastSet[T] {
	return &FastSet[T]{
		arr:    make([]T, 0),
		idxMap: make(map[T]int),
	}
}

func (fs *FastSet[T]) Add(value T) {
	if _, exists := fs.idxMap[value]; !exists {
		fs.arr = append(fs.arr, value)
		fs.idxMap[value] = len(fs.arr) - 1
	}
}

func (fs *FastSet[T]) Remove(value T) {
	if idx, exists := fs.idxMap[value]; exists {
		lastIdx := len(fs.arr) - 1
		lastValue := fs.arr[lastIdx]

		fs.arr[idx] = lastValue
		fs.idxMap[lastValue] = idx

		fs.arr = fs.arr[:lastIdx]
		delete(fs.idxMap, value)
	}
}

func (fs *FastSet[T]) Contains(value T) bool {
	_, exists := fs.idxMap[value]
	return exists
}

func (fs *FastSet[T]) Size() int {
	return len(fs.arr)
}

func (fs *FastSet[T]) Peek() (T, bool) {
	if len(fs.arr) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	return fs.arr[len(fs.arr)-1], true
}

func (fs *FastSet[T]) Clear() {
	fs.arr = fs.arr[:0]
	fs.idxMap = make(map[T]int)
}

func (fs *FastSet[T]) ToSlice() []T {
	return fs.arr
}
