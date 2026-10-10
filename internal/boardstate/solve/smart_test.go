package solve

import (
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/utils/fastset"
	"github.com/michaelzhan1/sudoku2/internal/utils/priorityqueue"
)

func TestResolveClosedGroupsInRow(t *testing.T) {
	testcases := []struct {
		name               string
		row                [9]int
		possibleCandidates map[int][]int
		exp                bool
	}{
		{
			name: "two_number_group",
			row:  [9]int{0, 0, 0, 4, 5, 6, 7, 8, 9},
			possibleCandidates: map[int][]int{
				0: {1, 2},
				1: {1, 2, 3},
				2: {3},
			},
			exp: true,
		},
		{
			name: "three_number_group",
			row:  [9]int{0, 0, 0, 0, 5, 6, 7, 8, 9},
			possibleCandidates: map[int][]int{
				0: {1, 2, 4},
				1: {1, 3, 4},
				2: {2, 4},
				3: {1},
			},
			exp: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			var possibleSets [9]*fastset.FastSet[int]
			for i := range 9 {
				possibleSets[i] = fastset.NewFastSet[int]()
			}
			for i, vals := range testcase.possibleCandidates {
				for _, val := range vals {
					possibleSets[i].Add(val)
				}
			}

			ss := &SudokuSolver{
				board: &board.Board{
					testcase.row,
				},
				possible: [9][9]*fastset.FastSet[int]{possibleSets},
				pq:       priorityqueue.NewPriorityQueue(func(a, b [2]int) bool { return a[0] < b[0] }),
			}

			res := ss.resolveHiddenGroupsInRow(0)
			if res != testcase.exp {
				t.Errorf("expected resolveClosedGroupsInRow to return %v, got %v", testcase.exp, res)
			}

			newSize := ss.possible[0][1].Size()
			if newSize != 2 {
				t.Errorf("expected possible candidates for cell 1 to have size 2, got %d", newSize)
			}
		})
	}
}

func TestResolveClosedGroupsInCol(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2, 4},
		{1, 0}: {1, 2, 5},
		{2, 0}: {3, 4, 5},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveHiddenGroupsInCol(0) {
		t.Fatal("expected hidden pair in column to be resolved")
	}

	assertCandidates(t, ss, [2]int{0, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{1, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{2, 0}, []int{3, 4, 5})
}

func TestResolveClosedGroupsInBox(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2, 3, 4, 5, 6, 7},
		{0, 1}: {1, 2, 3, 4, 5, 6, 7},
		{1, 0}: {3, 4, 5, 6, 7},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveHiddenGroupsInBox(0, 0) {
		t.Fatal("expected hidden pair in box to be resolved")
	}

	assertCandidates(t, ss, [2]int{0, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{0, 1}, []int{1, 2})
	assertCandidates(t, ss, [2]int{1, 0}, []int{3, 4, 5, 6, 7})
}

func TestResolveClosedGroupsDoesNotChangeWithoutGroup(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2},
		{0, 1}: {2, 3},
		{0, 2}: {1, 3},
	}
	ss := newSolverForCandidates(possible)

	if ss.resolveHiddenGroupsInRow(0) {
		t.Fatal("expected no hidden group to be resolved")
	}

	assertCandidates(t, ss, [2]int{0, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{0, 1}, []int{2, 3})
	assertCandidates(t, ss, [2]int{0, 2}, []int{1, 3})
}

func TestResolveNakedGroupsInRow(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2},
		{0, 1}: {1, 2},
		{0, 2}: {1, 3, 4},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveNakedGroupsInRow(0) {
		t.Fatal("expected naked pair in row to be resolved")
	}

	assertCandidates(t, ss, [2]int{0, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{0, 1}, []int{1, 2})
	assertCandidates(t, ss, [2]int{0, 2}, []int{3, 4})
}

func TestResolveNakedGroupsInBox(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2},
		{0, 1}: {1, 3},
		{0, 2}: {2, 3},
		{1, 0}: {3, 4},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveNakedGroupsInBox(0, 0) {
		t.Fatal("expected naked triple in box to be resolved")
	}

	assertCandidates(t, ss, [2]int{0, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{0, 1}, []int{1, 3})
	assertCandidates(t, ss, [2]int{0, 2}, []int{2, 3})
	assertCandidates(t, ss, [2]int{1, 0}, []int{4})
}

func TestResolveNakedGroupsDoesNotChangeWithoutGroup(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2},
		{0, 1}: {2, 3},
		{0, 2}: {1, 3},
	}
	ss := newSolverForCandidates(possible)

	if ss.resolveNakedGroupsInRow(0) {
		t.Fatal("expected no naked group to be resolved")
	}
}

func TestResolveSpearsInBoxUsesDistinctRowsAndColumns(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {5},
		{0, 1}: {5},
		{0, 3}: {5, 6},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveSpearsInBox(0, 0) {
		t.Fatal("expected pointing pair to be resolved")
	}

	assertCandidates(t, ss, [2]int{0, 3}, []int{6})
}

func TestResolveSpearsInBoxDoesNotUseDuplicateOccurrencesAsSingle(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {5},
		{1, 1}: {5},
		{0, 3}: {5, 6},
	}
	ss := newSolverForCandidates(possible)

	if ss.resolveSpearsInBox(0, 0) {
		t.Fatal("expected no pointing pair when candidates span rows and columns")
	}

	assertCandidates(t, ss, [2]int{0, 3}, []int{5, 6})
}

func TestResolveDoubleSpearsInBoxRowClaimsCandidate(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {6},
		{0, 1}: {6},
		{1, 2}: {6, 7},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveDoubleSpearsInBoxRow(0) {
		t.Fatal("expected claiming pair in row to be resolved")
	}

	assertCandidates(t, ss, [2]int{1, 2}, []int{7})
}

func TestResolveDoubleSpearsInBoxColClaimsCandidate(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {8},
		{1, 0}: {8},
		{2, 1}: {8, 9},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveDoubleSpearsInBoxCol(0) {
		t.Fatal("expected claiming pair in column to be resolved")
	}

	assertCandidates(t, ss, [2]int{2, 1}, []int{9})
}

func newSolverForCandidates(candidates map[[2]int][]int) *SudokuSolver {
	var possible [9][9]*fastset.FastSet[int]
	for row := range 9 {
		for col := range 9 {
			possible[row][col] = fastset.NewFastSet[int]()
		}
	}
	for cell, values := range candidates {
		for _, value := range values {
			possible[cell[0]][cell[1]].Add(value)
		}
	}

	return &SudokuSolver{
		board:    &board.Board{},
		possible: possible,
		pq:       priorityqueue.NewPriorityQueue(func(a, b [2]int) bool { return a[0] < b[0] }),
	}
}

func assertCandidates(t *testing.T, ss *SudokuSolver, cell [2]int, expected []int) {
	t.Helper()
	actual := ss.possible[cell[0]][cell[1]].ToSlice()
	if len(actual) != len(expected) {
		t.Fatalf("cell (%d, %d): expected candidates %v, got %v", cell[0], cell[1], expected, actual)
	}
	for _, value := range expected {
		if !ss.possible[cell[0]][cell[1]].Contains(value) {
			t.Errorf("cell (%d, %d): expected candidate %d in %v", cell[0], cell[1], value, actual)
		}
	}
}
