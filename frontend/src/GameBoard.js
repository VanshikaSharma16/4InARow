import "./board.css";
import { sendMove } from "./socket";

export default function GameBoard({ board }) {
  return (
    <div className="board">
      {board.map((row, r) =>
        row.map((cell, c) => (
          <div
            key={`${r}-${c}`}
            className={`cell player${cell}`}
            onClick={() => sendMove(c)}
          />
        ))
      )}
    </div>
  );
}
