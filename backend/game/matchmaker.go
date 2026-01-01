package game

import (
	"log"
	"time"
)

var WaitingPlayer string
var ActiveGames = make(map[string]*Game)

func FindMatch(username string) *Game {

	// If someone is waiting → pair them
	if WaitingPlayer != "" && WaitingPlayer != username {

		var board [6][7]int
		game := &Game{
			ID:          GenerateGameID(),
			Player1:     WaitingPlayer,
			Player2:     username,
			Turn:        1,
			Board:       board,
			GameOver:    false,
			Winner:      0,
			LastSeen:    make(map[string]time.Time),
			Connections: make(map[string]bool),
		}

		ActiveGames[game.Player1] = game
		ActiveGames[game.Player2] = game

		WaitingPlayer = ""
		return game
	}

	// Else wait
	WaitingPlayer = username
	return nil
}

// FindGameByID finds a game by its ID
func FindGameByID(gameID string) *Game {
	for _, g := range ActiveGames {
		if g.ID == gameID {
			return g
		}
	}
	return nil
}

// FindGameByUsername finds a game for a username (for reconnection)
func FindGameByUsername(username string) *Game {
	return ActiveGames[username]
}

func StartBotIfNoPlayer(username string) {
	time.AfterFunc(10*time.Second, func() {
		// Check if player is still waiting
		if WaitingPlayer == username {
			log.Printf("Bot joining game for player: %s", username)
			
			var board [6][7]int
			game := &Game{
				ID:          GenerateGameID(),
				Player1:     username,
				Player2:     "BOT",
				Turn:        1,
				Board:       board,
				GameOver:    false,
				Winner:      0,
				LastSeen:    make(map[string]time.Time),
				Connections: make(map[string]bool),
			}

			ActiveGames[username] = game
			WaitingPlayer = ""
			log.Printf("Bot game created for %s (Game ID: %s)", username, game.ID)
		}
	})
}

// ForfeitGame marks the game as forfeited when a player disconnects
func ForfeitGame(g *Game, username string) {
	if g.GameOver {
		return
	}

	g.GameOver = true

	// Determine winner (opponent wins)
	if username == g.Player1 {
		g.Winner = 2
	} else {
		g.Winner = 1
	}

	log.Printf("Game forfeited by %s, winner: %d", username, g.Winner)
}
