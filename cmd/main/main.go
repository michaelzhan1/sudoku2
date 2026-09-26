package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/boardstate"
	"github.com/michaelzhan1/sudoku2/internal/cli"
)

func main() {
	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Sudoku!")
	fmt.Println("Difficulty: how many clues to show (30–81). More clues = easier.")

	clues := 35 // default medium
	fmt.Print("Enter number of clues [30-81, default 35]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		if n, err := strconv.Atoi(input); err == nil {
			clues = n
		}
	}

	fmt.Println("\nGenerating puzzle…")
	puzzle, solution, err := board.GenerateBoard(clues, rng)
	if err != nil {
		fmt.Println("Error generating puzzle:", err)
		return
	}
	boardState, err := boardstate.NewBoardState(puzzle, solution)
	if err != nil {
		fmt.Println("Error creating board state:", err)
		return
	}

	const instructions = `Commands:
  <row> <col> <val>  — place a number (1-indexed, val 1-9)
  clear <row> <col>  — remove your entry
  hint <row> <col>   — reveal one cell from the solution
  check <row> <col>  - check one cell
  solve <mode>       — auto-solve the entire puzzle
  reset              - reset the puzzle
  quit               — exit`

	msg := "" // status message shown below the board each turn

	for {
		cli.ClearScreen()
		fmt.Println(instructions)
		fmt.Println()
		fmt.Println(boardState)

		if msg != "" {
			fmt.Println(msg)
			msg = ""
		}

		if boardState.IsComplete() {
			fmt.Println("Congratulations! You solved the puzzle!")
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
			if len(parts) != 2 {
				msg = "Usage: solve <mode>, possible modes: brute, smart"
				continue
			}
			mode := strings.ToLower(parts[1])
			switch mode {
			case "brute":
				ok := boardState.BruteForceSolve()
				if !ok {
					msg = "Unable to solve the puzzle."
				}

			case "smart":
				err := boardState.SmartSolve()
				if err != nil {
					msg = fmt.Sprintf("Unable to solve the puzzle: %s", err.Error())
				}

			default:
				msg = "Unknown solve mode. Possible modes: brute, smart"
			}

		case "hint":
			if len(parts) != 3 {
				msg = "Usage: hint <row> <col>"
				continue
			}
			row, col, ok := cli.ParseRowCol(parts[1], parts[2])
			if !ok {
				msg = "Row and column must each be between 1 and 9."
				continue
			}
			hint, err := boardState.GiveHint(row, col)
			if err != nil {
				msg = fmt.Sprintf("Error giving hint for row %d, column %d: %s", row+1, col+1, err.Error())
			} else {
				msg = fmt.Sprintf("Hint: row %d, column %d has value %d", row+1, col+1, hint)
			}

		case "check":
			if len(parts) != 3 {
				msg = "Usage: check <row> <col>"
				continue
			}
			row, col, ok := cli.ParseRowCol(parts[1], parts[2])
			if !ok {
				msg = "Row and column must each be between 1 and 9."
				continue
			}
			correct, err := boardState.Check(row, col)
			if err != nil {
				msg = fmt.Sprintf("Unable to check row %d, column %d: %s", row+1, col+1, err.Error())
			} else {
				if correct {
					msg = fmt.Sprintf("Row %d, column %d is correct", row+1, col+1)
				} else {
					msg = fmt.Sprintf("Row %d, column %d is incorrect", row+1, col+1)
				}
			}

		case "clear":
			if len(parts) != 3 {
				msg = "Usage: clear <row> <col>"
				continue
			}
			row, col, ok := cli.ParseRowCol(parts[1], parts[2])
			if !ok {
				msg = "Row and column must each be between 1 and 9."
				continue
			}
			err := boardState.Clear(row, col)
			if err != nil {
				msg = fmt.Sprintf("Unable to clear row %d, column %d: %s", row+1, col+1, err.Error())
			}

		case "reset":
			boardState.Reset()
			msg = "Puzzle reset."

		default:
			// Expect: <row> <col> <val>
			if len(parts) != 3 {
				msg = "Unknown command. Type <row> <col> <val> to place a number."
				continue
			}
			row, col, val, ok := cli.ParseRowColVal(parts[0], parts[1], parts[2])
			if !ok {
				msg = "Row, column, and value must each be between 1 and 9."
				continue
			}
			err := boardState.Set(row, col, val)
			if err != nil {
				msg = fmt.Sprintf("Unable to place %d at row %d, column %d: %s", val, row+1, col+1, err.Error())
			}
		}
	}
}
