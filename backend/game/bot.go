package game

import "math/rand"

// copyBoard creates a copy of the board
func copyBoard(board [6][7]int) [6][7]int {
	var newBoard [6][7]int
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			newBoard[i][j] = board[i][j]
		}
	}
	return newBoard
}

// CountThreats counts how many 3-in-a-row threats a player has
func CountThreats(board [6][7]int, player int) int {
	count := 0
	// Check all possible 3-in-a-row patterns
	// This is a simplified version - full implementation would check all directions
	for row := 0; row < 6; row++ {
		for col := 0; col < 5; col++ {
			// Horizontal threat
			if board[row][col] == player && board[row][col+1] == player && board[row][col+2] == player {
				count++
			}
		}
	}
	return count
}

func BotMove(g *Game) int {
	// 1️⃣ PRIORITY: Try winning move (bot can win)
	for c := 0; c < 7; c++ {
		tempBoard := copyBoard(g.Board)
		if DropDisc(&tempBoard, c, 2) && CheckWin(tempBoard, 2) {
			return c
		}
	}

	// 2️⃣ PRIORITY: Block player's winning move
	for c := 0; c < 7; c++ {
		tempBoard := copyBoard(g.Board)
		if DropDisc(&tempBoard, c, 1) && CheckWin(tempBoard, 1) {
			return c
		}
	}

	// 3️⃣ STRATEGY: Create own threat (3 in a row that can become 4)
	bestCol := -1
	maxThreats := -1
	for c := 0; c < 7; c++ {
		tempBoard := copyBoard(g.Board)
		if DropDisc(&tempBoard, c, 2) {
			threats := CountThreats(tempBoard, 2)
			if threats > maxThreats {
				maxThreats = threats
				bestCol = c
			}
		}
	}
	if bestCol != -1 && maxThreats > 0 {
		return bestCol
	}

	// 4️⃣ STRATEGY: Block player's threat (prevent 3 in a row)
	for c := 0; c < 7; c++ {
		tempBoard := copyBoard(g.Board)
		if DropDisc(&tempBoard, c, 1) {
			threats := CountThreats(tempBoard, 1)
			if threats > 0 {
				return c
			}
		}
	}

	// 5️⃣ PREFERENCE: Center columns are more valuable
	centerCols := []int{3, 2, 4, 1, 5, 0, 6}
	for _, c := range centerCols {
		if g.Board[0][c] == 0 {
			return c
		}
	}

	// 6️⃣ FALLBACK: Random valid column
	valid := []int{}
	for c := 0; c < 7; c++ {
		if g.Board[0][c] == 0 {
			valid = append(valid, c)
		}
	}
	if len(valid) > 0 {
	return valid[rand.Intn(len(valid))]
	}

	return 3 // Default to center
}
