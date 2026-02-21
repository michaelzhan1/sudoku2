package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/michaelzhan1/sudoku2/sudoku"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Sudoku!")
	fmt.Println("Difficulty: how many clues to show (17–81). More clues = easier.")

	clues := 35 // default medium
	fmt.Print("Enter number of clues [default 35]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		if n, err := strconv.Atoi(input); err == nil {
			clues = n
		}
	}

	fmt.Println("\nGenerating puzzle…")
	puzzle, solution := sudoku.Generate(clues, rng)

	// Keep the original puzzle so we know which cells are fixed.
	original := puzzle

	fmt.Println("\nLet's play! Commands:")
	fmt.Println("  <row> <col> <val>  — place a number (1-indexed, val 1-9)")
	fmt.Println("  clear <row> <col>  — remove your entry")
	fmt.Println("  hint               — reveal one cell from the solution")
	fmt.Println("  solve              — auto-solve the entire puzzle")
	fmt.Println("  quit               — exit")
	fmt.Println()

	for {
		printBoard(&puzzle, &original)

		if sudoku.IsComplete(&puzzle) {
			fmt.Println("\n🎉 Congratulations! You solved the puzzle!")
			break
		}

		fmt.Print("\n> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		switch strings.ToLower(parts[0]) {
		case "quit", "q", "exit":
			fmt.Println("Goodbye!")
			return

		case "solve":
			puzzle = solution
			printBoard(&puzzle, &original)
			fmt.Println("\nPuzzle solved!")
			return

		case "hint":
			giveHint(&puzzle, &solution, rng)

		case "clear":
			if len(parts) != 3 {
				fmt.Println("Usage: clear <row> <col>")
				continue
			}
			row, col, ok := parseRowCol(parts[1], parts[2])
			if !ok {
				continue
			}
			if original[row][col] != 0 {
				fmt.Println("That cell is part of the original puzzle and cannot be cleared.")
				continue
			}
			puzzle.Clear(row, col)

		default:
			// Expect: <row> <col> <val>
			if len(parts) != 3 {
				fmt.Println("Unknown command. Type <row> <col> <val> to place a number.")
				continue
			}
			row, col, ok := parseRowCol(parts[0], parts[1])
			if !ok {
				continue
			}
			val, err := strconv.Atoi(parts[2])
			if err != nil || val < 1 || val > 9 {
				fmt.Println("Value must be between 1 and 9.")
				continue
			}
			if original[row][col] != 0 {
				fmt.Println("That cell is part of the original puzzle and cannot be changed.")
				continue
			}
			if !puzzle.IsValidPlacement(row, col, val) {
				fmt.Printf("Warning: %d at (%d,%d) conflicts with another cell!\n", val, row+1, col+1)
			}
			puzzle.Set(row, col, val)
		}
	}
}

func parseRowCol(rowStr, colStr string) (row, col int, ok bool) {
	r, err1 := strconv.Atoi(rowStr)
	c, err2 := strconv.Atoi(colStr)
	if err1 != nil || err2 != nil || r < 1 || r > 9 || c < 1 || c > 9 {
		fmt.Println("Row and column must each be between 1 and 9.")
		return 0, 0, false
	}
	return r - 1, c - 1, true
}

// printBoard prints the board, highlighting user-entered cells with brackets.
func printBoard(b, original *sudoku.Board) {
	fmt.Println()
	fmt.Println("    1 2 3   4 5 6   7 8 9")
	fmt.Println("  +-------+-------+-------+")
	for row := 0; row < 9; row++ {
		fmt.Printf("%d | ", row+1)
		for col := 0; col < 9; col++ {
			val := b[row][col]
			switch {
			case val == 0:
				fmt.Print(". ")
			case original[row][col] != 0:
				fmt.Printf("%d ", val)
			default:
				// User-entered or hint cell
				fmt.Printf("%d ", val)
			}
			if col == 2 || col == 5 {
				fmt.Print("| ")
			}
		}
		fmt.Println("|")
		if row == 2 || row == 5 {
			fmt.Println("  +-------+-------+-------+")
		}
	}
	fmt.Println("  +-------+-------+-------+")
}

// giveHint reveals a random empty cell from the solution.
func giveHint(puzzle, solution *sudoku.Board, rng *rand.Rand) {
	var empty [][2]int
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle[r][c] == 0 {
				empty = append(empty, [2]int{r, c})
			}
		}
	}
	if len(empty) == 0 {
		fmt.Println("No empty cells to hint!")
		return
	}
	pos := empty[rng.Intn(len(empty))]
	r, c := pos[0], pos[1]
	puzzle[r][c] = solution[r][c]
	fmt.Printf("Hint: placed %d at (%d,%d)\n", solution[r][c], r+1, c+1)
}
