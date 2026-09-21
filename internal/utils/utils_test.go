package utils_test

import (
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/utils"
)

func TestCheckBounds(t *testing.T) {
	testcases := []struct {
		name string
		row  int
		col  int
		exp  bool
	}{
		{"row_below_bounds", -1, 0, false},
		{"col_below_bounds", 0, -1, false},
		{"row_above_bounds", 9, 8, false},
		{"col_above_bounds", 8, 9, false},
		{"within_bounds_low", 0, 0, true},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			out := utils.CheckBounds(testcase.row, testcase.col)
			if out != testcase.exp {
				t.Errorf("Incorrect output from CheckBounds. got=%v, want=%v", out, testcase.exp)
			}
		})
	}
}

func TestCheckValue(t *testing.T) {
	testcases := []struct {
		name  string
		value int
		exp   bool
	}{
		{"value_below_bounds", 0, false},
		{"value_above_bounds", 10, false},
		{"within_bounds", 2, true},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			out := utils.CheckValue(testcase.value)
			if out != testcase.exp {
				t.Errorf("Incorrect output from CheckValue. got=%v, want=%v", out, testcase.exp)
			}
		})
	}
}
