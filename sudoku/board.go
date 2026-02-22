package sudoku

import "fmt"

// Board represents a 9x9 sudoku board.
// 0 means the cell is empty.
type Board [9][9]int

// String returns a human-readable representation of the board.
func (b *Board) String() string {
	result := "  +-------+-------+-------+\n"
	for row := range 9 {
		result += "  | "
		for col := range 9 {
			if b[row][col] == 0 {
				result += ". "
			} else {
				result += fmt.Sprintf("%d ", b[row][col])
			}
			if col == 2 || col == 5 {
				result += "| "
			}
		}
		result += "|\n"
		if row == 2 || row == 5 {
			result += "  +-------+-------+-------+\n"
		}
	}
	result += "  +-------+-------+-------+"
	return result
}

// IsEmpty returns true if the cell at (row, col) is empty.
func (b *Board) IsEmpty(row, col int) bool {
	return b[row][col] == 0
}

// Set places a value at (row, col). Returns false if out of range.
func (b *Board) Set(row, col, val int) bool {
	if row < 0 || row > 8 || col < 0 || col > 8 || val < 1 || val > 9 {
		return false
	}
	b[row][col] = val
	return true
}

// Clear removes the value at (row, col).
func (b *Board) Clear(row, col int) {
	b[row][col] = 0
}

// IsValidPlacement checks whether val can be placed at (row, col) without
// violating sudoku rules (ignores the current value at that cell).
func (b *Board) IsValidPlacement(row, col, val int) bool {
	// Check row
	for c := range 9 {
		if c != col && b[row][c] == val {
			return false
		}
	}
	// Check column
	for r := range 9 {
		if r != row && b[r][col] == val {
			return false
		}
	}
	// Check 3x3 box
	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3
	for dr := range 3 {
		for dc := range 3 {
			r, c := boxRow+dr, boxCol+dc
			if (r != row || c != col) && b[r][c] == val {
				return false
			}
		}
	}
	return true
}

// IsFull returns true when there are no empty cells.
func (b *Board) IsFull() bool {
	for row := range 9 {
		for col := range 9 {
			if b[row][col] == 0 {
				return false
			}
		}
	}
	return true
}
