package game

// CheckDraw checks if the board is full (draw condition)
func CheckDraw(board [6][7]int) bool {
	for col := 0; col < 7; col++ {
		if board[0][col] == 0 {
			return false // At least one column has space
		}
	}
	return true // All columns are full
}

func MakeMove(g *Game, column int, player int) string {

	if g.GameOver {
		return "Game already finished"
	}

	if g.Turn != player {
		return "Not your turn"
	}

	if column < 0 || column > 6 {
		return "Invalid column"
	}

	// drop disc
	if !DropDisc(&g.Board, column, player) {
		return "Column full"
	}

	// check win
	if CheckWin(g.Board, player) {
		g.GameOver = true
		g.Winner = player
		return "WIN"
	}

	// check draw (board full)
	if CheckDraw(g.Board) {
		g.GameOver = true
		g.Winner = 0 // 0 = draw
		return "DRAW"
	}

	// switch turn
	if g.Turn == 1 {
		g.Turn = 2
	} else {
		g.Turn = 1
	}

	return "OK"
}
