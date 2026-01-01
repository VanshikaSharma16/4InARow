package game

func DropDisc(board *[6][7]int, column int, player int) bool {
	for row := 5; row >= 0; row-- {
		if board[row][column] == 0 {
			board[row][column] = player
			return true
		}
	}
	return false
}
