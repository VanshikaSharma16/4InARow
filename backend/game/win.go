package game

func CheckHorizontal(board [6][7]int, player int) bool {
	for row := 0; row < 6; row++ {
		for col := 0; col < 4; col++ {
			if board[row][col] == player &&
				board[row][col+1] == player &&
				board[row][col+2] == player &&
				board[row][col+3] == player {
				return true
			}
		}
	}
	return false
}

func CheckVertical(board [6][7]int, player int) bool {
	for col := 0; col < 7; col++ {
		for row := 0; row < 3; row++ {
			if board[row][col] == player &&
				board[row+1][col] == player &&
				board[row+2][col] == player &&
				board[row+3][col] == player {
				return true
			}
		}
	}
	return false
}

func CheckDiagonalRight(board [6][7]int, player int) bool {
	for row := 0; row < 3; row++ {
		for col := 0; col < 4; col++ {
			if board[row][col] == player &&
				board[row+1][col+1] == player &&
				board[row+2][col+2] == player &&
				board[row+3][col+3] == player {
				return true
			}
		}
	}
	return false
}

func CheckDiagonalLeft(board [6][7]int, player int) bool {
	for row := 0; row < 3; row++ {
		for col := 3; col < 7; col++ {
			if board[row][col] == player &&
				board[row+1][col-1] == player &&
				board[row+2][col-2] == player &&
				board[row+3][col-3] == player {
				return true
			}
		}
	}
	return false
}

func CheckWin(board [6][7]int, player int) bool {
	return CheckHorizontal(board, player) ||
		CheckVertical(board, player) ||
		CheckDiagonalRight(board, player) ||
		CheckDiagonalLeft(board, player)
}
