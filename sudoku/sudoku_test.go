package sudoku_test

import (
	"math/rand"
	"testing"

	"github.com/michaelzhan1/sudoku2/sudoku"
)

// --- Board tests ---

func TestBoardSetAndGet(t *testing.T) {
	var b sudoku.Board
	if !b.Set(0, 0, 5) {
		t.Fatal("Set should succeed for valid args")
	}
	if b[0][0] != 5 {
		t.Errorf("expected 5, got %d", b[0][0])
	}
}

func TestBoardSetOutOfRange(t *testing.T) {
	var b sudoku.Board
	if b.Set(-1, 0, 5) {
		t.Fatal("Set should fail for negative row")
	}
	if b.Set(0, 9, 5) {
		t.Fatal("Set should fail for col=9")
	}
	if b.Set(0, 0, 10) {
		t.Fatal("Set should fail for val=10")
	}
	if b.Set(0, 0, 0) {
		t.Fatal("Set should fail for val=0")
	}
}

func TestBoardClear(t *testing.T) {
	var b sudoku.Board
	b.Set(3, 4, 7)
	b.Clear(3, 4)
	if !b.IsEmpty(3, 4) {
		t.Error("cell should be empty after Clear")
	}
}

func TestBoardIsValidPlacement(t *testing.T) {
	var b sudoku.Board
	b.Set(0, 0, 5)

	// Same value in same row — invalid
	if b.IsValidPlacement(0, 8, 5) {
		t.Error("5 already in row 0; placement should be invalid")
	}
	// Same value in same column — invalid
	if b.IsValidPlacement(8, 0, 5) {
		t.Error("5 already in col 0; placement should be invalid")
	}
	// Same value in same box — invalid
	if b.IsValidPlacement(1, 1, 5) {
		t.Error("5 already in top-left box; placement should be invalid")
	}
	// Different value — valid
	if !b.IsValidPlacement(0, 8, 6) {
		t.Error("6 not in row/col/box; placement should be valid")
	}
}

func TestBoardIsValidPlacementSameCell(t *testing.T) {
	var b sudoku.Board
	b.Set(4, 4, 9)
	// Checking the same cell with its own value must be valid
	// (we ignore the cell itself when checking).
	if !b.IsValidPlacement(4, 4, 9) {
		t.Error("placement in same cell should ignore existing value")
	}
}

func TestBoardIsFull(t *testing.T) {
	var b sudoku.Board
	if b.IsFull() {
		t.Error("empty board should not be full")
	}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			b[r][c] = 1 // just fill with 1s for fullness check
		}
	}
	if !b.IsFull() {
		t.Error("board with all cells set should be full")
	}
}

// --- Generator tests ---

func TestGenerateProducesValidPuzzle(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	puzzle, solution := sudoku.Generate(35, rng)

	// Solution must be complete and valid.
	if !sudoku.IsComplete(&solution) {
		t.Error("generated solution is not complete/valid")
	}

	// Count clues in puzzle.
	clueCount := 0
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle[r][c] != 0 {
				clueCount++
			}
		}
	}
	if clueCount < 17 {
		t.Errorf("puzzle has fewer than 17 clues: %d", clueCount)
	}
	// Puzzle cells must match the solution where set.
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle[r][c] != 0 && puzzle[r][c] != solution[r][c] {
				t.Errorf("puzzle[%d][%d]=%d but solution[%d][%d]=%d",
					r, c, puzzle[r][c], r, c, solution[r][c])
			}
		}
	}
}

func TestGenerateClampedClues(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	// Requesting more clues than 81 should produce a fully filled board.
	puzzle, solution := sudoku.Generate(100, rng)
	if puzzle != solution {
		t.Error("requesting 100 clues should produce a puzzle equal to the solution")
	}
}

// --- Solver tests ---

func TestIsCompleteOnSolvedPuzzle(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	_, solution := sudoku.Generate(81, rng)
	if !sudoku.IsComplete(&solution) {
		t.Error("full generated solution should pass IsComplete")
	}
}

func TestIsCompleteOnEmptyBoard(t *testing.T) {
	var b sudoku.Board
	if sudoku.IsComplete(&b) {
		t.Error("empty board should not be complete")
	}
}

func TestIsCompleteDetectsConflict(t *testing.T) {
	rng := rand.New(rand.NewSource(13))
	_, solution := sudoku.Generate(81, rng)
	// Introduce a conflict: swap two cells in the same row.
	solution[0][0], solution[0][1] = solution[0][1], solution[0][0]
	if sudoku.IsComplete(&solution) {
		t.Error("board with conflicting cells should not be complete")
	}
}
