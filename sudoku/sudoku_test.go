package sudoku_test

import (
	"math/rand/v2"
	"testing"

	"github.com/michaelzhan1/sudoku2/sudoku"
)

// --- Board tests ---

func TestBoardSetAndGet(t *testing.T) {
	var b sudoku.Board
	if !b.Set(0, 0, 5) {
		t.Fatal("Set should succeed for valid args")
	}
	if b.Get(0, 0) != 5 {
		t.Errorf("expected 5, got %d", b.Get(0, 0))
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
	for r := range 9 {
		for c := range 9 {
			b.Set(r, c, 1) // just fill with 1s for fullness check
		}
	}
	if !b.IsFull() {
		t.Error("board with all cells set should be full")
	}
}

func TestBoardIsCompleteTracking(t *testing.T) {
	// Verify that IsComplete reflects incremental tracking correctly.
	rng := rand.New(rand.NewPCG(99, 0))
	puzzle, solution := sudoku.Generate(35, rng)

	// Puzzle is incomplete — IsComplete must be false.
	if sudoku.IsComplete(&puzzle) {
		t.Error("partial puzzle should not be complete")
	}

	// Fill every empty cell from the solution.
	for r := range 9 {
		for c := range 9 {
			if puzzle.IsEmpty(r, c) {
				puzzle.Set(r, c, solution.Get(r, c))
			}
		}
	}
	if !sudoku.IsComplete(&puzzle) {
		t.Error("fully and correctly filled puzzle should be complete")
	}

	// Overwrite one cell with a wrong value to introduce a conflict,
	// then verify IsComplete returns false again.
	val := puzzle.Get(0, 0)
	wrong := val%9 + 1 // a different value 1–9
	puzzle.Set(0, 0, wrong)
	if sudoku.IsComplete(&puzzle) {
		t.Error("board with a conflict should not be complete")
	}

	// Restore correct value — should be complete again.
	puzzle.Set(0, 0, val)
	if !sudoku.IsComplete(&puzzle) {
		t.Error("restored board should be complete again")
	}
}

// --- Generator tests ---

func TestGenerateProducesValidPuzzle(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 0))
	puzzle, solution := sudoku.Generate(35, rng)

	// Solution must be complete and valid.
	if !sudoku.IsComplete(&solution) {
		t.Error("generated solution is not complete/valid")
	}

	// Count clues in puzzle.
	clueCount := 0
	for r := range 9 {
		for c := range 9 {
			if puzzle.Get(r, c) != 0 {
				clueCount++
			}
		}
	}
	if clueCount < 17 {
		t.Errorf("puzzle has fewer than 17 clues: %d", clueCount)
	}
	// Puzzle cells must match the solution where set.
	for r := range 9 {
		for c := range 9 {
			if puzzle.Get(r, c) != 0 && puzzle.Get(r, c) != solution.Get(r, c) {
				t.Errorf("puzzle[%d][%d]=%d but solution[%d][%d]=%d",
					r, c, puzzle.Get(r, c), r, c, solution.Get(r, c))
			}
		}
	}
}

func TestGenerateClampedClues(t *testing.T) {
	rng := rand.New(rand.NewPCG(0, 0))
	// Requesting more clues than 81 should produce a puzzle identical to the solution.
	puzzle, solution := sudoku.Generate(100, rng)
	for r := range 9 {
		for c := range 9 {
			if puzzle.Get(r, c) != solution.Get(r, c) {
				t.Errorf("requesting 100 clues: puzzle[%d][%d]=%d but solution[%d][%d]=%d",
					r, c, puzzle.Get(r, c), r, c, solution.Get(r, c))
			}
		}
	}
}

// --- Solver tests ---

func TestIsCompleteOnSolvedPuzzle(t *testing.T) {
	rng := rand.New(rand.NewPCG(7, 0))
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
	rng := rand.New(rand.NewPCG(13, 0))
	_, solution := sudoku.Generate(81, rng)
	// Introduce a conflict: swap two adjacent cells in the same row.
	v00, v01 := solution.Get(0, 0), solution.Get(0, 1)
	solution.Set(0, 0, v01)
	solution.Set(0, 1, v00)
	if sudoku.IsComplete(&solution) {
		t.Error("board with conflicting cells should not be complete")
	}
}
