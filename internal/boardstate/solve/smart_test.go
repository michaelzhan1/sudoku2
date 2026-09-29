package solve

import (
	"testing"

	"github.com/michaelzhan1/sudoku2/internal/board"
	"github.com/michaelzhan1/sudoku2/internal/utils/fastset"
)

func TestResolveClosedGroupsInRow(t *testing.T) {
	testcases := []struct {
		name               string
		row                [9]int
		possibleCandidates map[int][]int
		exp                bool
	}{
		{
			name: "two_number_group",
			row:  [9]int{0, 0, 0, 4, 5, 6, 7, 8, 9},
			possibleCandidates: map[int][]int{
				0: {1, 2},
				1: {1, 2, 3},
				2: {3},
			},
			exp: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			var possibleSets [9]*fastset.FastSet[int]
			for i := range 9 {
				possibleSets[i] = fastset.NewFastSet[int]()
			}
			for i, vals := range testcase.possibleCandidates {
				for _, val := range vals {
					possibleSets[i].Add(val)
				}
			}

			ss := &SudokuSolver{
				board: &board.Board{
					testcase.row,
				},
				possible: [9][9]*fastset.FastSet[int]{possibleSets},
			}

			res := ss.resolveClosedGroupsInRow(0)
			if res != testcase.exp {
				t.Errorf("expected resolveClosedGroupsInRow to return %v, got %v", testcase.exp, res)
			}

			newSize := ss.possible[0][1].Size()
			if newSize != 2 {
				t.Errorf("expected possible candidates for cell 1 to have size 2, got %d", newSize)
			}
		})
	}
}
