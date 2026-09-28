package boardstate

import (
	"errors"

	"github.com/michaelzhan1/sudoku2/internal/boardstate/solve"
)

var ErrSolveFailed = errors.New("solve failed")

// BruteForceSolve solves the board using DFS
func (b *BoardState) BruteForceSolve() error {
	b.Reset()
	b.remaining = 0
	err := solve.SolveDFS(&b.board)
	if err != nil {
		return err
	}

	if !b.IsComplete() {
		return ErrSolveFailed
	}
	return nil
}

func (b *BoardState) SmartSolve() error {
	b.Reset()
	b.remaining = 0
	solver := solve.NewSudokuSolver(&b.board)
	err := solver.Solve()
	if err != nil {
		return err
	}

	// TODO: reenable this check
	// if !b.IsComplete() {
	// 	return ErrSolveFailed
	// }
	return nil
}
