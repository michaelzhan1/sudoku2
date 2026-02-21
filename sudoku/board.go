package sudoku

import "fmt"

// Board represents a 9x9 Sudoku board with O(1) fullness and conflict tracking.
// The zero value is a valid empty board.
type Board struct {
	cells       [9][9]int
	filledCells int // number of non-zero cells (0–81)
	conflicts   int // number of conflicting peer pairs (same value in same row/col/box)
}

// Get returns the value at (row, col). 0 means the cell is empty.
func (b *Board) Get(row, col int) int {
	return b.cells[row][col]
}

// String returns a human-readable representation of the board.
func (b *Board) String() string {
	result := "  +-------+-------+-------+\n"
	for row := range 9 {
		result += "  | "
		for col := range 9 {
			if b.cells[row][col] == 0 {
				result += ". "
			} else {
				result += fmt.Sprintf("%d ", b.cells[row][col])
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
	return b.cells[row][col] == 0
}

// Set places val at (row, col), updating the fullness and conflict counters.
// Returns false if any argument is out of range.
func (b *Board) Set(row, col, val int) bool {
	if row < 0 || row > 8 || col < 0 || col > 8 || val < 1 || val > 9 {
		return false
	}
	old := b.cells[row][col]
	if old == val {
		return true
	}
	// Remove old value's contribution.
	if old != 0 {
		b.conflicts -= b.countPeers(row, col, old)
		b.filledCells--
	}
	b.cells[row][col] = val
	b.filledCells++
	b.conflicts += b.countPeers(row, col, val)
	return true
}

// Clear removes the value at (row, col), updating the tracking counters.
func (b *Board) Clear(row, col int) {
	old := b.cells[row][col]
	if old == 0 {
		return
	}
	b.conflicts -= b.countPeers(row, col, old)
	b.cells[row][col] = 0
	b.filledCells--
}

// IsValidPlacement checks whether val can be placed at (row, col) without
// violating sudoku rules (ignores the current value at that cell).
func (b *Board) IsValidPlacement(row, col, val int) bool {
	// Check row
	for c := range 9 {
		if c != col && b.cells[row][c] == val {
			return false
		}
	}
	// Check column
	for r := range 9 {
		if r != row && b.cells[r][col] == val {
			return false
		}
	}
	// Check 3x3 box
	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3
	for dr := range 3 {
		for dc := range 3 {
			r, c := boxRow+dr, boxCol+dc
			if (r != row || c != col) && b.cells[r][c] == val {
				return false
			}
		}
	}
	return true
}

// IsFull returns true when there are no empty cells.
func (b *Board) IsFull() bool {
	return b.filledCells == 81
}

// countPeers counts cells that share the same row, column, or 3×3 box with
// (row, col) and contain val. Row/column peers are counted directly; box peers
// exclude cells already in the same row or column to avoid double-counting.
func (b *Board) countPeers(row, col, val int) int {
	n := 0
	for c := range 9 {
		if c != col && b.cells[row][c] == val {
			n++
		}
	}
	for r := range 9 {
		if r != row && b.cells[r][col] == val {
			n++
		}
	}
	boxRow, boxCol := (row/3)*3, (col/3)*3
	for dr := range 3 {
		for dc := range 3 {
			r, c := boxRow+dr, boxCol+dc
			// Exclude cells already counted via the row or column loops.
			if r != row && c != col && b.cells[r][c] == val {
				n++
			}
		}
	}
	return n
}
