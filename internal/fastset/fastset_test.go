package fastset_test

import (
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/fastset"
)

func TestNewFastSet(t *testing.T) {
	fs := fastset.NewFastSet[int]()
	if fs == nil {
		t.Errorf("Expected non-nil FastSet, got nil")
	}

	if fs.Size() != 0 {
		t.Errorf("Expected size 0, got %d", fs.Size())
	}
}

func TestOperations(t *testing.T) {
	type operation string

	var (
		addOp      operation = "add"
		removeOp   operation = "remove"
		containsOp operation = "contains"
		sizeOp     operation = "size"
		peekOp     operation = "peek"
		clearOp    operation = "clear"
		toSliceOp  operation = "toSlice"
	)

	type opGroup struct {
		op  operation
		val int
		exp any
	}

	type testcase struct {
		name string
		ops  []opGroup
	}

	testcases := []testcase{
		{
			name: "add_only",
			ops: []opGroup{
				{op: addOp, val: 1},
				{op: addOp, val: 2},
				{op: addOp, val: 3},
				{op: sizeOp, exp: 3},
				{op: containsOp, val: 2, exp: true},
				{op: containsOp, val: 4, exp: false},
				{op: peekOp, exp: 3},
				{op: toSliceOp, exp: []int{1, 2, 3}},
			},
		},
		{
			name: "add_and_remove",
			ops: []opGroup{
				{op: addOp, val: 1},
				{op: addOp, val: 2},
				{op: removeOp, val: 1},
				{op: sizeOp, exp: 1},
				{op: containsOp, val: 1, exp: false},
				{op: containsOp, val: 2, exp: true},
				{op: peekOp, exp: 2},
				{op: toSliceOp, exp: []int{2}},
			},
		},
		{
			name: "clear",
			ops: []opGroup{
				{op: addOp, val: 1},
				{op: addOp, val: 2},
				{op: clearOp},
				{op: sizeOp, exp: 0},
				{op: containsOp, val: 1, exp: false},
				{op: peekOp, exp: nil},
				{op: toSliceOp, exp: []int{}},
			},
		},
		{
			name: "add_duplicates",
			ops: []opGroup{
				{op: addOp, val: 1},
				{op: addOp, val: 1},
				{op: addOp, val: 2},
				{op: sizeOp, exp: 2},
				{op: containsOp, val: 1, exp: true},
				{op: containsOp, val: 2, exp: true},
				{op: peekOp, exp: 2},
				{op: toSliceOp, exp: []int{1, 2}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			fs := fastset.NewFastSet[int]()

			for _, op := range tc.ops {
				switch op.op {
				case addOp:
					fs.Add(op.val)
				case removeOp:
					fs.Remove(op.val)
				case containsOp:
					result := fs.Contains(op.val)
					if result != op.exp {
						t.Errorf("Contains(%d) = %v; want %v", op.val, result, op.exp)
					}
				case sizeOp:
					result := fs.Size()
					if result != op.exp {
						t.Errorf("Size() = %d; want %d", result, op.exp)
					}
				case peekOp:
					result, ok := fs.Peek()
					if !ok && op.exp != nil {
						t.Errorf("Peek() = nil; want %v", op.exp)
					} else if ok && result != op.exp {
						t.Errorf("Peek() = %v; want %v", result, op.exp)
					}
				case clearOp:
					fs.Clear()
				case toSliceOp:
					result := fs.ToSlice()
					if len(result) != len(op.exp.([]int)) {
						t.Errorf("ToSlice() length = %d; want %d", len(result), len(op.exp.([]int)))
					} else {
						for i, v := range result {
							if v != op.exp.([]int)[i] {
								t.Errorf("ToSlice()[%d] = %d; want %d", i, v, op.exp.([]int)[i])
							}
						}
					}
				default:
					t.Errorf("Unknown operation: %s", op.op)
				}
			}
		})
	}
}
