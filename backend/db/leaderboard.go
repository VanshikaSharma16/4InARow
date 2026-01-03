package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type LeaderboardEntry struct {
	Player string `bson:"_id" json:"player"`
	Wins   int    `bson:"wins" json:"wins"`
}

func GetLeaderboard() ([]LeaderboardEntry, error) {
	if Client == nil {
		// MongoDB not connected, return empty leaderboard
		return []LeaderboardEntry{}, nil
	}

	collection := Client.Database("connect4").Collection("games")

	pipeline := []bson.M{
		// Filter out draws (winner == 0)
		{
			"$match": bson.M{
				"winner": bson.M{"$ne": 0},
			},
		},
		{
			"$project": bson.M{
				"winnerPlayer": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []interface{}{"$winner", 1}},
						"$player1",
						"$player2",
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":  "$winnerPlayer",
				"wins": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"wins": -1},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []LeaderboardEntry
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
