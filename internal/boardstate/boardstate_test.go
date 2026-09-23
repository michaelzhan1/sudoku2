package boardstate_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/boardstate"
)

// assigning x := solutionFixture is a copy due to [9][9]int being a value type
var solutionFixture = board.Board{
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

func TestNewBoardState(t *testing.T) {
	t.Run("valid_board", func(t *testing.T) {
		solution := solutionFixture

		testcases := []struct {
			name         string
			board        board.Board
			expRemaining int
		}{
			{
				name: "one_clue",
				board: board.Board{
					{1, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				expRemaining: 80,
			},
			{
				name:         "empty_board",
				board:        board.Board{},
				expRemaining: 81,
			},
			{
				name:         "full_board",
				board:        solution,
				expRemaining: 0,
			},
			{
				name: "multiple_clues",
				board: board.Board{
					{1, 2, 3, 0, 0, 0, 0, 0, 0},
					{4, 5, 6, 0, 0, 0, 0, 0, 0},
					{7, 8, 9, 0, 0, 0, 0, 0, 0},
					{2, 3, 4, 0, 0, 0, 0, 0, 0},
					{5, 6, 7, 0, 0, 0, 0, 0, 0},
					{8, 9, 1, 0, 0, 0, 0, 0, 0},
					{3, 4, 5, 0, 0, 0, 0, 0, 0},
					{6, 7, 8, 0, 0, 0, 0, 0, 0},
					{9, 1, 2, 0, 0, 0, 0, 0, 0},
				},
				expRemaining: 54,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(testcase.board, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}
				if bs == nil {
					t.Errorf("NewBoardState returned nil for valid board")
				}

				if bs.Remaining() != testcase.expRemaining {
					t.Errorf("NewBoardState did not set remaining correctly, got %d, want %d", bs.Remaining(), testcase.expRemaining)
				}
				if bs.Board() != testcase.board {
					t.Errorf("NewBoardState did not set board correctly")
				}
				if bs.Original() != testcase.board {
					t.Errorf("NewBoardState did not set original correctly")
				}
				if bs.Solution() != solution {
					t.Errorf("NewBoardState did not set solution correctly")
				}
			})
		}
	})

	t.Run("invalid_solution", func(t *testing.T) {
		testcases := []struct {
			name  string
			board board.Board
		}{
			{
				name: "zero_in_solution",
				board: board.Board{
					{1, 2, 3, 4, 5, 6, 7, 8, 9},
					{4, 5, 6, 7, 8, 9, 1, 2, 3},
					{7, 8, 9, 1, 2, 3, 4, 5, 6},
					{2, 3, 4, 5, 6, 7, 8, 9, 1},
					{5, 6, 7, 8, 9, 1, 2, 3, 4},
					{8, 9, 1, 2, 3, 4, 5, 6, 7},
					{3, 4, 5, 6, 7, 8, 9, 1, 2},
					{6, 7, 8, 9, 1, 2, 3, 4, 5},
					{9, 1, 2, 3, 4, 5, 6, 7, 0}, // zero
				},
			},
			{
				name: "invalid_value_in_solution",
				board: board.Board{
					{1, 2, 3, 4, 5, 6, 7, 8, 9},
					{4, 5, 6, 7, 8, 9, 1, 2, 3},
					{7, 8, 9, 1, 2, 3, 4, 5, 6},
					{2, 3, 4, 5, 6, 7, 8, 9, 1},
					{5, 6, 7, 8, 9, 1, 2, 3, 4},
					{8, 9, 1, 2, 3, 4, 5, 6, -1}, // invalid value
					{3, 4, 5, 6, 7, 8, 9, 1, 2},
					{6, 7, 8, 9, 1, 2, 3, 4, 5},
					{9, 1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			{
				name: "invalid_placement_in_solution",
				board: board.Board{
					{1, 2, 3, 4, 5, 6, 7, 8, 9},
					{4, 5, 6, 7, 8, 9, 1, 2, 3},
					{7, 8, 9, 1, 2, 3, 4, 5, 6},
					{2, 3, 4, 5, 6, 7, 8, 9, 1},
					{5, 6, 7, 8, 9, 1, 2, 3, 4},
					{8, 9, 1, 2, 3, 4, 5, 6, 7},
					{3, 4, 5, 6, 7, 8, 9, 1, 2},
					{6, 7, 8, 9, 1, 2, 3, 4, 5},
					{9, 1, 2, 3, 4, 5, 6, 7, 1}, // duplicate 1
				},
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(board.Board{}, testcase.board)
				if !errors.Is(err, boardstate.ErrInvalidSolution) {
					t.Errorf("NewBoardState did not return ErrInvalidSolution for %s, got %v", testcase.name, err)
				}
				if bs != nil {
					t.Errorf("NewBoardState returned non-nil for %s", testcase.name)
				}
			})
		}
	})

	t.Run("invalid_original", func(t *testing.T) {
		solution := solutionFixture

		testcases := []struct {
			name  string
			board board.Board
		}{
			{
				name: "invalid_in_bounds_value",
				board: board.Board{
					{1, 3, 0, 0, 0, 0, 0, 0, 0}, // still 1-9, but mismatch
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
			},
			{
				name: "invalid_out_of_bounds_value",
				board: board.Board{
					{1, -1, 0, 0, 0, 0, 0, 0, 0}, // invalid value
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(testcase.board, solution)
				if !errors.Is(err, boardstate.ErrInvalidOriginal) {
					t.Errorf("NewBoardState did not return ErrInvalidOriginal for %s, got %v", testcase.name, err)
				}
				if bs != nil {
					t.Errorf("NewBoardState returned non-nil for %s", testcase.name)
				}
			})
		}
	})
}

func TestReset(t *testing.T) {
	solution := solutionFixture

	testcases := []struct {
		name             string
		original         board.Board
		row, col, value  int
		initialRemaining int
	}{
		{
			name:             "empty_board",
			row:              0,
			col:              0,
			value:            1,
			initialRemaining: 81,
		},
		{
			name: "one_clue",
			original: board.Board{
				{1, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			row:              0,
			col:              1,
			value:            2,
			initialRemaining: 80,
		},
		{
			name: "multiple_clues",
			original: board.Board{
				{1, 2, 3, 0, 0, 0, 0, 0, 0},
				{4, 5, 6, 0, 0, 0, 0, 0, 0},
				{7, 8, 9, 0, 0, 0, 0, 0, 0},
				{2, 3, 4, 0, 0, 0, 0, 0, 0},
				{5, 6, 7, 0, 0, 0, 0, 0, 0},
				{8, 9, 1, 0, 0, 0, 0, 0, 0},
				{3, 4, 5, 0, 0, 0, 0, 0, 0},
				{6, 7, 8, 0, 0, 0, 0, 0, 0},
				{9, 1, 2, 0, 0, 0, 0, 0, 0},
			},
			row:              0,
			col:              3,
			value:            4,
			initialRemaining: 54,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			bs, err := boardstate.NewBoardState(testcase.original, solution)
			if err != nil {
				t.Errorf("NewBoardState returned error for valid board: %v", err)
			}
			if bs.Remaining() != testcase.initialRemaining {
				t.Errorf("NewBoardState did not set remaining correctly, got %d, want %d", bs.Remaining(), testcase.initialRemaining)
			}

			bs.Set(testcase.row, testcase.col, testcase.value) // set a value
			if bs.Remaining() != testcase.initialRemaining-1 {
				t.Errorf("Set did not update remaining correctly, got %d, want %d", bs.Remaining(), testcase.initialRemaining-1)
			}

			bs.Reset()
			if bs.Board() != testcase.original {
				t.Errorf("Reset did not reset board correctly")
			}
			if bs.Remaining() != testcase.initialRemaining {
				t.Errorf("Reset did not reset remaining correctly, got %d, want %d", bs.Remaining(), testcase.initialRemaining)
			}
		})
	}
}

func TestString(t *testing.T) {
	solution := board.Board{
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
		name  string
		board board.Board
		sets  []struct {
			row, col, val int
		}
		exp string
	}{
		{
			name:  "empty_board",
			board: board.Board{},
			exp: "    1 2 3   4 5 6   7 8 9\n" +
				"  +-------+-------+-------+\n" +
				"1 | . . . | . . . | . . . |\n" +
				"2 | . . . | . . . | . . . |\n" +
				"3 | . . . | . . . | . . . |\n" +
				"  +-------+-------+-------+\n" +
				"4 | . . . | . . . | . . . |\n" +
				"5 | . . . | . . . | . . . |\n" +
				"6 | . . . | . . . | . . . |\n" +
				"  +-------+-------+-------+\n" +
				"7 | . . . | . . . | . . . |\n" +
				"8 | . . . | . . . | . . . |\n" +
				"9 | . . . | . . . | . . . |\n" +
				"  +-------+-------+-------+\n",
		},
		{
			name:  "full_board",
			board: solution,
			exp: "    1 2 3   4 5 6   7 8 9\n" +
				"  +-------+-------+-------+\n" +
				"1 | 1 2 3 | 4 5 6 | 7 8 9 |\n" +
				"2 | 4 5 6 | 7 8 9 | 1 2 3 |\n" +
				"3 | 7 8 9 | 1 2 3 | 4 5 6 |\n" +
				"  +-------+-------+-------+\n" +
				"4 | 2 3 4 | 5 6 7 | 8 9 1 |\n" +
				"5 | 5 6 7 | 8 9 1 | 2 3 4 |\n" +
				"6 | 8 9 1 | 2 3 4 | 5 6 7 |\n" +
				"  +-------+-------+-------+\n" +
				"7 | 3 4 5 | 6 7 8 | 9 1 2 |\n" +
				"8 | 6 7 8 | 9 1 2 | 3 4 5 |\n" +
				"9 | 9 1 2 | 3 4 5 | 6 7 8 |\n" +
				"  +-------+-------+-------+\n",
		},
		{
			name: "partial_board",
			sets: []struct {
				row, col, val int
			}{
				{0, 0, 1},
			},
			exp: "    1 2 3   4 5 6   7 8 9\n" +
				"  +-------+-------+-------+\n" +
				"1 | 1 . . | . . . | . . . |\n" +
				"2 | . . . | . . . | . . . |\n" +
				"3 | . . . | . . . | . . . |\n" +
				"  +-------+-------+-------+\n" +
				"4 | . . . | . . . | . . . |\n" +
				"5 | . . . | . . . | . . . |\n" +
				"6 | . . . | . . . | . . . |\n" +
				"  +-------+-------+-------+\n" +
				"7 | . . . | . . . | . . . |\n" +
				"8 | . . . | . . . | . . . |\n" +
				"9 | . . . | . . . | . . . |\n" +
				"  +-------+-------+-------+\n",
		},
	}

	// remove all ansi escape sequences
	var ansiEscape = regexp.MustCompile(`\033\[[0-?]*[ -/]*[@-~]`)

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			bs, err := boardstate.NewBoardState(testcase.board, solution)
			if err != nil {
				t.Errorf("NewBoardState returned error for valid board: %v", err)
			}

			for _, set := range testcase.sets {
				err = bs.Set(set.row, set.col, set.val)
				if err != nil {
					t.Errorf("Set returned error for valid parameters: %v", err)
				}
			}

			str := ansiEscape.ReplaceAllString(bs.String(), "")
			if str != testcase.exp {
				t.Errorf("String returned unexpected result, got:\n%s\nwant:\n%s", str, testcase.exp)
			}
		})
	}
}

func TestSet(t *testing.T) {
	solution := solutionFixture

	t.Run("valid_set", func(t *testing.T) {
		testcases := []struct {
			name string
			sets []struct {
				row, col, val, expRemaining int
			}
			expRemaining int
		}{
			{
				name: "set_empty_cell",
				sets: []struct{ row, col, val, expRemaining int }{
					{0, 0, 1, 80},
					{0, 1, 2, 79},
				},
			},
			{
				name: "set_cell_with_clue",
				sets: []struct{ row, col, val, expRemaining int }{
					{0, 0, 1, 80},
					{0, 0, 2, 80},
				},
				expRemaining: 80,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(board.Board{}, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}

				for _, set := range testcase.sets {
					err = bs.Set(set.row, set.col, set.val)
					if err != nil {
						t.Errorf("Set returned error for valid parameters: %v", err)
					}
					if bs.Remaining() != set.expRemaining {
						t.Errorf("Set did not update remaining correctly, got %d, want %d", bs.Remaining(), set.expRemaining)
					}
				}
			})
		}
	})

	t.Run("invalid_set", func(t *testing.T) {
		testcases := []struct {
			name          string
			original      board.Board
			row, col, val int
			expErr        error
		}{
			{
				name:   "out_of_bounds_row",
				row:    9,
				col:    0,
				val:    1,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:   "out_of_bounds_col",
				row:    0,
				col:    9,
				val:    1,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:   "out_of_bounds_val",
				row:    0,
				col:    0,
				val:    10,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:     "locked_value",
				original: board.Board{{1, 0, 0, 0, 0, 0, 0, 0, 0}},
				row:      0,
				col:      0,
				val:      1,
				expErr:   boardstate.ErrLockedValue,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(testcase.original, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}
				err = bs.Set(testcase.row, testcase.col, testcase.val)
				if !errors.Is(err, testcase.expErr) {
					t.Errorf("Set did not return expected error for %s, got %v", testcase.name, err)
				}
			})
		}
	})
}

func TestClear(t *testing.T) {
	solution := solutionFixture

	type action string
	var actionClear action = "clear"
	var actionSet action = "set"

	t.Run("valid_clear", func(t *testing.T) {
		testcases := []struct {
			name    string
			actions []struct {
				action        action
				row, col, val int
				expRemaining  int
			}
		}{
			{
				name: "clear_empty_cell",
				actions: []struct {
					action        action
					row, col, val int
					expRemaining  int
				}{
					{action: actionClear, row: 0, col: 0, val: 0, expRemaining: 81},
				},
			},
			{
				name: "clear_cell_with_value",
				actions: []struct {
					action        action
					row, col, val int
					expRemaining  int
				}{
					{action: actionSet, row: 0, col: 0, val: 1, expRemaining: 80},
					{action: actionClear, row: 0, col: 0, val: 0, expRemaining: 81},
				},
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(board.Board{}, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}

				for _, action := range testcase.actions {
					switch action.action {
					case actionSet:
						err = bs.Set(action.row, action.col, action.val)
						if err != nil {
							t.Errorf("Set returned error for valid parameters: %v", err)
						}
					case actionClear:
						err = bs.Clear(action.row, action.col)
						if err != nil {
							t.Errorf("Clear returned error for valid parameters: %v", err)
						}
					}

					if bs.Remaining() != action.expRemaining {
						t.Errorf("%s did not update remaining correctly, got %d, want %d", action.action, bs.Remaining(), action.expRemaining)
					}
				}
			})
		}
	})

	t.Run("invalid_clear", func(t *testing.T) {
		testcases := []struct {
			name     string
			original board.Board
			row, col int
			expErr   error
		}{
			{
				name:   "out_of_bounds_row",
				row:    9,
				col:    0,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:   "out_of_bounds_col",
				row:    0,
				col:    9,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:     "locked_value",
				original: board.Board{{1, 0, 0, 0, 0, 0, 0, 0, 0}},
				row:      0,
				col:      0,
				expErr:   boardstate.ErrLockedValue,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(testcase.original, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}
				err = bs.Clear(testcase.row, testcase.col)
				if !errors.Is(err, testcase.expErr) {
					t.Errorf("Clear did not return expected error for %s, got %v", testcase.name, err)
				}
			})
		}
	})
}

func TestIsComplete(t *testing.T) {
	solution := solutionFixture

	testcases := []struct {
		name     string
		original board.Board
		sets     []struct {
			row, col, val int
		}
		exp bool
	}{
		{
			name: "empty_board",
			exp:  false,
		},
		{
			name:     "full_board",
			original: solution,
			exp:      true,
		},
		{
			name: "partial_board",
			original: board.Board{
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{4, 5, 6, 7, 8, 9, 1, 2, 3},
				{7, 8, 9, 1, 2, 3, 4, 5, 6},
				{2, 3, 4, 5, 6, 7, 8, 9, 1},
				{5, 6, 7, 8, 9, 1, 2, 3, 4},
				{8, 9, 1, 2, 3, 4, 5, 6, 7},
				{3, 4, 5, 6, 7, 8, 9, 1, 2},
				{6, 7, 8, 9, 1, 2, 3, 4, 5},
				{9, 1, 2, 3, 4, 5, 6, 7, 0}, // one cell empty
			},
			exp: false,
		},
		{
			name: "incorrect_full_board",
			original: board.Board{
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{4, 5, 6, 7, 8, 9, 1, 2, 3},
				{7, 8, 9, 1, 2, 3, 4, 5, 6},
				{2, 3, 4, 5, 6, 7, 8, 9, 1},
				{5, 6, 7, 8, 9, 1, 2, 3, 4},
				{8, 9, 1, 2, 3, 4, 5, 6, 0}, // missing value
				{3, 4, 5, 6, 7, 8, 9, 1, 2},
				{6, 7, 8, 9, 1, 2, 3, 4, 5},
				{9, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			sets: []struct {
				row, col, val int
			}{
				{5, 8, 4}, // fill in the missing value incorrectly
			},
			exp: false,
		},
		{
			name: "correct_full_board",
			original: board.Board{
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{4, 5, 6, 7, 8, 9, 1, 2, 3},
				{7, 8, 9, 1, 2, 3, 4, 5, 6},
				{2, 3, 4, 5, 6, 7, 8, 9, 1},
				{5, 6, 7, 8, 9, 1, 2, 3, 4},
				{8, 9, 1, 2, 3, 4, 5, 6, 0}, // missing value
				{3, 4, 5, 6, 7, 8, 9, 1, 2},
				{6, 7, 8, 9, 1, 2, 3, 4, 5},
				{9, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			sets: []struct {
				row, col, val int
			}{
				{5, 8, 7}, // fill in the missing value correctly
			},
			exp: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			bs, err := boardstate.NewBoardState(testcase.original, solution)
			if err != nil {
				t.Errorf("NewBoardState returned error for valid board: %v", err)
			}

			for _, set := range testcase.sets {
				err = bs.Set(set.row, set.col, set.val)
				if err != nil {
					t.Errorf("Set returned error for valid parameters: %v", err)
				}
			}

			if bs.IsComplete() != testcase.exp {
				t.Errorf("IsComplete returned unexpected result, got %v, want %v", bs.IsComplete(), testcase.exp)
			}
		})
	}
}

func TestGiveHint(t *testing.T) {
	solution := solutionFixture

	t.Run("valid_hint", func(t *testing.T) {
		testcases := []struct {
			name     string
			original board.Board
			row, col int
			expVal   int
		}{
			{
				name:   "empty_board",
				row:    0,
				col:    0,
				expVal: 1,
			},
			{
				name: "partial_board",
				original: board.Board{
					{1, 2, 3, 4, 5, 6, 7, 8, 9},
					{4, 5, 6, 7, 8, 9, 1, 2, 3},
					{7, 8, 9, 1, 2, 3, 4, 5, 6},
					{2, 3, 4, 5, 6, 7, 8, 9, 1},
					{5, 6, 7, 8, 9, 1, 2, 3, 4},
					{8, 9, 1, 2, 3, 4, 5, 6, 0}, // missing value
					{3, 4, 5, 6, 7, 8, 9, 1, 2},
					{6, 7, 8, 9, 1, 2, 3, 4, 5},
					{9, 1, 2, 3, 4, 5, 6, 7, 8},
				},
				row:    5,
				col:    8,
				expVal: 7,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(testcase.original, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}

				val, err := bs.GiveHint(testcase.row, testcase.col)
				if err != nil {
					t.Errorf("GiveHint returned error for valid parameters: %v", err)
				}
				if val != testcase.expVal {
					t.Errorf("GiveHint returned unexpected value, got %d, want %d", val, testcase.expVal)
				}
			})
		}
	})

	t.Run("invalid_hint", func(t *testing.T) {
		testcases := []struct {
			name     string
			original board.Board
			row, col int
			expErr   error
		}{
			{
				name:   "out_of_bounds_row",
				row:    9,
				col:    0,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:   "out_of_bounds_col",
				row:    0,
				col:    9,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:     "locked_cell",
				original: board.Board{{1, 0, 0, 0, 0, 0, 0, 0, 0}},
				row:      0,
				col:      0,
				expErr:   boardstate.ErrLockedValue,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(testcase.original, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}
				_, err = bs.GiveHint(testcase.row, testcase.col)
				if !errors.Is(err, testcase.expErr) {
					t.Errorf("GiveHint did not return expected error for %s, got %v", testcase.name, err)
				}
			})
		}
	})
}

func TestCheck(t *testing.T) {
	solution := solutionFixture

	t.Run("valid_check", func(t *testing.T) {
		testcases := []struct {
			name string
			sets []struct {
				row, col, value int
			}
			row, col int
			expVal   bool
		}{
			{
				name: "correct_value",
				sets: []struct {
					row, col, value int
				}{
					{0, 0, 1},
				},
				row:    0,
				col:    0,
				expVal: true,
			},
			{
				name: "incorrect_value",
				sets: []struct {
					row, col, value int
				}{
					{0, 0, 2},
				},
				row:    0,
				col:    0,
				expVal: false,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(board.Board{}, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}
				for _, set := range testcase.sets {
					err = bs.Set(set.row, set.col, set.value)
					if err != nil {
						t.Errorf("Set returned error for valid parameters: %v", err)
					}
				}

				val, err := bs.Check(testcase.row, testcase.col)
				if err != nil {
					t.Errorf("Check returned error for valid parameters: %v", err)
				}
				if val != testcase.expVal {
					t.Errorf("Check returned unexpected value, got %v, want %v", val, testcase.expVal)
				}
			})
		}
	})

	t.Run("invalid_check", func(t *testing.T) {
		testcases := []struct {
			name     string
			row, col int
			expErr   error
		}{
			{
				name:   "out_of_bounds_row",
				row:    9,
				col:    0,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:   "out_of_bounds_col",
				row:    0,
				col:    9,
				expErr: boardstate.ErrInvalidParameters,
			},
			{
				name:   "empty_cell",
				row:    0,
				col:    0,
				expErr: boardstate.ErrEmptyCell,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.name, func(t *testing.T) {
				bs, err := boardstate.NewBoardState(board.Board{}, solution)
				if err != nil {
					t.Errorf("NewBoardState returned error for valid board: %v", err)
				}
				_, err = bs.Check(testcase.row, testcase.col)
				if !errors.Is(err, testcase.expErr) {
					t.Errorf("Check did not return expected error for %s, got %v", testcase.name, err)
				}
			})
		}
	})
}
