package solve

import (
	"errors"
	"fmt"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/utils/fastset"
	"github.com/michaelzhan1/sudoku2/internal/utils/priorityqueue"
)

var ErrEmptyCell = errors.New("cell is empty, cannot update possible values")
var ErrUncertainCell = errors.New("cell is uncertain, cannot fully resolve")
var ErrUnresolvable = errors.New("sudoku is unsolvable")
var ErrUnexpected = errors.New("unexpected error")

type completionStatus struct {
	row [9]bool
	col [9]bool
	box [3][3]bool
}

type SudokuSolver struct {
	board     board.Board
	possible  [9][9]*fastset.FastSet[int]
	pq        *priorityqueue.PriorityQueue[[2]int]
	completed completionStatus
}

func NewSudokuSolver(b board.Board) *SudokuSolver {
	ss := &SudokuSolver{
		board:     b,
		possible:  [9][9]*fastset.FastSet[int]{},
		pq:        nil,
		completed: completionStatus{},
	}

	// assign priority queue later to avoid circular dependency on ss.possible
	ss.pq = priorityqueue.NewPriorityQueue[[2]int](func(a, b [2]int) bool {
		rowA, colA := a[0], a[1]
		rowB, colB := b[0], b[1]
		return ss.possible[rowA][colA].Size() < ss.possible[rowB][colB].Size()
	})

	// check completion status
	for i := range 9 {
		rowComplete := true
		colComplete := true
		boxComplete := true

		for j := range 9 {
			// row
			if ss.board[i][j] == 0 {
				rowComplete = false
			}
			// col
			if ss.board[j][i] == 0 {
				colComplete = false
			}
			// box
			boxRow := (i / 3) * 3
			boxCol := (i % 3) * 3
			if ss.board[boxRow+(j/3)][boxCol+(j%3)] == 0 {
				boxComplete = false
			}
		}
		ss.completed.row[i] = rowComplete
		ss.completed.col[i] = colComplete
		ss.completed.box[i/3][i%3] = boxComplete
	}

	for i := range 9 {
		for j := range 9 {
			ss.possible[i][j] = fastset.NewFastSet[int]()

			if ss.board[i][j] == 0 {
				for k := 1; k <= 9; k++ {
					ss.possible[i][j].Add(k)
				}
				ss.pq.Push([2]int{i, j})
			}
		}
	}

	// remove impossible values
	for i := range 9 {
		for j := range 9 {
			if ss.board[i][j] == 0 {
				continue
			}

			ss.updateLinked(i, j)
		}
	}

	return ss
}

func (ss *SudokuSolver) Board() board.Board {
	return ss.board
}

func (ss *SudokuSolver) Solve() error {
	changed := true
	for !ss.pq.IsEmpty() {
		if !changed {
			// return fmt.Errorf("%w: infinite loop, no more progress can be made", ErrUnresolvable)
			return nil
		}
		changed = false

		// resolve all single-possibility cells
		for top, ok := ss.pq.Peek(); ok && ss.possible[top[0]][top[1]].Size() <= 1; top, ok = ss.pq.Peek() {
			cell, ok := ss.pq.Pop()
			if !ok {
				// should not happen
				return fmt.Errorf("%w: priority queue is empty", ErrUnexpected)
			}
			row, col := cell[0], cell[1]

			if ss.board[row][col] != 0 {
				continue
			}

			if ss.possible[row][col].Size() == 0 {
				return fmt.Errorf("%w: cell (%d, %d) has no possible values", ErrUnresolvable, row, col)
			}

			if ss.possible[row][col].Size() == 1 {
				err := ss.resolveCertainCell(row, col)
				if err != nil {
					return err
				}
				changed = true
				ss.updateLinked(row, col)
				continue
			}
		}
	}

	return nil
}

func (ss *SudokuSolver) resolveCertainCell(row, col int) error {
	if ss.possible[row][col].Size() != 1 {
		return ErrUncertainCell
	}

	ss.board[row][col], _ = ss.possible[row][col].Peek()
	err := ss.updateLinked(row, col)
	if err != nil {
		ss.board[row][col] = 0
		return err
	}

	ss.possible[row][col].Clear()
	return nil
}

// updateLinked updates the possible values for all cells based on (row, col).
// It does not update the cell (row, col) itself.
func (ss *SudokuSolver) updateLinked(row, col int) error {
	val := ss.board[row][col]
	if val == 0 {
		return ErrEmptyCell
	}

	for k := range 9 {
		if k != row && ss.board[k][col] == 0 {
			ss.possible[k][col].Remove(val)
			ss.pq.Push([2]int{k, col})
		}
		if k != col && ss.board[row][k] == 0 {
			ss.possible[row][k].Remove(val)
			ss.pq.Push([2]int{row, k})
		}
	}

	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3
	for r := boxRow; r < boxRow+3; r++ {
		for c := boxCol; c < boxCol+3; c++ {
			if (r != row || c != col) && ss.board[r][c] == 0 {
				ss.possible[r][c].Remove(val)
				ss.pq.Push([2]int{r, c})
			}
		}
	}

	return nil
}
