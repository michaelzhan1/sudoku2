package boardstate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/utils"
)

var ErrInvalidOriginal = errors.New("invalid original board; must be a valid sudoku board")
var ErrInvalidSolution = errors.New("invalid solution; must be a valid and complete sudoku solution")
var ErrInvalidParameters = errors.New("Row, column, and value must each be between 1 and 9")
var ErrLockedValue = errors.New("Value is locked")
var ErrEmptyCell = errors.New("Cell is empty")

// BoardState represents the state of a sudoku board
type BoardState struct {
	// the board itself
	board    board.Board
	original board.Board
	solution board.Board

	// number of empty cells remaining
	remaining int
}

// NewBoardState creates a new BoardState from a given Board
func NewBoardState(original, solution board.Board) (*BoardState, error) {
	remaining := 0
	for row := range 9 {
		for col := range 9 {
			if solution[row][col] == 0 || !utils.CheckValue(solution[row][col]) || !solution.IsValidPlacement(row, col, solution[row][col]) {
				return nil, ErrInvalidSolution
			}
			if original[row][col] != 0 && original[row][col] != solution[row][col] {
				return nil, ErrInvalidOriginal
			}

			if original[row][col] == 0 {
				remaining++
			}
		}
	}

	copy := original // copy
	return &BoardState{
		board:     copy,
		remaining: remaining,
		original:  original,
		solution:  solution,
	}, nil
}

// Reset resets the board to its original state
func (b *BoardState) Reset() {
	b.board = b.original
	b.remaining = 0
	for row := range 9 {
		for col := range 9 {
			if b.board[row][col] == 0 {
				b.remaining++
			}
		}
	}
}

// Board returns the current board
func (b *BoardState) Board() board.Board {
	return b.board
}

// Original returns the original board (the puzzle)
func (b *BoardState) Original() board.Board {
	return b.original
}

// Solution returns the solution board
func (b *BoardState) Solution() board.Board {
	return b.solution
}

// Remaining returns the number of empty cells remaining
func (b *BoardState) Remaining() int {

	return b.remaining
}

// String returns a human-readable representation of the board
func (b *BoardState) String() string {
	var sb strings.Builder
	sb.WriteString("    1 2 3   4 5 6   7 8 9\n")
	sb.WriteString("  +-------+-------+-------+\n")
	for row := range 9 {
		sb.WriteString(fmt.Sprintf("%d | ", row+1))
		for col := range 9 {
			val := b.board[row][col]
			switch {
			case val == 0:
				sb.WriteString(". ")
			case b.original[row][col] != 0:
				sb.WriteString(fmt.Sprintf("\033[1m%d\033[0m ", val))
			default:
				sb.WriteString(fmt.Sprintf("%d ", val))
			}
			if col == 2 || col == 5 {
				sb.WriteString("| ")
			}
		}
		sb.WriteString("|\n")
		if row == 2 || row == 5 {
			sb.WriteString("  +-------+-------+-------+\n")
		}
	}
	sb.WriteString("  +-------+-------+-------+\n")

	return sb.String()
}

// Set places a value at (row, col), and returns an error if it fails.
func (b *BoardState) Set(row, col, val int) error {
	if !utils.CheckBounds(row, col) || !utils.CheckValue(val) {
		return ErrInvalidParameters
	}
	if b.original[row][col] != 0 {
		return ErrLockedValue
	}
	if b.board.IsEmpty(row, col) {
		b.remaining--
	}
	b.board[row][col] = val
	return nil
}

// Clear removes the value at (row, col), and returns an error if it fails.
func (b *BoardState) Clear(row, col int) error {
	if !utils.CheckBounds(row, col) {
		return ErrInvalidParameters
	}
	if b.original[row][col] != 0 {
		return ErrLockedValue
	}
	if !b.board.IsEmpty(row, col) {
		b.remaining++
	}
	b.board[row][col] = 0
	return nil
}

// IsComplete returns true when the board is full and all placements are correct.
func (b *BoardState) IsComplete() bool {
	if b.remaining != 0 {
		return false
	}
	for row := range 9 {
		for col := range 9 {
			if b.board[row][col] != b.solution[row][col] {
				return false
			}
		}
	}

	return true
}

// GiveHint provides a hint for the cell at (row, col) by filling it with the correct value.
// It then updates the original board state as if the value was always provided.
func (b *BoardState) GiveHint(row, col int) (int, error) {
	if !utils.CheckBounds(row, col) {
		return 0, ErrInvalidParameters
	}
	if b.original[row][col] != 0 {
		return 0, ErrLockedValue
	}
	hint := b.solution[row][col]
	b.board[row][col] = hint
	b.original[row][col] = hint
	b.remaining--
	return hint, nil
}

// Check returns if a given cell has the correct value, and if the value is unable to be checked
func (b *BoardState) Check(row, col int) (bool, error) {
	if !utils.CheckBounds(row, col) {
		return false, ErrInvalidParameters
	}
	if b.board[row][col] == 0 {
		return false, ErrEmptyCell
	}

	good := b.board[row][col] == b.solution[row][col]
	if good {
		b.original[row][col] = b.solution[row][col]
	}
	return good, nil
}
