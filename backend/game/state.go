package game

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Game struct {
	ID       string
	Player1  string
	Player2  string
	Turn     int
	Board    [6][7]int
	GameOver bool
	Winner   int

	LastSeen    map[string]time.Time
	Connections map[string]bool // Track active connections
}

// GenerateGameID creates a unique game ID
func GenerateGameID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
