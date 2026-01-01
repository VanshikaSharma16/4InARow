package game

func NewGame(p1, p2 string) *Game {
	var board [6][7]int

	return &Game{
		Player1: p1,
		Player2: p2,
		Turn:    1,
		Board:   board,
	}
}
