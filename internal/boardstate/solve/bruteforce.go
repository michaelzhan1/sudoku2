package solve

import "github.com/michaelzhan1/sudoku2/internal/board"

// SolveDFS is a helper function for BruteForceSolve that contains the DFS logic
func SolveDFS(b *board.Board, i int) bool {
	if i >= 81 {
		return true
	}

	row := i / 9
	col := i % 9

	if b[row][col] != 0 {
		return SolveDFS(b, i+1)
	}

	for num := 1; num <= 9; num++ {
		b[row][col] = num
		if b.IsValidPlacement(row, col, num) {
			if SolveDFS(b, i+1) {
				return true
			}
		}
		b[row][col] = 0
	}

	return false
}
