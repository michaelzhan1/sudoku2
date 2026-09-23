package board_test

import (
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/board"
)

func TestIsValidPlacement(t *testing.T) {
	t.Run("invalid_clues", func(t *testing.T) {
		testcases := []struct {
			name          string
			row, col, val int
		}{
			{"row_below_bounds", -1, 0, 1},
			{"col_below_bounds", 0, -1, 1},
			{"val_below_bounds", 0, 0, -1},
			{"row_above_bounds", 9, 0, 1},
			{"col_above_bounds", 0, 9, 1},
			{"val_above_bounds", 0, 0, 10},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				b := board.Board{}
				out := b.IsValidPlacement(testcase.row, testcase.col, testcase.val)
				if out {
					t.Errorf("Expected IsValidPlacement to return false for input (%d, %d, %d)", testcase.row, testcase.col, testcase.val)
				}
			})
		}
	})

	t.Run("row_check", func(t *testing.T) {
		// assume col1 with val1 exist, and try to place val2 at col2
		testcases := []struct {
			name                   string
			col1, col2, val1, val2 int
			exp                    bool
		}{
			{"invalid_placement", 0, 1, 5, 5, false},
			{"valid_placement", 0, 1, 5, 6, true},
			{"same_cell", 0, 0, 1, 1, true}, // same cell, should be valid
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				b := board.Board{}
				b[0][testcase.col1] = testcase.val1
				out := b.IsValidPlacement(0, testcase.col2, testcase.val2)
				if out != testcase.exp {
					t.Errorf("Expected IsValidPlacement to return %v for input (%d, %d, %d) with existing value %d at col %d", testcase.exp, 0, testcase.col2, testcase.val2, testcase.val1, testcase.col1)
				}
			})
		}
	})

	t.Run("column_check", func(t *testing.T) {
		// assume row1 with val1 exist, and try to place val2 at row2
		testcases := []struct {
			name                   string
			row1, row2, val1, val2 int
			exp                    bool
		}{
			{"invalid_placement", 0, 1, 5, 5, false},
			{"valid placement", 0, 1, 5, 6, true},
			{"same_cell", 0, 0, 1, 1, true}, // same cell, should be valid
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				b := board.Board{}
				b[testcase.row1][0] = testcase.val1
				out := b.IsValidPlacement(testcase.row2, 0, testcase.val2)
				if out != testcase.exp {
					t.Errorf("Expected IsValidPlacement to return %v for input (%d, %d, %d) with existing value %d at row %d", testcase.exp, testcase.row2, 0, testcase.val2, testcase.val1, testcase.row1)
				}
			})
		}
	})
}

func TestIsEmpty(t *testing.T) {
	testcases := []struct {
		name          string
		row, col, val int
		exp           bool
	}{
		{"row_below_bounds", -1, 0, 0, false},
		{"col_below_bounds", 0, -1, 0, false},
		{"row_above_bounds", 9, 8, 0, false},
		{"col_above_bounds", 8, 9, 0, false},
		{"empty_cell", 0, 0, 0, true},
		{"non_empty_cell", 0, 0, 5, false},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			b := board.Board{}
			if testcase.val > 0 {
				b[testcase.row][testcase.col] = testcase.val
			}
			out := b.IsEmpty(testcase.row, testcase.col)
			if out != testcase.exp {
				t.Errorf("Incorrect output from IsEmpty. got=%v, want=%v", out, testcase.exp)
			}
		})
	}
}

func TestGenerateBoard(t *testing.T) {
	rng := rand.New(rand.NewPCG(0, 0))

	t.Run("invalid_clues", func(t *testing.T) {
		testcases := []struct {
			name  string
			clues int
		}{
			{"clues_below_bound", 29},
			{"clues_above_bound", 82},
		}
		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				_, _, err := board.GenerateBoard(testcase.clues, rng)
				if !errors.Is(err, board.ErrInvalidClues) {
					t.Errorf("Expected GenerateBoard to fail with clue input %d", testcase.clues)
				}
			})
		}
	})

	t.Run("valid_clues", func(t *testing.T) {
		testcases := []struct {
			name  string
			clues int
		}{
			{"clues_lower_bound", 30},
			{"clues_upper_bound", 81},
			{"clues_mid_range", 40},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				board, solution, err := board.GenerateBoard(testcase.clues, rng)
				if err != nil {
					t.Errorf("Expected GenerateBoard to succeed with clue input %d, but got error: %v", testcase.clues, err)
				}
				emptyCount := 0
				for i := range 9 {
					for j := range 9 {
						if board[i][j] != 0 && board[i][j] != solution[i][j] {
							t.Errorf("Mismatch between board and solution at (%d, %d): board=%d, solution=%d", i, j, board[i][j], solution[i][j])
						}
						if !solution.IsValidPlacement(i, j, solution[i][j]) {
							t.Errorf("Invalid placement in solution at (%d, %d): value=%d", i, j, solution[i][j])
						}

						if board[i][j] == 0 {
							emptyCount++
						}
					}
				}
				if emptyCount != 81-testcase.clues {
					t.Errorf("Expected %d empty cells in the board, but got %d", 81-testcase.clues, emptyCount)
				}
			})
		}
	})
}
