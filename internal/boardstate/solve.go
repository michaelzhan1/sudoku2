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

func (b *BoardState) SmartSolve() error {
	b.Reset()
	solver := solve.NewSudokuSolver(b.board)
	err := solver.Solve()
	if err != nil {
		return err
	}

	b.board = solver.Board()
	b.remaining = 0
	return nil
}
