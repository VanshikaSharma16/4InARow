package db

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client
var Client *mongo.Client

// loadEnvFile loads environment variables from .env file if it exists
func loadEnvFile() {
	envFile, err := os.Open(".env")
	if err != nil {
		// .env file doesn't exist, that's okay
		return
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}
			os.Setenv(key, value)
		}
	}
}

func ConnectMongo() {
	// Try to load .env file if it exists
	loadEnvFile()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Println("⚠️  WARNING: MONGODB_URI environment variable is not set.")
		log.Println("   Game results will not be saved, but the server will still run.")
		log.Println("   To enable MongoDB:\n" +
			"   1. Set it in your terminal: export MONGODB_URI=\"your-connection-string\"\n" +
			"   2. Create a .env file in the backend directory with: MONGODB_URI=your-connection-string")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("⚠️  WARNING: Failed to connect to MongoDB: %v", err)
		log.Println("   Game results will not be saved, but the server will still run.")
		log.Println("   Check your MONGODB_URI and network connection.")
		return
	}

	// Force authentication + connectivity check
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("⚠️  WARNING: Failed to ping MongoDB: %v", err)
		log.Println("   Game results will not be saved, but the server will still run.")
		log.Println("   Check your MongoDB connection string and IP whitelist settings.")
		return
	}

	Client = client
	log.Println("✅ MongoDB connected successfully")
}
