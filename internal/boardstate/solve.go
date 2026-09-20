package boardstate

import "github.com/michaelzhan1/sudoku2/internal/boardstate/solve"

// BruteForceSolve solves the board using DFS
func (b *BoardState) BruteForceSolve() bool {
	b.Reset()
	res := solve.SolveDFS(&b.board, 0)
	if res {
		b.remaining = 0
	}
	return res
}
