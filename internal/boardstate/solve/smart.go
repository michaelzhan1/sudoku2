package solve

import (
	"errors"
	"fmt"
	"sort"

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

func NewSudokuSolver(b *board.Board) (*SudokuSolver, error) {
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

			if err := ss.updateLinked(i, j); err != nil {
				return nil, err
			}
		}
	}

	return ss, nil
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
			if !ss.completed.row[i] {
				didChange, err := ss.fillSinglePossibilityInRow(i)
				if err != nil {
					return err
				}
				if didChange {
					changed = true
					break
				}
			}
			if !ss.completed.col[i] {
				didChange, err := ss.fillSinglePossibilityInCol(i)
				if err != nil {
					return err
				}
				if didChange {
					changed = true
					break
				}
			}
			if !ss.completed.box[i/3][i%3] {
				didChange, err := ss.fillSinglePossibilityInBox(i/3, i%3)
				if err != nil {
					return err
				}
				if didChange {
					changed = true
					break
				}
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

		// use naked groups
		for i := range 9 {
			if ss.resolveNakedGroupsInRow(i) {
				changed = true
				break
			}
			if ss.resolveNakedGroupsInCol(i) {
				changed = true
				break
			}
			if ss.resolveNakedGroupsInBox(i/3, i%3) {
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

func (ss *SudokuSolver) fillSinglePossibilityInRow(row int) (bool, error) {
	if row < 0 || row > 8 {
		return false, nil
	}
	if ss.completed.row[row] {
		return false, nil
	}

	cells := [][2]int{}
	for col := range 9 {
		if ss.board[row][col] == 0 {
			cells = append(cells, [2]int{row, col})
		}
	}

	return ss.resolveSinglePossibility(cells)
}

func (ss *SudokuSolver) fillSinglePossibilityInCol(col int) (bool, error) {
	if col < 0 || col > 8 {
		return false, nil
	}
	if ss.completed.col[col] {
		return false, nil
	}

	cells := [][2]int{}
	for row := range 9 {
		if ss.board[row][col] == 0 {
			cells = append(cells, [2]int{row, col})
		}
	}

	return ss.resolveSinglePossibility(cells)
}

func (ss *SudokuSolver) fillSinglePossibilityInBox(boxRow, boxCol int) (bool, error) {
	if boxRow < 0 || boxRow > 2 || boxCol < 0 || boxCol > 2 {
		return false, nil
	}
	if ss.completed.box[boxRow][boxCol] {
		return false, nil
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

func (ss *SudokuSolver) resolveSinglePossibility(cells [][2]int) (bool, error) {
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
			if err := ss.place(row, col, val); err != nil {
				return false, err
			}
			found = true
		}
	}

	return found, nil
}

func (ss *SudokuSolver) resolveSpearsInBox(boxRow, boxCol int) bool {
	if ss.completed.box[boxRow][boxCol] {
		return false
	}

	valueToRows := make(map[int]*fastset.FastSet[int]) // value -> possible rows in the box
	valueToCols := make(map[int]*fastset.FastSet[int]) // value -> possible columns in the box

	for r := boxRow * 3; r < boxRow*3+3; r++ {
		for c := boxCol * 3; c < boxCol*3+3; c++ {
			if ss.board[r][c] != 0 {
				continue
			}

			for _, val := range ss.possible[r][c].ToSlice() {
				if valueToRows[val] == nil {
					valueToRows[val] = fastset.NewFastSet[int]()
					valueToCols[val] = fastset.NewFastSet[int]()
				}
				valueToRows[val].Add(r)
				valueToCols[val].Add(c)
			}
		}
	}

	updated := false
	for val, rows := range valueToRows {
		if rows.Size() == 1 {
			row, _ := rows.Peek()
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
		if cols.Size() == 1 {
			col, _ := cols.Peek()
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

	for row := boxRow * 3; row < boxRow*3+3; row++ {
		for val := 1; val <= 9; val++ {
			boxCol := -1
			count := 0
			for col := 0; col < 9; col++ {
				if ss.board[row][col] == 0 && ss.possible[row][col].Contains(val) {
					count++
					if boxCol == -1 {
						boxCol = col / 3
					} else if boxCol != col/3 {
						boxCol = -2
						break
					}
				}
			}
			if count == 0 || boxCol < 0 {
				continue
			}

			updated := false
			for r := boxRow * 3; r < boxRow*3+3; r++ {
				if r == row {
					continue
				}
				for col := boxCol * 3; col < boxCol*3+3; col++ {
					if ss.board[r][col] == 0 && ss.possible[r][col].Contains(val) {
						ss.possible[r][col].Remove(val)
						ss.pq.Push([2]int{r, col})
						updated = true
					}
				}
			}
			if updated {
				return true
			}
		}
	}
	return false
}

func (ss *SudokuSolver) resolveDoubleSpearsInBoxCol(boxCol int) bool {
	if boxCol < 0 || boxCol > 2 {
		return false
	}

	for col := boxCol * 3; col < boxCol*3+3; col++ {
		for val := 1; val <= 9; val++ {
			boxRow := -1
			count := 0
			for row := 0; row < 9; row++ {
				if ss.board[row][col] == 0 && ss.possible[row][col].Contains(val) {
					count++
					if boxRow == -1 {
						boxRow = row / 3
					} else if boxRow != row/3 {
						boxRow = -2
						break
					}
				}
			}
			if count == 0 || boxRow < 0 {
				continue
			}

			updated := false
			for row := boxRow * 3; row < boxRow*3+3; row++ {
				for c := boxCol * 3; c < boxCol*3+3; c++ {
					if c == col {
						continue
					}
					if ss.board[row][c] == 0 && ss.possible[row][c].Contains(val) {
						ss.possible[row][c].Remove(val)
						ss.pq.Push([2]int{row, c})
						updated = true
					}
				}
			}
			if updated {
				return true
			}
		}
	}
	return false
}

func (ss *SudokuSolver) resolveNakedGroupsInRow(row int) bool {
	cells := make([][2]int, 0, 9)
	for col := range 9 {
		cells = append(cells, [2]int{row, col})
	}
	return ss.resolveNakedGroups(cells)
}

func (ss *SudokuSolver) resolveNakedGroupsInCol(col int) bool {
	cells := make([][2]int, 0, 9)
	for row := range 9 {
		cells = append(cells, [2]int{row, col})
	}
	return ss.resolveNakedGroups(cells)
}

func (ss *SudokuSolver) resolveNakedGroupsInBox(boxRow, boxCol int) bool {
	cells := make([][2]int, 0, 9)
	for row := boxRow * 3; row < boxRow*3+3; row++ {
		for col := boxCol * 3; col < boxCol*3+3; col++ {
			cells = append(cells, [2]int{row, col})
		}
	}
	return ss.resolveNakedGroups(cells)
}

func (ss *SudokuSolver) resolveNakedGroups(cells [][2]int) bool {
	eligible := make([]int, 0, len(cells))
	for pos, cell := range cells {
		row, col := cell[0], cell[1]
		if ss.board[row][col] == 0 && ss.possible[row][col].Size() >= 2 && ss.possible[row][col].Size() <= 4 {
			eligible = append(eligible, pos)
		}
	}

	for groupSize := 2; groupSize <= 4; groupSize++ {
		chosen := make([]int, 0, groupSize)

		var findGroup func(int) bool
		findGroup = func(next int) bool {
			if len(chosen) == groupSize {
				members := fastset.NewFastSet[int]()
				positions := fastset.NewFastSet[int]()
				for _, index := range chosen {
					row, col := cells[index][0], cells[index][1]
					positions.Add(index)
					for _, val := range ss.possible[row][col].ToSlice() {
						members.Add(val)
					}
				}

				if members.Size() != groupSize {
					return false
				}

				changed := false
				for pos, cell := range cells {
					if positions.Contains(pos) || ss.board[cell[0]][cell[1]] != 0 {
						continue
					}
					for _, val := range ss.possible[cell[0]][cell[1]].ToSlice() {
						if members.Contains(val) {
							ss.possible[cell[0]][cell[1]].Remove(val)
							ss.pq.Push(cell)
							changed = true
						}
					}
				}
				return changed
			}

			for index := next; index <= len(eligible)-(groupSize-len(chosen)); index++ {
				chosen = append(chosen, eligible[index])
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
	sort.Ints(values)

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
