package sudoku

// Solve attempts to solve the board in-place using backtracking.
// Returns true if a solution was found.
func Solve(b *Board) bool {
	return countSolutions(b, 1) == 1
}

// IsComplete returns true when the board is full and all placements are valid.
func IsComplete(b *Board) bool {
	if !b.IsFull() {
		return false
	}
	for row := range 9 {
		for col := range 9 {
			val := b[row][col]
			b[row][col] = 0
			if !b.IsValidPlacement(row, col, val) {
				b[row][col] = val
				return false
			}
			b[row][col] = val
		}
	}
	return true
}
