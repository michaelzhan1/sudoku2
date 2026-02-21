package sudoku

// Solve attempts to solve the board in-place using backtracking.
// Returns true if a solution was found.
func Solve(b *Board) bool {
	return countSolutions(b, 1) == 1
}

// IsComplete returns true when the board is full and has no conflicting cells.
// Runs in O(1) using the board's tracking fields.
func IsComplete(b *Board) bool {
	return b.IsFull() && b.conflicts == 0
}
