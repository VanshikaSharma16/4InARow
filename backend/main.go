package main

import (
	"encoding/json"
	"log"
	"net/http"

	"connect4/db"
	"connect4/websocket"
)

func main() {
	db.ConnectMongo()

	http.HandleFunc("/ws", websocket.HandleWS)
	http.HandleFunc("/leaderboard", withCORS(leaderboardHandler))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := db.GetLeaderboard()
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		handler(w, r)
	}
}
