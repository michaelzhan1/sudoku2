package board

// BruteForceSolve solves the board using DFS
func (b *BoardState) BruteForceSolve() bool {
	b.Reset()
	res := solveDFS(&b.board, 0)
	if res {
		b.remaining = 0
	}
	return res
}

// solveDFS is a helper function for BruteForceSolve that contains the DFS logic
func solveDFS(board *Board, i int) bool {
	if i >= 81 {
		return true
	}

	row := i / 9
	col := i % 9

	if board[row][col] != 0 {
		return solveDFS(board, i+1)
	}

	for num := 1; num <= 9; num++ {
		board[row][col] = num
		if board.IsValidPlacement(row, col, num) {
			if solveDFS(board, i+1) {
				return true
			}
		}
		board[row][col] = 0
	}

	return false
}
