package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestInsert() {
	collection := Client.Database("connect4").Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"winner":    "Player1",
		"createdAt": time.Now(),
		"reason":    "Test win insert",
	}

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Game result saved with ID:", result.InsertedID)
}
