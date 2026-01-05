import { useEffect, useState, useRef } from "react";
import Leaderboard from "./Leaderboard";

const WS_URL = "wss://connect4-backend-a4kq.onrender.com/ws";

const emptyBoard = () =>
  Array.from({ length: 6 }, () => Array(7).fill(0));

function App() {
  const [socket, setSocket] = useState(null);
  const [board, setBoard] = useState(emptyBoard());
  const [turn, setTurn] = useState(null);
  const [gameOver, setGameOver] = useState(false);
  const [winner, setWinner] = useState(null);
  const [username, setUsername] = useState("");
  const [opponent, setOpponent] = useState("");
  const [gameId, setGameId] = useState("");
  const [status, setStatus] = useState(""); // waiting, playing, game_over
  const [isConnected, setIsConnected] = useState(false);
  const [showUsernameInput, setShowUsernameInput] = useState(true);
  const [player1, setPlayer1] = useState("");
  const [player2, setPlayer2] = useState("");
  const gameOverRef = useRef(false); // Use ref to track game over state for onclose handler
  const [leaderboardRefresh, setLeaderboardRefresh] = useState(0); // Trigger leaderboard refresh

  const connectWebSocket = (user, gameIdParam = "") => {
    // Reset game state when starting new connection (unless reconnecting to existing game)
    if (!gameIdParam) {
      // New game - reset everything
      setBoard(emptyBoard());
      setTurn(null);
      setGameOver(false);
      setWinner(null);
      setOpponent("");
      setGameId("");
      setPlayer1("");
      setPlayer2("");
    }
    gameOverRef.current = false;
    if (socket) {
      socket.close();
    }

    let wsUrl = `${WS_URL}?username=${encodeURIComponent(user)}`;
    if (gameIdParam) {
      wsUrl += `&gameId=${encodeURIComponent(gameIdParam)}`;
    }

    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log("Connected to server");
      setIsConnected(true);
      setStatus("connecting");
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.type === "waiting") {
        setStatus("waiting");
        setShowUsernameInput(false);
        // Reset board when waiting for opponent (new game)
        setBoard(emptyBoard());
        setTurn(null);
        setGameOver(false);
        setWinner(null);
      }

      if (data.type === "game_started") {
        setStatus("playing");
        setOpponent(data.opponent);
        setShowUsernameInput(false);
      }

      if (data.type === "reconnected") {
        setStatus("playing");
        setGameId(data.gameId);
        if (data.opponent) setOpponent(data.opponent);
        setShowUsernameInput(false);
      }

      if (data.type === "state") {
        setBoard(data.board);
        setTurn(data.turn);
        setGameId(data.gameId || gameId);
        if (data.player1) setPlayer1(data.player1);
        if (data.player2) setPlayer2(data.player2);
        // Update opponent if not set yet
        if (data.player1 && data.player2 && !opponent) {
          const opp = username === data.player1 ? data.player2 : data.player1;
          setOpponent(opp);
        }
        if (data.gameOver) {
          setGameOver(true);
          gameOverRef.current = true; // Update ref
          setStatus("game_over");
          if (data.winner === 0) {
            setWinner("draw");
          } else {
            setWinner(data.winner);
          }
          // Trigger leaderboard refresh when game ends
          setLeaderboardRefresh(prev => prev + 1);
        }
      }

      if (data.type === "game_over") {
        setGameOver(true);
        gameOverRef.current = true; // Update ref
        setStatus("game_over");
        if (data.winner === 0 || data.result === "draw") {
          setWinner("draw");
        } else {
          setWinner(data.winner);
        }
        // Trigger leaderboard refresh when game ends
        setLeaderboardRefresh(prev => prev + 1);
      }

      if (data.type === "error") {
        alert(data.error);
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      setIsConnected(false);
    };

    ws.onclose = () => {
      console.log("Disconnected from server");
      setIsConnected(false);
      // Only set disconnected status if game is not over
      // Use ref to check current state (avoids stale closure issue)
      if (!gameOverRef.current) {
        setStatus(prevStatus => {
          // Don't override game_over status
          if (prevStatus !== "game_over") {
            return "disconnected";
          }
          return prevStatus;
        });
      }
    };

    setSocket(ws);
  };

  const handleStartGame = () => {
    if (!username.trim()) {
      alert("Please enter a username");
      return;
    }
    connectWebSocket(username.trim());
  };

  const handleReconnect = () => {
    if (!username.trim() || !gameId.trim()) {
      alert("Please enter username and game ID");
      return;
    }
    connectWebSocket(username.trim(), gameId.trim());
  };

  function handleMove(column) {
    if (!socket || gameOver || status !== "playing") return;
    if (turn === null) return;

    // Check if it's player's turn - determine player number from player1/player2
    let playerNum = 0;
    if (username === player1) {
      playerNum = 1;
    } else if (username === player2) {
      playerNum = 2;
    } else {
      // Fallback: if player1/player2 not set yet, use old logic
      playerNum = username === opponent ? 2 : 1;
    }
    
    if (turn !== playerNum) {
      alert("Not your turn!");
      return;
    }

    socket.send(
      JSON.stringify({
        type: "move",
        column,
      })
    );
  }

  function resetGame() {
    setBoard(emptyBoard());
    setTurn(null);
    setGameOver(false);
    gameOverRef.current = false; // Reset ref
    setWinner(null);
    setOpponent("");
    setGameId("");
    setStatus("");
    setShowUsernameInput(true);
    setPlayer1("");
    setPlayer2("");
    if (socket) {
      socket.close();
    }
    setSocket(null);
  }

  const getStatusMessage = () => {
    if (status === "waiting") {
      return "Waiting for opponent... Bot will join in 10 seconds if no player found.";
    }
    if (status === "connecting") {
      return "Connecting...";
    }
    if (status === "playing") {
      // Determine player number from player1/player2
      let playerNum = 0;
      if (username === player1) {
        playerNum = 1;
      } else if (username === player2) {
        playerNum = 2;
      } else {
        // Fallback
        playerNum = username === opponent ? 2 : 1;
      }
      
      if (gameOver) {
        if (winner === "draw") {
          return "Game Over - It's a Draw!";
        }
        const winnerName = winner === 1 ? player1 : player2;
        const isYou = (winner === playerNum);
        return `Game Over - ${winnerName}${isYou ? " (You)" : ""} wins!`;
      }
      
      const currentPlayerName = turn === 1 ? player1 : player2;
      const isYourTurn = (turn === playerNum);
      return `Turn: ${currentPlayerName}${isYourTurn ? " (You)" : ""}`;
    }
    if (status === "game_over") {
      if (winner === "draw") {
        return "Game Over - It's a Draw!";
      }
      return `Game Over - Player ${winner} wins!`;
    }
    if (status === "disconnected") {
      return "Disconnected. You can reconnect within 30 seconds using your game ID.";
    }
    return "";
  };

  return (
    <div style={{ padding: 20, maxWidth: 800, margin: "0 auto" }}>
      <h1>4 in a Row - Connect Four</h1>

      {showUsernameInput && (
        <div style={{ marginBottom: 20, padding: 15, border: "1px solid #ccc", borderRadius: 5 }}>
          <h3>Enter Username</h3>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Enter your username"
            style={{ padding: 8, marginRight: 10, width: 200 }}
            onKeyPress={(e) => e.key === "Enter" && handleStartGame()}
          />
          <button onClick={handleStartGame} style={{ padding: 8, marginRight: 10 }}>
            Start New Game
          </button>
          <div style={{ marginTop: 10 }}>
            <h4>Reconnect to existing game:</h4>
            <input
              type="text"
              value={gameId}
              onChange={(e) => setGameId(e.target.value)}
              placeholder="Enter game ID"
              style={{ padding: 8, marginRight: 10, width: 200 }}
            />
            <button onClick={handleReconnect} style={{ padding: 8 }}>
              Reconnect
            </button>
          </div>
        </div>
      )}

      {status && (
        <div style={{ marginBottom: 15, padding: 10, background: "#f0f0f0", borderRadius: 5 }}>
          <strong>Status:</strong> {getStatusMessage()}
          {gameId && (
            <div style={{ marginTop: 5, fontSize: 12, color: "#666" }}>
              Game ID: {gameId} (Save this to reconnect)
            </div>
          )}
        </div>
      )}

      {opponent && (
        <div style={{ marginBottom: 10 }}>
          <strong>Opponent:</strong> {opponent}
        </div>
      )}

      <div style={{ display: "flex", flexDirection: "column", gap: 5, marginBottom: 20, alignItems: "center" }}>
        {/* Column headers - clickable to drop discs */}
        <div style={{ display: "flex", gap: 5, marginBottom: 5 }}>
          {Array.from({ length: 7 }).map((_, c) => (
            <div
              key={`header-${c}`}
              onClick={() => handleMove(c)}
              style={{
                width: 50,
                height: 30,
                background: "#ddd",
                border: "2px solid #333",
                borderRadius: 5,
                cursor: gameOver || status !== "playing" ? "not-allowed" : "pointer",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: 12,
                fontWeight: "bold",
              }}
            >
              ↓
            </div>
          ))}
        </div>
        {/* Game board */}
        {board.map((row, r) => (
          <div key={r} style={{ display: "flex", gap: 5 }}>
            {row.map((cell, c) => (
              <div
                key={`${r}-${c}`}
                style={{
                  width: 50,
                  height: 50,
                  borderRadius: "50%",
                  background:
                    cell === 1 ? "red" : cell === 2 ? "yellow" : "#eee",
                  border: "2px solid #333",
                  opacity: cell === 0 ? 0.3 : 1,
                }}
              />
            ))}
          </div>
        ))}
      </div>

      {gameOver && (
        <button onClick={resetGame} style={{ padding: 10, fontSize: 16, marginBottom: 20 }}>
          New Game
        </button>
      )}

      <hr />
      <Leaderboard refreshTrigger={leaderboardRefresh} />
      
      <footer style={{ 
        marginTop: 40, 
        padding: 20, 
        textAlign: "center", 
        color: "#666", 
        fontSize: 14,
        borderTop: "1px solid #eee"
      }}>
        <p>Developed by <strong>Vanshika Sharma</strong></p>
      </footer>
    </div>
  );
}

export default App;
