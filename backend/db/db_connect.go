package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo" //https://pkg.go.dev/go.mongodb.org/mongo-driver#readme-requirements
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var IdeaCollection *mongo.Collection

func Connect() {
	// Fetch the Mongo URI from environment variable
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("⚠️ MONGO_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed:", err)
	}

	Client = client
	IdeaCollection = client.Database("fluffydb").Collection("ideas")
	log.Println("✅ Connected to MongoDB")
}
