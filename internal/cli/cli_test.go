package cli_test

import (
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/cli"
)

func TestClearScreen(t *testing.T) {
	// just make sure it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ClearScreen panicked: %v", r)
		}
	}()
	cli.ClearScreen()
}

func TestParseRowCol(t *testing.T) {
	testcases := []struct {
		name   string
		rowStr string
		colStr string
		expRow int
		expCol int
		expOk  bool
	}{
		{"valid_input", "1", "2", 0, 1, true},
		{"invalid_row", "a", "2", 0, 0, false},
		{"invalid_col", "1", "b", 0, 0, false},
		{"out_of_bounds_row", "10", "2", 0, 0, false},
		{"out_of_bounds_col", "1", "10", 0, 0, false},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			row, col, ok := cli.ParseRowCol(testcase.rowStr, testcase.colStr)
			if row != testcase.expRow || col != testcase.expCol || ok != testcase.expOk {
				t.Errorf("Incorrect output from ParseRowCol. got=(%d, %d, %v), want=(%d, %d, %v)", row, col, ok, testcase.expRow, testcase.expCol, testcase.expOk)
			}
		})
	}
}

func TestParseRowColVal(t *testing.T) {
	testcases := []struct {
		name   string
		rowStr string
		colStr string
		valStr string
		expRow int
		expCol int
		expVal int
		expOk  bool
	}{
		{"valid_input", "1", "2", "3", 0, 1, 3, true},
		{"invalid_row", "a", "2", "3", 0, 0, 0, false},
		{"invalid_col", "1", "b", "3", 0, 0, 0, false},
		{"invalid_val", "1", "2", "c", 0, 0, 0, false},
		{"out_of_bounds_row", "10", "2", "3", 0, 0, 0, false},
		{"out_of_bounds_col", "1", "10", "3", 0, 0, 0, false},
		{"out_of_bounds_val", "1", "2", "10", 0, 0, 0, false},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			row, col, val, ok := cli.ParseRowColVal(testcase.rowStr, testcase.colStr, testcase.valStr)
			if row != testcase.expRow || col != testcase.expCol || val != testcase.expVal || ok != testcase.expOk {
				t.Errorf("Incorrect output from ParseRowColVal. got=(%d, %d, %d, %v), want=(%d, %d, %d, %v)", row, col, val, ok, testcase.expRow, testcase.expCol, testcase.expVal, testcase.expOk)
			}
		})
	}
}
