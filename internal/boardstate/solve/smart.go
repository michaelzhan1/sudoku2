package solve

import (
	"errors"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/fastset"
)

var ErrEmptyCell = errors.New("cell is empty, cannot update possible values")

type SudokuSolver struct {
	board     board.Board
	possible  [9][9]*fastset.FastSet[int]
	remaining int
}

func NewSudokuSolver(b board.Board) *SudokuSolver {
	ss := &SudokuSolver{
		board:     b,
		possible:  [9][9]*fastset.FastSet[int]{},
		remaining: 0,
	}

	for i := range 9 {
		for j := range 9 {
			ss.possible[i][j] = fastset.NewFastSet[int]()

			if ss.board[i][j] == 0 {
				ss.remaining++
				for k := 1; k <= 9; k++ {
					ss.possible[i][j].Add(k)
				}
			}
		}
	}

	// remove impossible values
	for i := range 9 {
		for j := range 9 {
			if ss.board[i][j] == 0 {
				continue
			}

			ss.updatePossible(i, j, ss.board[i][j])
		}
	}

	return ss
}

func (ss *SudokuSolver) Solve() bool {

}

func (ss *SudokuSolver) resolveCertainCell(row, col int) bool {
	if ss.possible[row][col].Size() != 1 {
		return false
	}

	ss.board[row][col], _ = ss.possible[row][col].Peek()
	ss.remaining--
	ss.possible[row][col].Clear()
	ss.updatePossible(row, col, ss.board[row][col])
	return true
}

// updatePossible updates the possible values for all cells after placing val in (row, col).
// It does not update the cell (row, col) itself.
func (ss *SudokuSolver) updatePossible(row, col int) error {
	val := ss.board[row][col]
	if val == 0 {
		return ErrEmptyCell
	}

	for k := range 9 {
		if k != row {
			ss.possible[k][col].Remove(val)
		}
		if k != col {
			ss.possible[row][k].Remove(val)
		}
	}

	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3
	for r := boxRow; r < boxRow+3; r++ {
		for c := boxCol; c < boxCol+3; c++ {
			if r != row || c != col {
				ss.possible[r][c].Remove(val)
			}
		}
	}

	return nil
}
