package boardstate_test

import (
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/boardstate"
)

func TestBruteForceSolve(t *testing.T) {
	solution := solutionFixture

	testcases := []struct {
		name     string
		original board.Board
	}{
		{
			name: "trivial board",
			original: board.Board{
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{4, 5, 6, 7, 8, 9, 1, 2, 3},
				{7, 8, 9, 1, 2, 3, 4, 5, 6},
				{2, 3, 4, 5, 6, 7, 8, 9, 1},
				{5, 6, 7, 8, 9, 1, 2, 3, 4},
				{8, 9, 1, 2, 3, 4, 5, 6, 7},
				{3, 4, 5, 6, 7, 8, 9, 1, 2},
				{6, 7, 8, 9, 1, 2, 3, 4, 5},
				{9, 1, 2, 3, 4, 5, 6, 7, 0},
			},
		},
		{
			name: "slightly harder board",
			original: board.Board{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 5, 6, 7, 8, 9, 1, 2, 3},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 3, 4, 5, 6, 7, 8, 9, 1},
				{0, 6, 7, 0, 9, 1, 2, 3, 4},
				{0, 9, 0, 0, 3, 4, 5, 6, 7},
				{0, 4, 5, 6, 7, 8, 9, 0, 2},
				{0, 7, 8, 9, 1, 0, 3, 4, 5},
				{0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
		},
		{
			name:     "completed board",
			original: solution,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			bs, err := boardstate.NewBoardState(testcase.original, solution)
			if err != nil {
				t.Errorf("NewBoardState returned error for valid board: %v", err)
			}

			solved := bs.BruteForceSolve()
			if !solved {
				t.Errorf("BruteForceSolve returned false for solvable board")
			}
		})
	}
}
