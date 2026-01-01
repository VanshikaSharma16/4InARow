import { useEffect, useState } from "react";

function Leaderboard() {
  const [data, setData] = useState([]);

  useEffect(() => {
    fetch("http://localhost:8080/leaderboard")
      .then(res => res.json())
      .then(setData);
  }, []);

  return (
    <div>
      <h2>Leaderboard</h2>
      <ul>
        {data.map((item, index) => (
          <li key={index}>
            {item.player} — {item.wins} wins
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Leaderboard;
