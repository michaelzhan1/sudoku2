package cli

import (
	"fmt"
	"strconv"

	"github.com/michaelzhan1/sudoku2/internal/utils"
)

// clearScreen clears the terminal
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ParseRowCol parses row and column strings into 0-indexed integers.
func ParseRowCol(rowStr, colStr string) (row, col int, ok bool) {
	r, err1 := strconv.Atoi(rowStr)
	c, err2 := strconv.Atoi(colStr)
	if err1 != nil || err2 != nil || !utils.CheckBounds(r-1, c-1) {
		return 0, 0, false
	}
	return r - 1, c - 1, true
}

// ParseRowColVal parses row, column, and value strings
func ParseRowColVal(rowStr, colStr, valStr string) (row, col, val int, ok bool) {
	r, err1 := strconv.Atoi(rowStr)
	c, err2 := strconv.Atoi(colStr)
	v, err3 := strconv.Atoi(valStr)
	if err1 != nil || err2 != nil || err3 != nil || !utils.CheckBounds(r-1, c-1) || !utils.CheckValue(v) {
		return 0, 0, 0, false
	}
	return r - 1, c - 1, v, true
}
