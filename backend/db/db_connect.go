package db
import (
 	"context"
    "fmt"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo" //https://pkg.go.dev/go.mongodb.org/mongo-driver#readme-requirements
    "go.mongodb.org/mongo-driver/mongo/options"
   )

func Connect() { 	// MongoDB URI (replace with your connection string)
	//Copilot Promt: "Write a Go program that connects to a MongoDB database and performs a simple operation."
    uri := "mongodb://localhost:27017"

    // Set client options
    clientOptions := options.Client().ApplyURI(uri)

    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Connect to MongoDB
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Ping the database to test the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Could not connect to MongoDB:", err)
    }

    fmt.Println("Connected to MongoDB!")

    // Example: Get a reference to a collection
    collection := client.Database("testdb").Collection("testcollection")
    fmt.Printf("Using collection: %v\n", collection.Name())

    // Close the connection when you're done
    err = client.Disconnect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connection to MongoDB closed.")
}