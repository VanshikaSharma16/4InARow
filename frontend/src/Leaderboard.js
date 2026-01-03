import { useEffect, useState } from "react";

function Leaderboard({ refreshTrigger = 0 }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchLeaderboard = () => {
    setLoading(true);
    setError(null);

    fetch("https://connect4-backend-a4kq.onrender.com/leaderboard")
      .then(res => {
        if (!res.ok) {
          throw new Error("Failed to fetch leaderboard");
        }
        return res.json();
      })
      .then(result => {
        if (Array.isArray(result)) {
          setData(result);
        } else {
          setData([]);
        }
      })
      .catch(err => {
        console.error("Leaderboard fetch error:", err);
        setError("Unable to load leaderboard. Make sure the backend server is running.");
        setData([]);
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchLeaderboard();
  }, []);

  useEffect(() => {
    if (refreshTrigger > 0) {
      const timer = setTimeout(fetchLeaderboard, 500);
      return () => clearTimeout(timer);
    }
  }, [refreshTrigger]);

  return (
    <div>
      <h2>Leaderboard</h2>

      {loading && <p>Loading...</p>}

      {error && (
        <p style={{ color: "#999", fontStyle: "italic" }}>
          {error}
        </p>
      )}

      {!loading && !error && data.length === 0 && (
        <p style={{ color: "#999", fontStyle: "italic" }}>
          No games played yet.
        </p>
      )}

      {!loading && !error && data.length > 0 && (
        <ul>
          {data.map((item, index) => (
            <li key={index}>
              {item.player} — {item.wins} wins
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default Leaderboard;
