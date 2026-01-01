# 4 in a Row - Connect Four Game

A real-time multiplayer Connect Four game built with Go backend and React frontend, featuring competitive bot AI, matchmaking, reconnection support, and leaderboard tracking.

## 🎮 Features

- **Real-time Multiplayer**: Play against other players or a competitive bot
- **Smart Matchmaking**: Automatic pairing with 10-second bot fallback
- **Competitive Bot AI**: Strategic bot that blocks wins and creates threats
- **Reconnection Support**: Rejoin games within 30 seconds using game ID
- **Leaderboard**: Track wins and see top players
- **MongoDB Persistence**: Game results stored in MongoDB
- **WebSocket Communication**: Real-time updates via WebSockets

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

3. **Set up MongoDB connection:**
   - Create a `.env` file in the `backend` directory
   - Add your MongoDB connection string:
     ```
     MONGODB_URI=mongodb+srv://username:password@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
     ```
   - See `MONGODB_SETUP.md` for detailed MongoDB Atlas setup instructions

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
│   │   └── handler.go  # Main WebSocket handler
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

### REST API
- `GET /leaderboard` - Get leaderboard data

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
- Check your `.env` file has correct `MONGODB_URI`
- Verify your IP is whitelisted in MongoDB Atlas
- Ensure database user has correct permissions
- See `backend/MONGODB_SETUP.md` for detailed help

### WebSocket Connection Issues
- Ensure backend is running on port 8080
- Check browser console for connection errors
- Verify CORS settings if accessing from different origin

### Bot Not Joining
- Wait 10 seconds after starting a game
- Check backend logs for "Bot joining game" messages
- Ensure no other player is waiting

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
- Completed games are persisted to MongoDB
- Leaderboard aggregates wins from MongoDB
- Bot moves are calculated server-side
- All game logic is validated server-side for security

---

**Enjoy playing 4 in a Row! 🎮**

# connect4
# connect4
# connect4
# 4InARow
# 4InARow
