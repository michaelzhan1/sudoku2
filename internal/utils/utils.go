package utils

import "errors"

var ErrInvalidClues = errors.New("invalid number of clues; must be between 30 and 81")

// CheckBounds checks if the given row and column are within the bounds of a 9x9 sudoku board
func CheckBounds(row, col int) bool {
	return row >= 0 && row < 9 && col >= 0 && col < 9
}

// CheckValue checks if the given value is a valid sudoku number (1-9)
func CheckValue(val int) bool {
	return val >= 1 && val <= 9
}
