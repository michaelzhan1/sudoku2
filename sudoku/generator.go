package sudoku

import "math/rand/v2"

// Generate creates a new Sudoku puzzle by:
//  1. Building a fully solved board.
//  2. Removing `clues` cells while ensuring a unique solution.
//
// clues must be between 17 and 81. Values outside this range are clamped.
func Generate(clues int, rng *rand.Rand) (puzzle Board, solution Board) {
	if clues < 17 {
		clues = 17
	}
	if clues > 81 {
		clues = 81
	}

	// Build a complete solution.
	solution = Board{}
	fillBoard(&solution, rng)

	// Start the puzzle as a copy of the solution.
	puzzle = solution

	// Determine how many cells to remove.
	toRemove := 81 - clues

	// Try removing cells in random order, keeping the puzzle uniquely solvable.
	positions := rng.Perm(81)
	removed := 0
	for _, pos := range positions {
		if removed == toRemove {
			break
		}
		row, col := pos/9, pos%9
		saved := puzzle.Get(row, col)
		puzzle.Clear(row, col)
		if countSolutions(&puzzle, 2) == 1 {
			removed++
		} else {
			puzzle.Set(row, col, saved)
		}
	}
	return puzzle, solution
}

// fillBoard fills a board using backtracking with random value ordering.
func fillBoard(b *Board, rng *rand.Rand) bool {
	for row := range 9 {
		for col := range 9 {
			if !b.IsEmpty(row, col) {
				continue
			}
			vals := rng.Perm(9)
			for _, v := range vals {
				val := v + 1
				if b.IsValidPlacement(row, col, val) {
					b.Set(row, col, val)
					if fillBoard(b, rng) {
						return true
					}
					b.Clear(row, col)
				}
			}
			return false
		}
	}
	return true
}

// countSolutions counts the number of solutions up to limit using backtracking.
func countSolutions(b *Board, limit int) int {
	for row := range 9 {
		for col := range 9 {
			if !b.IsEmpty(row, col) {
				continue
			}
			count := 0
			for v := range 9 {
				val := v + 1
				if b.IsValidPlacement(row, col, val) {
					b.Set(row, col, val)
					count += countSolutions(b, limit-count)
					b.Clear(row, col)
					if count >= limit {
						return count
					}
				}
			}
			return count
		}
	}
	return 1 // board is full — this is a solution
}
