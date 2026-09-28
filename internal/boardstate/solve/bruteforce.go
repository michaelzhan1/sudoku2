package solve

import "github.com/michaelzhan1/sudoku2/internal/board"

func SolveDFS(b *board.Board) error {
	return solveDFS(b, 0)
}

func solveDFS(b *board.Board, i int) error {
	if i >= 81 {
		return nil
	}

	row := i / 9
	col := i % 9

	if b[row][col] != 0 {
		return solveDFS(b, i+1)
	}

	for num := 1; num <= 9; num++ {
		b[row][col] = num
		if b.IsValidPlacement(row, col, num) {
			if solveDFS(b, i+1) == nil {
				return nil
			}
		}
		b[row][col] = 0
	}

	return ErrUnsolvable
}
