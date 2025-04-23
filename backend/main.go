package main

import (
	"fluffy-enigma/api" // Import the api package
	"fluffy-enigma/db"  // Import the db package
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db.Connect()           // Connect to MongoDB
	api.StartIdeaFetcher() // Call the API response function

	// Start the API server in a goroutine
	go api.StartAPIServer()

	// Set up channel to listen for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("🚀 Application running. Press Ctrl+C to exit.")

	// Block until we receive a signal
	<-sigChan
	log.Println("Shutting down gracefully...")
}
