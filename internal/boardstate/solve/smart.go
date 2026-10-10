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
var ErrUnsolvable = errors.New("sudoku is unsolvable")
var ErrUnexpected = errors.New("unexpected error")

type completionStatus struct {
	row [9]bool
	col [9]bool
	box [3][3]bool
}

type SudokuSolver struct {
	board     *board.Board
	possible  [9][9]*fastset.FastSet[int]
	pq        *priorityqueue.PriorityQueue[[2]int]
	completed completionStatus
}

func NewSudokuSolver(b *board.Board) *SudokuSolver {
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

func (ss *SudokuSolver) Solve() error {
	changed := true
	// min heap on number of possible values in a cell
	for !ss.pq.IsEmpty() {
		if !changed {
			// TODO: restore this
			// return fmt.Errorf("%w: infinite loop, no more progress can be made", ErrUnresolvable)
			return nil
		}
		changed = false

		cell, ok := ss.pq.Pop()
		if !ok {
			return fmt.Errorf("%w: priority queue is empty", ErrUnexpected)
		}
		row, col := cell[0], cell[1]
		if ss.board[row][col] != 0 {
			continue
		}
		if ss.possible[row][col].Size() == 0 {
			return fmt.Errorf("%w: cell (%d, %d) has no possible values", ErrUnsolvable, row, col)
		}

		// resolve certain cells
		if ss.possible[row][col].Size() == 1 {
			err := ss.resolveCertainCell(row, col)
			if err != nil {
				return err
			}
			changed = true
			continue
		}

		// check for single-possibility cells in rows, cols, and boxes
		for i := range 9 {
			if !ss.completed.row[i] && ss.fillSinglePossibilityInRow(i) {
				changed = true
				break
			}
			if !ss.completed.col[i] && ss.fillSinglePossibilityInCol(i) {
				changed = true
				break
			}
			if !ss.completed.box[i/3][i%3] && ss.fillSinglePossibilityInBox(i/3, i%3) {
				changed = true
				break
			}
		}
		if changed {
			continue
		}

		// check single spears
		for i := range 9 {
			if ss.resolveSpearsInBox(i/3, i%3) {
				changed = true
				break
			}
		}
		if changed {
			continue
		}

		// use double spears
		for i := range 3 {
			if ss.resolveDoubleSpearsInBoxRow(i) {
				changed = true
				break
			}
			if ss.resolveDoubleSpearsInBoxCol(i) {
				changed = true
				break
			}
		}
		if changed {
			continue
		}

		// use hidden groups
		for i := range 9 {
			if ss.resolveHiddenGroupsInRow(i) {
				changed = true
				break
			}
			if ss.resolveHiddenGroupsInCol(i) {
				changed = true
				break
			}
			if ss.resolveHiddenGroupsInBox(i/3, i%3) {
				changed = true
				break
			}
		}
		if changed {
			continue
		}
	}

	return nil
}

func (ss *SudokuSolver) place(row, col, val int) error {
	if row < 0 || row > 8 || col < 0 || col > 8 || val < 1 || val > 9 {
		return fmt.Errorf("%w: invalid row, col, or value", ErrUnexpected)
	}
	if ss.board[row][col] != 0 {
		return fmt.Errorf("%w: cell (%d, %d) is already filled", ErrUnexpected, row, col)
	}
	if !ss.possible[row][col].Contains(val) {
		return fmt.Errorf("%w: value %d is not possible for cell (%d, %d)", ErrUnexpected, val, row, col)
	}

	ss.board[row][col] = val
	ss.possible[row][col].Clear()
	err := ss.updateLinked(row, col)
	if err != nil {
		ss.board[row][col] = 0
		return err
	}
	return nil
}

func (ss *SudokuSolver) fillSinglePossibilityInRow(row int) bool {
	if row < 0 || row > 8 {
		return false
	}
	if ss.completed.row[row] {
		return false
	}

	cells := [][2]int{}
	for col := range 9 {
		if ss.board[row][col] == 0 {
			cells = append(cells, [2]int{row, col})
		}
	}

	return ss.resolveSinglePossibility(cells)
}

func (ss *SudokuSolver) fillSinglePossibilityInCol(col int) bool {
	if col < 0 || col > 8 {
		return false
	}
	if ss.completed.col[col] {
		return false
	}

	cells := [][2]int{}
	for row := range 9 {
		if ss.board[row][col] == 0 {
			cells = append(cells, [2]int{row, col})
		}
	}

	return ss.resolveSinglePossibility(cells)
}

func (ss *SudokuSolver) fillSinglePossibilityInBox(boxRow, boxCol int) bool {
	if boxRow < 0 || boxRow > 2 || boxCol < 0 || boxCol > 2 {
		return false
	}
	if ss.completed.box[boxRow][boxCol] {
		return false
	}

	cells := [][2]int{}
	for row := boxRow * 3; row < boxRow*3+3; row++ {
		for col := boxCol * 3; col < boxCol*3+3; col++ {
			if ss.board[row][col] == 0 {
				cells = append(cells, [2]int{row, col})
			}
		}
	}

	return ss.resolveSinglePossibility(cells)
}

func (ss *SudokuSolver) resolveSinglePossibility(cells [][2]int) bool {
	valToPos := make(map[int]*fastset.FastSet[[2]int]) // value -> set of positions
	for i := range cells {
		row, col := cells[i][0], cells[i][1]
		if ss.board[row][col] != 0 {
			continue
		}

		for _, val := range ss.possible[row][col].ToSlice() {
			if valToPos[val] == nil {
				valToPos[val] = fastset.NewFastSet[[2]int]()
			}
			valToPos[val].Add([2]int{row, col})
		}
	}

	found := false
	for val, posSet := range valToPos {
		if posSet.Size() == 1 {
			cell, _ := posSet.Peek()
			row, col := cell[0], cell[1]
			// TODO: figure out the error handling
			err := ss.place(row, col, val)
			found = true
		}
	}

	return found
}

func (ss *SudokuSolver) resolveSpearsInBox(boxRow, boxCol int) bool {
	if ss.completed.box[boxRow][boxCol] {
		return false
	}

	valueToRows := make(map[int][]int) // value -> list of possible rows in the box
	valueToCols := make(map[int][]int) // value -> list of possible columns in the box

	for r := boxRow * 3; r < boxRow*3+3; r++ {
		for c := boxCol * 3; c < boxCol*3+3; c++ {
			if ss.board[r][c] != 0 {
				continue
			}

			for _, val := range ss.possible[r][c].ToSlice() {
				valueToRows[val] = append(valueToRows[val], r)
				valueToCols[val] = append(valueToCols[val], c)
			}
		}
	}

	updated := false
	for val, rows := range valueToRows {
		if len(rows) == 1 {
			row := rows[0]
			for c := 0; c < 9; c++ {
				if c/3 == boxCol {
					continue
				}
				if ss.board[row][c] == 0 && ss.possible[row][c].Contains(val) {
					ss.possible[row][c].Remove(val)
					ss.pq.Push([2]int{row, c})
					updated = true
				}
			}
		}
	}

	for val, cols := range valueToCols {
		if len(cols) == 1 {
			col := cols[0]
			for r := 0; r < 9; r++ {
				if r/3 == boxRow {
					continue
				}
				if ss.board[r][col] == 0 && ss.possible[r][col].Contains(val) {
					ss.possible[r][col].Remove(val)
					ss.pq.Push([2]int{r, col})
					updated = true
				}
			}
		}
	}

	return updated
}

func (ss *SudokuSolver) resolveDoubleSpearsInBoxRow(boxRow int) bool {
	if boxRow < 0 || boxRow > 2 {
		return false
	}

	// for each number, if it is in only 2 rows in 2 boxes, then it must be in the third row in the third box

	// value -> which row in the ith box
	valToRows := [3]map[int]*fastset.FastSet[int]{}
	for i := range 3 {
		valToRows[i] = make(map[int]*fastset.FastSet[int])
	}

	for boxCol := range 3 {
		for i := boxRow * 3; i < boxRow*3+3; i++ {
			for j := boxCol * 3; j < boxCol*3+3; j++ {
				for _, val := range ss.possible[i][j].ToSlice() {
					if valToRows[boxCol][val] == nil {
						valToRows[boxCol][val] = fastset.NewFastSet[int]()
					}

					valToRows[boxCol][val].Add(i)
				}
			}
		}
	}

	for val := 1; val <= 9; val++ {
		valid := true
		single := -1
		for i := range 3 {
			if valToRows[i][val] == nil || valToRows[i][val].Size() == 0 {
				valid = false
				break
			}

			if valToRows[i][val].Size() == 1 {
				if single != -1 {
					valid = false
					break
				}
				single = i
			}

			if valToRows[i][val].Size() > 2 {
				valid = false
				break
			}
		}
		if single == -1 || !valid {
			continue
		}

		var set1 *fastset.FastSet[int]
		var set2 *fastset.FastSet[int]

		for i := range 3 {
			if i == single {
				continue
			}

			if set1 == nil {
				set1 = valToRows[i][val]
			} else {
				set2 = valToRows[i][val]
			}
		}

		if set1.Eq(set2) {
			// remove val from the third box in the same row
			for j := range set1.ToSlice() {
				for i := single * 3; i < single*3+3; i++ {
					if ss.possible[i][j].Contains(val) {
						ss.possible[i][j].Remove(val)
						ss.pq.Push([2]int{i, j})
						return true
					}
				}
			}
		}
	}

	return false
}

func (ss *SudokuSolver) resolveDoubleSpearsInBoxCol(boxCol int) bool {
	if boxCol < 0 || boxCol > 2 {
		return false
	}

	// for each number, if it is in only 2 cols in 2 boxes, then it must be in the third col in the third box

	// value -> which col in the ith box
	valToCols := [3]map[int]*fastset.FastSet[int]{}
	for i := range 3 {
		valToCols[i] = make(map[int]*fastset.FastSet[int])
	}

	for boxRow := range 3 {
		for i := boxCol * 3; i < boxCol*3+3; i++ {
			for j := boxRow * 3; j < boxRow*3+3; j++ {
				for _, val := range ss.possible[i][j].ToSlice() {
					if valToCols[boxRow][val] == nil {
						valToCols[boxRow][val] = fastset.NewFastSet[int]()
					}

					valToCols[boxRow][val].Add(i)
				}
			}
		}
	}

	for val := 1; val <= 9; val++ {
		valid := true
		single := -1
		for i := range 3 {
			if valToCols[i][val] == nil || valToCols[i][val].Size() == 0 {
				valid = false
				break
			}

			if valToCols[i][val].Size() == 1 {
				if single != -1 {
					valid = false
					break
				}
				single = i
			}

			if valToCols[i][val].Size() > 2 {
				valid = false
				break
			}
		}
		if single == -1 || !valid {
			continue
		}

		var set1 *fastset.FastSet[int]
		var set2 *fastset.FastSet[int]

		for i := range 3 {
			if i == single {
				continue
			}

			if set1 == nil {
				set1 = valToCols[i][val]
			} else {
				set2 = valToCols[i][val]
			}
		}

		if set1.Eq(set2) {
			// remove val from the third box in the same col
			for i := range set1.ToSlice() {
				for j := single * 3; j < single*3+3; j++ {
					if ss.possible[i][j].Contains(val) {
						ss.possible[i][j].Remove(val)
						ss.pq.Push([2]int{i, j})
						return true
					}
				}
			}
		}
	}

	return false
}

func (ss *SudokuSolver) resolveHiddenGroupsInRow(row int) bool {
	cells := make([][2]int, 0, 9)
	for col := range 9 {
		cells = append(cells, [2]int{row, col})
	}
	return ss.resolveHiddenGroups(cells)
}

func (ss *SudokuSolver) resolveHiddenGroupsInCol(col int) bool {
	cells := make([][2]int, 0, 9)
	for row := range 9 {
		cells = append(cells, [2]int{row, col})
	}
	return ss.resolveHiddenGroups(cells)
}

func (ss *SudokuSolver) resolveHiddenGroupsInBox(boxRow, boxCol int) bool {
	cells := make([][2]int, 0, 9)
	for row := boxRow * 3; row < boxRow*3+3; row++ {
		for col := boxCol * 3; col < boxCol*3+3; col++ {
			cells = append(cells, [2]int{row, col})
		}
	}
	return ss.resolveHiddenGroups(cells)
}

func (ss *SudokuSolver) resolveHiddenGroups(cells [][2]int) bool {
	valToPos := make(map[int][]int)
	for pos, cell := range cells {
		row, col := cell[0], cell[1]
		if ss.board[row][col] != 0 {
			continue
		}
		for _, val := range ss.possible[row][col].ToSlice() {
			valToPos[val] = append(valToPos[val], pos)
		}
	}

	values := make([]int, 0, len(valToPos))
	for val := range valToPos {
		values = append(values, val)
	}

	for groupSize := 2; groupSize < len(values); groupSize++ {
		chosen := make([]int, 0, groupSize)

		// findGroup looks at all subsets and tests if they form a closed group
		var findGroup func(int) bool
		findGroup = func(next int) bool {
			if len(chosen) == groupSize {
				positions := fastset.NewFastSet[int]()
				members := fastset.NewFastSet[int]()
				for _, val := range chosen {
					members.Add(val)
					for _, col := range valToPos[val] {
						positions.Add(col)
					}
				}

				// not valid if more positions than groupSize
				if positions.Size() != groupSize {
					return false
				}

				changed := false
				for _, pos := range positions.ToSlice() {
					row, col := cells[pos][0], cells[pos][1]
					for _, val := range ss.possible[row][col].ToSlice() {
						if !members.Contains(val) {
							ss.possible[row][col].Remove(val)
							ss.pq.Push([2]int{row, col})
							changed = true
						}
					}
				}
				return changed
			}

			// strange indexing to avoid duplicates per loop
			for index := next; index <= len(values)-(groupSize-len(chosen)); index++ {
				chosen = append(chosen, values[index])
				if findGroup(index + 1) {
					return true
				}
				chosen = chosen[:len(chosen)-1]
			}
			return false
		}

		if findGroup(0) {
			return true
		}
	}

	return false
}

func (ss *SudokuSolver) resolveCertainCell(row, col int) error {
	if ss.possible[row][col].Size() != 1 {
		return ErrUncertainCell
	}

	val, _ := ss.possible[row][col].Peek()
	return ss.place(row, col, val)
}

// updateLinked updates the possible values for all cells based on (row, col).
// It does not update the cell (row, col) itself.
func (ss *SudokuSolver) updateLinked(row, col int) error {
	val := ss.board[row][col]
	if val == 0 {
		return ErrEmptyCell
	}

	rowComplete := true
	colComplete := true
	boxComplete := true

	// remove possible values from row and column
	for k := range 9 {
		if k != row && ss.board[k][col] == 0 {
			ss.possible[k][col].Remove(val)
			ss.pq.Push([2]int{k, col})
		}
		if k != col && ss.board[row][k] == 0 {
			ss.possible[row][k].Remove(val)
			ss.pq.Push([2]int{row, k})
		}

		if ss.board[row][k] == 0 {
			rowComplete = false
		}
		if ss.board[k][col] == 0 {
			colComplete = false
		}
	}

	// remove possible values from box
	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3
	for r := boxRow; r < boxRow+3; r++ {
		for c := boxCol; c < boxCol+3; c++ {
			if (r != row || c != col) && ss.board[r][c] == 0 {
				ss.possible[r][c].Remove(val)
				ss.pq.Push([2]int{r, c})
			}
			if ss.board[r][c] == 0 {
				boxComplete = false
			}
		}
	}

	if rowComplete {
		ss.completed.row[row] = true
	}
	if colComplete {
		ss.completed.col[col] = true
	}
	if boxComplete {
		ss.completed.box[row/3][col/3] = true
	}

	return nil
}
