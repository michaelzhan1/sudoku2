package solve

import (
	"errors"
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
				pq:       priorityqueue.NewPriorityQueue(func(a, b queuedCell) bool { return a.possibilityCount < b.possibilityCount }),
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

func TestResolveNakedGroupsInCol(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2},
		{1, 0}: {1, 2},
		{2, 0}: {1, 3, 4},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveNakedGroupsInCol(0) {
		t.Fatal("expected naked pair in column to be resolved")
	}

	assertCandidates(t, ss, [2]int{0, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{1, 0}, []int{1, 2})
	assertCandidates(t, ss, [2]int{2, 0}, []int{3, 4})
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

func TestResolveSinglePossibilityPlacesOneSingleAtATime(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 0}: {1, 2},
		{0, 1}: {3, 4},
	}
	ss := newSolverForCandidates(possible)

	changed, err := ss.resolveSinglePossibility([][2]int{
		{0, 0},
		{0, 1},
	})
	if err != nil {
		t.Fatalf("resolveSinglePossibility returned an error: %v", err)
	}
	if !changed {
		t.Fatal("expected a hidden single to be placed")
	}

	placed := 0
	for _, cell := range [][2]int{{0, 0}, {0, 1}} {
		if ss.board[cell[0]][cell[1]] != 0 {
			placed++
		}
	}
	if placed != 1 {
		t.Fatalf("expected exactly one single to be placed, got %d", placed)
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

func TestResolveXWingRow(t *testing.T) {
	possible := map[[2]int][]int{
		{0, 2}: {5},
		{0, 7}: {5},
		{1, 2}: {5},
		{1, 7}: {5},
		{2, 2}: {4, 5},
		{2, 7}: {5, 6},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveXWingRow() {
		t.Fatal("expected x-wing in rows to be resolved")
	}

	assertCandidates(t, ss, [2]int{2, 2}, []int{4})
	assertCandidates(t, ss, [2]int{2, 7}, []int{6})
}

func TestResolveXWingCol(t *testing.T) {
	possible := map[[2]int][]int{
		{3, 0}: {5},
		{8, 0}: {5},
		{3, 1}: {5},
		{8, 1}: {5},
		{3, 5}: {4, 5},
		{8, 5}: {5, 6},
	}
	ss := newSolverForCandidates(possible)

	if !ss.resolveXWingCol() {
		t.Fatal("expected x-wing in columns to be resolved")
	}

	assertCandidates(t, ss, [2]int{3, 5}, []int{4})
	assertCandidates(t, ss, [2]int{8, 5}, []int{6})
}

func TestFillSinglePossibilityInUnits(t *testing.T) {
	testcases := []struct {
		name       string
		candidates map[[2]int][]int
		resolve    func(*SudokuSolver) (bool, error)
		expected   [2]int
	}{
		{
			name: "row",
			candidates: map[[2]int][]int{
				{0, 0}: {3},
				{0, 1}: {1, 2},
				{0, 2}: {1, 2},
			},
			resolve:  func(ss *SudokuSolver) (bool, error) { return ss.fillSinglePossibilityInRow(0) },
			expected: [2]int{0, 0},
		},
		{
			name: "column",
			candidates: map[[2]int][]int{
				{0, 0}: {3},
				{1, 0}: {1, 2},
				{2, 0}: {1, 2},
			},
			resolve:  func(ss *SudokuSolver) (bool, error) { return ss.fillSinglePossibilityInCol(0) },
			expected: [2]int{0, 0},
		},
		{
			name: "box",
			candidates: map[[2]int][]int{
				{0, 0}: {3},
				{0, 1}: {1, 2},
				{0, 2}: {1, 2},
			},
			resolve:  func(ss *SudokuSolver) (bool, error) { return ss.fillSinglePossibilityInBox(0, 0) },
			expected: [2]int{0, 0},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			ss := newSolverForCandidates(testcase.candidates)

			changed, err := testcase.resolve(ss)
			if err != nil {
				t.Fatalf("fillSinglePossibility returned an error: %v", err)
			}
			if !changed {
				t.Fatal("expected a single possibility to be resolved")
			}
			if ss.board[testcase.expected[0]][testcase.expected[1]] != 3 {
				t.Fatalf("expected cell %v to contain 3", testcase.expected)
			}
		})
	}
}

func TestFillSinglePossibilityRejectsInvalidOrCompletedUnits(t *testing.T) {
	ss := newSolverForCandidates(nil)
	ss.completed.row[0] = true
	ss.completed.col[0] = true
	ss.completed.box[0][0] = true

	testcases := []struct {
		name    string
		resolve func() (bool, error)
	}{
		{
			name:    "invalid row",
			resolve: func() (bool, error) { return ss.fillSinglePossibilityInRow(-1) },
		},
		{
			name:    "completed row",
			resolve: func() (bool, error) { return ss.fillSinglePossibilityInRow(0) },
		},
		{
			name:    "invalid column",
			resolve: func() (bool, error) { return ss.fillSinglePossibilityInCol(9) },
		},
		{
			name:    "completed column",
			resolve: func() (bool, error) { return ss.fillSinglePossibilityInCol(0) },
		},
		{
			name:    "invalid box",
			resolve: func() (bool, error) { return ss.fillSinglePossibilityInBox(3, 0) },
		},
		{
			name:    "completed box",
			resolve: func() (bool, error) { return ss.fillSinglePossibilityInBox(0, 0) },
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			changed, err := testcase.resolve()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if changed {
				t.Fatal("expected no change")
			}
		})
	}
}

func TestPlaceValidatesCellAndCandidate(t *testing.T) {
	ss := newSolverForCandidates(map[[2]int][]int{
		{0, 0}: {1},
	})
	ss.board[0][1] = 2

	testcases := []struct {
		name string
		row  int
		col  int
		val  int
	}{
		{name: "invalid coordinates", row: -1, col: 0, val: 1},
		{name: "filled cell", row: 0, col: 1, val: 2},
		{name: "impossible candidate", row: 0, col: 0, val: 2},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if err := ss.place(testcase.row, testcase.col, testcase.val); !errors.Is(err, ErrUnexpected) {
				t.Fatalf("expected ErrUnexpected, got %v", err)
			}
		})
	}
}

func TestResolveCertainCellRequiresOneCandidate(t *testing.T) {
	ss := newSolverForCandidates(map[[2]int][]int{
		{0, 0}: {1, 2},
	})

	if err := ss.resolveCertainCell(0, 0); !errors.Is(err, ErrUncertainCell) {
		t.Fatalf("expected ErrUncertainCell, got %v", err)
	}
}

func TestUpdateLinkedRequiresFilledCell(t *testing.T) {
	ss := newSolverForCandidates(nil)

	if err := ss.updateLinked(0, 0); !errors.Is(err, ErrEmptyCell) {
		t.Fatalf("expected ErrEmptyCell, got %v", err)
	}
}

func TestResolveSpearsSkipsCompletedBox(t *testing.T) {
	ss := newSolverForCandidates(map[[2]int][]int{
		{0, 0}: {5},
		{0, 3}: {5},
	})
	ss.completed.box[0][0] = true

	if ss.resolveSpearsInBox(0, 0) {
		t.Fatal("expected completed box to be skipped")
	}
}

func TestResolveDoubleSpearsRejectsInvalidUnit(t *testing.T) {
	ss := newSolverForCandidates(nil)

	if ss.resolveDoubleSpearsInBoxRow(-1) {
		t.Fatal("expected invalid box row to be rejected")
	}
	if ss.resolveDoubleSpearsInBoxCol(3) {
		t.Fatal("expected invalid box column to be rejected")
	}
}

func TestSmartSolveReportedPuzzle(t *testing.T) {
	puzzle := board.Board{
		{0, 8, 0, 0, 1, 2, 9, 4, 3},
		{2, 0, 9, 7, 0, 0, 8, 1, 0},
		{6, 0, 0, 0, 3, 0, 5, 0, 2},
		{0, 0, 3, 4, 0, 1, 7, 0, 0},
		{0, 0, 0, 0, 0, 5, 4, 0, 8},
		{0, 0, 7, 0, 0, 0, 0, 9, 0},
		{0, 0, 0, 0, 0, 0, 6, 0, 0},
		{4, 0, 0, 5, 2, 0, 0, 8, 0},
		{9, 5, 8, 3, 0, 6, 1, 2, 0},
	}
	ss, err := NewSudokuSolver(&puzzle)
	if err != nil {
		t.Fatal(err)
	}
	if err := ss.Solve(); err != nil {
		t.Fatal(err)
	}
	expected := board.Board{
		{7, 8, 5, 6, 1, 2, 9, 4, 3},
		{2, 3, 9, 7, 5, 4, 8, 1, 6},
		{6, 1, 4, 9, 3, 8, 5, 7, 2},
		{8, 2, 3, 4, 9, 1, 7, 6, 5},
		{1, 9, 6, 2, 7, 5, 4, 3, 8},
		{5, 4, 7, 8, 6, 3, 2, 9, 1},
		{3, 7, 2, 1, 8, 9, 6, 5, 4},
		{4, 6, 1, 5, 2, 7, 3, 8, 9},
		{9, 5, 8, 3, 4, 6, 1, 2, 7},
	}
	if puzzle != expected {
		t.Errorf("smart solve produced %v, want %v", puzzle, expected)
	}
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
		pq:       priorityqueue.NewPriorityQueue(func(a, b queuedCell) bool { return a.possibilityCount < b.possibilityCount }),
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
