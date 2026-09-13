package board

import (
	"fmt"
	"strings"

	"github.com/michaelzhan1/sudoku2/internal/utils"
)

// BoardState represents the state of a sudoku board
type BoardState struct {
	// the board itself
	board    Board
	original Board
	solution Board

	// number of empty cells remaining
	remaining int
}

// NewBoardState creates a new BoardState from a given Board
func NewBoardState(board, solution Board) *BoardState {
	remaining := 0
	for row := range 9 {
		for col := range 9 {
			if board[row][col] == 0 {
				remaining++
			}
		}
	}

	original := board // copy
	return &BoardState{
		board:     board,
		remaining: remaining,
		original:  original,
		solution:  solution,
	}
}

// Board returns the current board
func (b *BoardState) Board() Board {
	return b.board
}

// Original returns the original board (the puzzle)
func (b *BoardState) Original() Board {
	return b.original
}

// Solution returns the solution board
func (b *BoardState) Solution() Board {
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

// Set places a value at (row, col). Returns true if successful.
func (b *BoardState) Set(row, col, val int) bool {
	if !utils.CheckBounds(row, col) || !utils.CheckValue(val) {
		return false
	}
	if b.original[row][col] != 0 {
		return false
	}
	if b.board.IsEmpty(row, col) {
		b.remaining--
	}
	b.board[row][col] = val
	return true
}

// Clear removes the value at (row, col).
func (b *BoardState) Clear(row, col int) bool {
	if !utils.CheckBounds(row, col) {
		return false
	}
	if b.original[row][col] != 0 {
		return false
	}
	if !b.board.IsEmpty(row, col) {
		b.remaining++
	}
	b.board[row][col] = 0
	return true
}

// IsFull returns true when there are no empty cells.
func (b *BoardState) IsFull() bool {
	return b.remaining == 0
}

// IsComplete returns true when the board is full and all placements are valid.
func (b *BoardState) IsComplete() bool {
	// TODO: can optimize by tracking number of errors in the board state
	if !b.IsFull() {
		return false
	}
	for row := range 9 {
		for col := range 9 {
			val := b.board[row][col]
			b.board[row][col] = 0
			if !b.board.IsValidPlacement(row, col, val) {
				b.board[row][col] = val
				return false
			}
			b.board[row][col] = val
		}
	}

	return true
}

// GiveHint provides a hint for the cell at (row, col) by filling it with the correct value.
// It then updates the original board state as if the value was always provided.
func (b *BoardState) GiveHint(row, col int) (int, bool) {
	if !utils.CheckBounds(row, col) {
		return 0, false
	}
	if b.board[row][col] != 0 {
		return 0, false
	}
	hint := b.solution[row][col]
	b.board[row][col] = hint
	b.original[row][col] = hint
	b.remaining--
	return hint, true
}
