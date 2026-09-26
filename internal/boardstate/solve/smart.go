package solve

import (
	"errors"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/utils/fastset"
	"github.com/michaelzhan1/sudoku2/internal/utils/priorityqueue"
)

var ErrEmptyCell = errors.New("cell is empty, cannot update possible values")
var ErrUncertainCell = errors.New("cell is uncertain, cannot fully resolve")
var ErrUnresolvableCell = errors.New("cel found with zero possible values, sudoku is unsolvable")
var ErrInfiniteLoop = errors.New("infinite loop detected, sudoku is unsolvable")

type SudokuSolver struct {
	board    board.Board
	possible [9][9]*fastset.FastSet[int]
	pq       *priorityqueue.PriorityQueue[[2]int]
}

func NewSudokuSolver(b board.Board) *SudokuSolver {
	ss := &SudokuSolver{
		board:    b,
		possible: [9][9]*fastset.FastSet[int]{},
		pq:       nil,
	}

	// assign priority queue later to avoid circular dependency on ss.possible
	ss.pq = priorityqueue.NewPriorityQueue[[2]int](func(a, b [2]int) bool {
		rowA, colA := a[0], a[1]
		rowB, colB := b[0], b[1]
		return ss.possible[rowA][colA].Size() < ss.possible[rowB][colB].Size()
	})

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
			return ErrInfiniteLoop
		}
		cell, _ := ss.pq.Pop() // should not error
		row, col := cell[0], cell[1]

		if ss.board[row][col] != 0 {
			continue
		}

		if ss.possible[row][col].Size() == 0 {
			return ErrUnresolvableCell
		}

		if ss.possible[row][col].Size() == 1 {
			err := ss.resolveCertainCell(row, col)
			if err != nil {
				return err
			}
			changed = true
			ss.updateLinked(row, col)
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
