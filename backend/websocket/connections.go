package websocket

import (
	"github.com/gorilla/websocket"
)

// ConnectionManager manages WebSocket connections for games
type ConnectionManager struct {
	connections map[string]map[string]*websocket.Conn // gameID -> username -> conn
}

var manager = &ConnectionManager{
	connections: make(map[string]map[string]*websocket.Conn),
}

// AddConnection adds a connection for a game
func (cm *ConnectionManager) AddConnection(gameID, username string, conn *websocket.Conn) {
	if cm.connections[gameID] == nil {
		cm.connections[gameID] = make(map[string]*websocket.Conn)
	}
	cm.connections[gameID][username] = conn
}

// RemoveConnection removes a connection
func (cm *ConnectionManager) RemoveConnection(gameID, username string) {
	if cm.connections[gameID] != nil {
		delete(cm.connections[gameID], username)
		if len(cm.connections[gameID]) == 0 {
			delete(cm.connections, gameID)
		}
	}
}

// BroadcastToGame sends a message to all players in a game
func (cm *ConnectionManager) BroadcastToGame(gameID string, message []byte) {
	if cm.connections[gameID] == nil {
		return
	}
	
	for username, conn := range cm.connections[gameID] {
		if conn != nil {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				// Connection might be closed, remove it
				delete(cm.connections[gameID], username)
			}
		}
	}
}

// SendToPlayer sends a message to a specific player
func (cm *ConnectionManager) SendToPlayer(gameID, username string, message []byte) {
	if cm.connections[gameID] != nil {
		if conn := cm.connections[gameID][username]; conn != nil {
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

