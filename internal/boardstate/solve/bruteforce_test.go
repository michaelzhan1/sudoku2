package solve_test

import (
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/boardstate/solve"
)

func TestSolveDFS(t *testing.T) {
	board := board.Board{
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
		{4, 5, 6, 7, 8, 9, 1, 2, 3},
		{7, 8, 9, 1, 2, 3, 4, 5, 6},
		{2, 3, 4, 5, 6, 7, 8, 9, 1},
		{5, 6, 7, 8, 9, 1, 2, 3, 4},
		{8, 9, 1, 2, 3, 4, 5, 6, 7},
		{3, 4, 5, 6, 7, 8, 9, 1, 2},
		{6, 7, 8, 9, 1, 2, 3, 4, 5},
		{9, 1, 2, 3, 4, 5, 6, 7, 8},
	}

	testcases := []struct {
		name string
		sets []struct {
			row, col, val int
		}
		exp bool
	}{
		{
			name: "full_board",
			exp:  true,
		},
		{
			name: "one_empty_cell",
			sets: []struct{ row, col, val int }{{0, 0, 0}},
			exp:  true,
		},
		{
			name: "two_empty_cells",
			sets: []struct{ row, col, val int }{{0, 0, 0}, {1, 1, 0}},
			exp:  true,
		},
		{
			name: "impossible_solve",
			sets: []struct {
				row, col, val int
			}{
				{0, 0, 0},
				{0, 1, 1},
			},
			exp: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testboard := board
			for _, set := range testcase.sets {
				testboard[set.row][set.col] = set.val
			}
			res := solve.SolveDFS(&testboard, 0)
			if res != testcase.exp {
				t.Errorf("expected SolveDFS to return %v, got %v", testcase.exp, res)
			}
		})
	}
}
