package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func SaveGameResult(player1, player2 string, winner int) error {
	if Client == nil {
		// MongoDB not connected, silently skip saving
		return nil
	}

	collection := Client.Database("connect4").Collection("games")

	// Only save if there's a winner (not a draw or forfeit)
	// For draws, winner is 0, we still save but mark it differently
	doc := bson.M{
		"player1":   player1,
		"player2":   player2,
		"winner":    winner,
		"isDraw":    winner == 0,
		"createdAt": time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), doc)
	return err
}
