package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"connect4/db"
	"connect4/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// ✅ username and gameID from query (for reconnection)
	username := r.URL.Query().Get("username")
	gameID := r.URL.Query().Get("gameId")

	if username == "" {
		username = "player1"
	}

	log.Printf("Player connected: %s (gameID: %s)", username, gameID)

	var g *game.Game

	// ✅ RECONNECTION: Try to find existing game
	if gameID != "" {
		g = game.FindGameByID(gameID)
		if g != nil {
			// Verify username matches
			if username != g.Player1 && username != g.Player2 {
				sendError(conn, "Username doesn't match this game")
				return
			}
			// Mark as reconnected
			if g.Connections == nil {
				g.Connections = make(map[string]bool)
			}
			g.Connections[username] = true
			delete(g.LastSeen, username)
			log.Printf("Player %s reconnected to game %s", username, gameID)
			sendMessage(conn, "reconnected", map[string]interface{}{
				"message": "Reconnected to game",
				"gameId":  g.ID,
			})
			sendState(conn, g)
			// Continue to game loop below
		}
	}

	// ✅ MATCHMAKING: If no game found, try to find or create one
	if g == nil {
		g = game.FindGameByUsername(username)
		if g == nil {
			g = game.FindMatch(username)
		}
	}

	if g == nil {
		// Send waiting message to client
		sendMessage(conn, "waiting", map[string]interface{}{
			"message": "Waiting for opponent... Bot will join in 10 seconds if no player found.",
		})

		// Start bot timer
		game.StartBotIfNoPlayer(username)

		// Wait until game is assigned (either match found or bot joined)
		for {
			time.Sleep(500 * time.Millisecond)
			g = game.ActiveGames[username]
			if g != nil {
				// Send game started message
				if g.Player2 == "BOT" {
					sendMessage(conn, "game_started", map[string]interface{}{
						"message":  "Bot joined! Game starting...",
						"opponent": "BOT",
					})
				} else {
					sendMessage(conn, "game_started", map[string]interface{}{
						"message":  "Match found! Game starting...",
						"opponent": g.Player2,
					})
				}
				// Send initial game state
				sendState(conn, g)
				break
			}
		}
	} else {
		// Match found immediately
		sendMessage(conn, "game_started", map[string]interface{}{
			"message":  "Match found! Game starting...",
			"opponent": g.Player2,
		})
		sendState(conn, g)
	}

	// Mark connection as active
	if g.Connections == nil {
		g.Connections = make(map[string]bool)
	}
	g.Connections[username] = true

	// Start disconnect monitoring goroutine
	go monitorDisconnection(g, username)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for %s: %v", username, err)
			handleDisconnect(g, username)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			sendError(conn, "Invalid JSON")
			continue
		}

		if data["type"] == "move" {

			colFloat, ok := data["column"].(float64)
			if !ok {
				sendError(conn, "Invalid column")
				continue
			}

			col := int(colFloat)

			// ✅ PLAYER MOVE - Check if it's player's turn
			playerNum := 1
			if username == g.Player2 {
				playerNum = 2
			}

			// Make the move using function version
			result := game.MakeMove(g, col, playerNum)
			if result != "OK" && result != "WIN" {
				sendError(conn, result)
				continue
			}

			sendState(conn, g)

			// ✅ PLAYER WIN or DRAW
			if result == "WIN" || result == "DRAW" {
				saveAndEndGame(conn, g)
				return
			}

			// 🤖 BOT MOVE (if bot is opponent and it's bot's turn)
			if g.Player2 == "BOT" && !g.GameOver && g.Turn == 2 {
				time.Sleep(700 * time.Millisecond) // feels human 😄

				botCol := game.BotMove(g)
				botResult := game.MakeMove(g, botCol, 2)

				sendState(conn, g)

				// Check if bot won or draw
				if botResult == "WIN" || botResult == "DRAW" {
					saveAndEndGame(conn, g)
					return
				}
			}
		}
	}
}

// ---------------- helpers ----------------

func saveAndEndGame(conn *websocket.Conn, g *game.Game) {
	err := db.SaveGameResult(g.Player1, g.Player2, g.Winner)
	if err != nil {
		log.Println("MongoDB save failed:", err)
	}
	sendGameOver(conn, g)
}

func sendState(conn *websocket.Conn, g *game.Game) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":     "state",
		"board":    g.Board,
		"turn":     g.Turn,
		"gameId":   g.ID,
		"gameOver": g.GameOver,
		"winner":   g.Winner,
	})
	conn.WriteMessage(websocket.TextMessage, data)
}

func sendError(conn *websocket.Conn, msg string) {
	data, _ := json.Marshal(map[string]string{
		"type":  "error",
		"error": msg,
	})
	conn.WriteMessage(websocket.TextMessage, data)
}

func sendGameOver(conn *websocket.Conn, g *game.Game) {
	result := "draw"
	if g.Winner == 1 {
		result = g.Player1
	} else if g.Winner == 2 {
		result = g.Player2
	}

	data, _ := json.Marshal(map[string]interface{}{
		"type":   "game_over",
		"winner": g.Winner,
		"result": result,
		"board":  g.Board,
	})
	conn.WriteMessage(websocket.TextMessage, data)
}

func sendMessage(conn *websocket.Conn, msgType string, data map[string]interface{}) {
	data["type"] = msgType
	jsonData, _ := json.Marshal(data)
	conn.WriteMessage(websocket.TextMessage, jsonData)
}

// monitorDisconnection checks if player disconnects and handles forfeit after 30 seconds
func monitorDisconnection(g *game.Game, username string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if g.GameOver {
			return
		}

		// Check if player is still connected
		if g.Connections != nil && g.Connections[username] {
			continue
		}

		// Check last seen time
		if lastSeen, exists := g.LastSeen[username]; exists {
			if time.Since(lastSeen) > 30*time.Second {
				log.Printf("Player %s forfeited (disconnected > 30s)", username)
				game.ForfeitGame(g, username)
				// Notify other player if connected
				notifyForfeit(g, username)
				return
			}
		}
	}
}

// handleDisconnect marks player as disconnected
func handleDisconnect(g *game.Game, username string) {
	if g.Connections != nil {
		g.Connections[username] = false
	}
	if g.LastSeen == nil {
		g.LastSeen = make(map[string]time.Time)
	}
	g.LastSeen[username] = time.Now()
	log.Printf("Player %s disconnected (30s grace period)", username)
}

// notifyForfeit notifies the opponent about forfeit (would need connection tracking)
func notifyForfeit(g *game.Game, username string) {
	// This would need a connection manager to notify the other player
	// For now, the next state update will show game over
}
