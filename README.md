# 4 in a Row - Connect Four Game

A real-time multiplayer Connect Four game built with Go backend and React frontend, featuring competitive bot AI, matchmaking, reconnection support, and leaderboard tracking.

## 🎮 Features

- **Real-time Multiplayer**: Play against other players or a competitive bot
- **Smart Matchmaking**: Automatic pairing with 10-second bot fallback
- **Competitive Bot AI**: Strategic bot that blocks wins and creates threats
- **Reconnection Support**: Rejoin games within 30 seconds using game ID
- **Auto-updating Leaderboard**: Leaderboard automatically refreshes after each game
- **MongoDB Persistence**: Game results stored in MongoDB (optional - server runs without it)
- **WebSocket Communication**: Real-time updates via WebSockets
- **Game State Management**: Automatic cleanup of finished games, fresh board for new games

## 🛠 Tech Stack

### Backend
- **Go 1.25+**
- **Gorilla WebSocket** - Real-time communication
- **MongoDB** - Game result persistence
- **MongoDB Go Driver** - Database operations

### Frontend
- **React** - UI framework
- **WebSocket API** - Real-time updates

## 📋 Prerequisites

- Go 1.25 or higher
- Node.js 14+ and npm
- MongoDB Atlas account (or local MongoDB instance)

## 🚀 Setup Instructions

### 1. Backend Setup

1. **Navigate to backend directory:**
   ```bash
   cd backend
   ```

2. **Install Go dependencies:**
   ```bash
   go mod download
   ```

3. **Set up MongoDB connection (Optional):**
   - The server will run without MongoDB, but game results won't be saved
   - To enable MongoDB:
     - Create a `.env` file in the `backend` directory
     - Add your MongoDB connection string:
       ```
       MONGODB_URI=mongodb+srv://username:password@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
       ```
     - See `MONGODB_SETUP.md` for detailed MongoDB Atlas setup instructions
     - See `MONGODB_TROUBLESHOOTING.md` if you encounter connection issues

4. **Run the backend server:**
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8080`

### 2. Frontend Setup

1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm start
   ```
   The frontend will start on `http://localhost:3000`

## 🎯 How to Play

1. **Start a New Game:**
   - Enter your username
   - Click "Start New Game"
   - Wait for an opponent (or bot will join in 10 seconds)

2. **Make Moves:**
   - Click on a column to drop your disc
   - Red discs = Player 1, Yellow discs = Player 2
   - First to get 4 in a row wins!

3. **Reconnect to Game:**
   - If you disconnect, save your Game ID
   - Enter your username and Game ID
   - Click "Reconnect" (within 30 seconds)

4. **View Leaderboard:**
   - Scroll down to see the leaderboard
   - Shows total wins per player
   - Automatically updates after each game ends (no page reload needed)

## 🏗 Project Structure

```
connect4/
├── backend/
│   ├── db/              # MongoDB operations
│   │   ├── mongo.go     # Connection setup
│   │   ├── game_result.go  # Save game results
│   │   └── leaderboard.go  # Leaderboard queries
│   ├── game/            # Game logic
│   │   ├── state.go     # Game state structure
│   │   ├── gameplay.go  # Move validation & win/draw detection
│   │   ├── win.go       # Win condition checks
│   │   ├── bot.go       # Bot AI strategy
│   │   └── matchmaker.go # Matchmaking & reconnection
│   ├── websocket/       # WebSocket handlers
│   │   ├── handler.go  # Main WebSocket handler
│   │   └── connections.go # Connection management
│   └── main.go          # Server entry point
├── frontend/
│   └── src/
│       ├── App.js       # Main React component
│       └── Leaderboard.js # Leaderboard component
└── README.md
```

## 🧠 Bot Strategy

The competitive bot uses a multi-priority strategy:

1. **Win if possible** - Check if bot can win immediately
2. **Block opponent win** - Prevent player from winning
3. **Create threats** - Build 3-in-a-row patterns
4. **Block threats** - Prevent opponent's 3-in-a-row
5. **Prefer center** - Center columns are more valuable
6. **Random fallback** - Random valid move if no strategy applies

## 🔌 API Endpoints

### WebSocket
- `ws://localhost:8080/ws?username=<username>&gameId=<gameId>` - Connect to game
  - **Parameters:**
    - `username` (required): Player's username
    - `gameId` (optional): Game ID for reconnection to existing game
  - **Messages:**
    - Client → Server: `{"type": "move", "column": 0-6}`
    - Server → Client: `{"type": "waiting" | "game_started" | "state" | "game_over" | "error", ...}`

### REST API
- `GET /leaderboard` - Get leaderboard data
  - **Response:** JSON array of `{player: string, wins: number}` sorted by wins (descending)
  - **Note:** Only includes players who have won games (draws are excluded)
  - **Example:**
    ```json
    [
      {"player": "alice", "wins": 5},
      {"player": "bob", "wins": 3},
      {"player": "charlie", "wins": 1}
    ]
    ```

## 🔄 Reconnection Flow

1. Player disconnects → Marked as disconnected with timestamp
2. 30-second grace period → Player can reconnect
3. Reconnection → Use username + gameId to rejoin
4. If not reconnected → Game forfeited, opponent wins

## 📊 Database Schema

### Games Collection
```json
{
  "player1": "string",
  "player2": "string",
  "winner": 1 | 2 | 0,  // 0 = draw
  "isDraw": boolean,
  "createdAt": ISODate
}
```

## 🐛 Troubleshooting

### MongoDB Connection Issues
- **Server runs without MongoDB**: If MongoDB connection fails, the server will still start with warnings
- Game results won't be saved, but all other features work normally
- To fix MongoDB connection:
  - Check your `.env` file has correct `MONGODB_URI`
  - Verify your IP is whitelisted in MongoDB Atlas (Network Access)
  - Ensure database user has correct permissions
  - See `backend/MONGODB_SETUP.md` for detailed setup instructions
  - See `backend/MONGODB_TROUBLESHOOTING.md` for common issues and fixes

### WebSocket Connection Issues
- Ensure backend is running on port 8080
- Check browser console for connection errors
- Verify CORS settings if accessing from different origin
- Make sure frontend is connecting to `ws://localhost:8080/ws`

### Bot Not Joining
- Wait 10 seconds after starting a game
- Check backend logs for "Bot joining game" messages
- Ensure no other player is waiting

### Old Board Showing in New Game
- If you see a previous game's board when starting a new game:
  - Make sure you're not reconnecting to an old game (check Game ID)
  - The system automatically cleans up finished games
  - Try refreshing the page if the issue persists

### Leaderboard Not Updating
- Leaderboard automatically refreshes after each game ends
- If it doesn't update:
  - Check if MongoDB is connected (leaderboard requires MongoDB)
  - Wait a moment - there's a 500ms delay to ensure data is saved
  - Check browser console for errors

## 🚧 Future Enhancements

- [ ] Kafka integration for game analytics
- [ ] Game replay functionality
- [ ] Tournament mode
- [ ] Chat functionality
- [ ] Spectator mode
- [ ] Mobile responsive design improvements

## 📝 License

This project is created as a backend engineering intern assignment.

## 👨‍💻 Development Notes

- Backend uses in-memory game state for active games
- Completed games are persisted to MongoDB (if connected)
- Leaderboard aggregates wins from MongoDB (empty if MongoDB not connected)
- Bot moves are calculated server-side
- All game logic is validated server-side for security
- Finished games are automatically cleaned up from memory after 2 seconds
- Server gracefully handles MongoDB connection failures
- Leaderboard auto-refreshes when games end (no page reload needed)
- Game state is properly reset when starting new games with the same username

---

**Enjoy playing 4 in a Row! 🎮**

# connect4
# 4InARow
