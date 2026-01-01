let socket = null;

export function connectSocket(username) {
  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    socket.send(JSON.stringify({
      type: "join",
      username: username
    }));
  };

  return socket;
}

export function sendMove(column) {
  if (socket) {
    socket.send(JSON.stringify({
      type: "move",
      column: column
    }));
  }
}
