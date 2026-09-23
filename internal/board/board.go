package board

import (
	"errors"
	"math/rand/v2"

	"github.com/michaelzhan1/sudoku2/internal/utils"
)

var ErrInvalidClues = errors.New("invalid number of clues; must be between 30 and 81")

// Board is a 9x9 sudoku board
type Board [9][9]int

// GenerateBoard generates a sudoku puzzle with its solution.
func GenerateBoard(clues int, rng *rand.Rand) (puzzle, solution Board, err error) {
	if clues < 30 || clues > 81 {
		return Board{}, Board{}, ErrInvalidClues
	}

	solution = Board{}
	solution.fillBoard(rng)
	puzzle = solution // copy

	toRemove := 81 - clues
	positions := rng.Perm(81)
	removed := 0
	for _, pos := range positions {
		if removed == toRemove {
			break
		}
		row, col := pos/9, pos%9
		saved := puzzle[row][col]
		puzzle[row][col] = 0
		if puzzle.countSolutions(2) == 1 {
			removed++
		} else {
			puzzle[row][col] = saved
		}
	}
	return puzzle, solution, nil
}

// IsValidPlacement checks whether val can be placed at (row, col) without
// violating sudoku rules (ignores the current value at that cell).
// The current value at (row, col) is ignored.
func (b *Board) IsValidPlacement(row, col, val int) bool {
	if !utils.CheckBounds(row, col) || !utils.CheckValue(val) {
		return false
	}

	// row check
	for c := range 9 {
		if c != col && b[row][c] == val {
			return false
		}
	}

	// column check
	for r := range 9 {
		if r != row && b[r][col] == val {
			return false
		}
	}

	// box check
	for dr := range 3 {
		for dc := range 3 {
			r, c := (row/3)*3+dr, (col/3)*3+dc
			if (r != row || c != col) && b[r][c] == val {
				return false
			}
		}
	}

	return true
}

// IsEmpty returns whether the cell at (row, col) is empty
func (b *Board) IsEmpty(row, col int) bool {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return false
	}
	return b[row][col] == 0
}

// fillBoard fills a board using backtracking with random values
func (b *Board) fillBoard(rng *rand.Rand) bool {
	for row := range 9 {
		for col := range 9 {
			if b[row][col] != 0 {
				continue
			}
			vals := rng.Perm(9)
			for _, v := range vals {
				val := v + 1
				if b.IsValidPlacement(row, col, val) {
					b[row][col] = val
					if b.fillBoard(rng) {
						return true
					}
					b[row][col] = 0
				}
			}
			return false
		}
	}
	return true
}

// countSolutions counts the number of solutions up to limit using backtracking.
func (b *Board) countSolutions(limit int) int {
	for row := range 9 {
		for col := range 9 {
			if b[row][col] != 0 {
				continue
			}
			count := 0
			for v := range 9 {
				val := v + 1
				if b.IsValidPlacement(row, col, val) {
					b[row][col] = val
					count += b.countSolutions(limit - count)
					b[row][col] = 0
					if count >= limit {
						return count
					}
				}
			}
			return count
		}
	}

	// full board case
	return 1
}
