package solve

// // SmartSolve attempts to solve the board using logical techniques
// func (b *BoardState) SmartSolve() bool {
// 	b.Reset()

// 	possible := [9][9][]int{}
// 	changed := true
// 	for changed {
// 		for i := range 9 {
// 			for j := range 9 {
// 				if b.board[i][j] == 0 {
// 					for num := 1; num <= 9; num++ {
// 						if b.board.IsValidPlacement(i, j, num) {
// 							possible[i][j] = append(possible[i][j], num)
// 						}
// 					}
// 				}
// 			}
// 		}

// 		for i := range 81 {
// 			if len(possible[i/9][i%9]) == 1 {
// 				b.board[i/9][i%9] = possible[i/9][i%9][0]
// 				possible[i/9][i%9] = nil
// 			}
// 		}
// 	}

// 	return true
// }
