package api

import (
	"context"
	"encoding/json"
	"fluffy-enigma/db"
	"fluffy-enigma/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Prompt: "Create an APIServer that serves the ideas stored in the MongoDB database.
// The API should have a single endpoint that returns all ideas in JSON format.
// Use the Go programming language and the MongoDB driver for Go to implement this functionality.
// The API should handle errors gracefully and return appropriate HTTP status codes.
// Additionally, implement a function to fetch new ideas from an external API and store them
// in the database at regular intervals."

// StartAPIServer starts the HTTP server
func StartAPIServer() {
	r := mux.NewRouter()
	r.Use(corsMiddleware) // Use the CORS middleware
	r.HandleFunc("/api/ideas", getIdeasHandler).Methods("GET")
	r.HandleFunc("/api/ideas/{id}", deleteIdeaHandler).Methods("DELETE", "OPTIONS")

	log.Println("🌐 API Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("❌ HTTP server error: %v", err)
	}
}

// Add this new middleware function, referenced this problems: https://stackoverflow.com/questions/78090943/website-using-react-that-works-with-a-go-backend-cors-issue
// https://stackoverflow.com/questions/71594818/cors-issue-react-axios-frontend-and-golang-backend
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getIdeasHandler returns all stored ideas
func getIdeasHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all ideas in the database
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by creation date, newest first

	cursor, err := db.IdeaCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Println("❌ Failed to fetch ideas:", err)
		http.Error(w, "Failed to fetch ideas", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode the documents
	var ideas []models.Idea
	if err := cursor.All(ctx, &ideas); err != nil {
		log.Println("❌ Failed to decode ideas:", err)
		http.Error(w, "Failed to decode ideas", http.StatusInternalServerError)
		return
	}

	// Set headers and encode to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ideas); err != nil {
		log.Println("❌ Failed to encode ideas to JSON:", err)
		http.Error(w, "Failed to encode ideas", http.StatusInternalServerError)
		return
	}

	log.Printf("📤 Returned %d ideas", len(ideas))
}

// deleteIdeaHandler deletes an idea by ID
func deleteIdeaHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.IdeaCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Failed to delete idea", http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.Error(w, "Idea not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("🗑️ Deleted idea with ID: %s", id)

	// After deletion, try to fetch a new idea to maintain the count at 10
	go fetchIdeasUntilLimit()
}

// StartIdeaFetcher starts a goroutine to fetch ideas from the external API
func StartIdeaFetcher() {
	// Initial fetch to populate the database
	log.Println("🔄 Starting initial idea fetch")
	fetchIdeasUntilLimit()

	// Set up periodic fetching every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("🔄 Running scheduled idea fetch")
			fetchIdeasUntilLimit()
		}
	}()

	log.Println("✅ Idea fetcher started successfully")
}

// fetchIdeasUntilLimit fetches ideas until the database has 10 ideas
func fetchIdeasUntilLimit() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check current count
	count, err := db.IdeaCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Println("⚠️ Failed to count documents:", err)
		return
	}

	if count >= 10 {
		log.Println("✅ Already have 10 ideas in the database")
		return
	}

	// Calculate how many more ideas we need
	neededIdeas := 10 - count
	log.Printf("📊 Currently have %d ideas, need %d more", count, neededIdeas)

	// Fetch the required number of ideas
	for i := 0; i < int(neededIdeas); i++ {
		if err := fetchIdea(); err != nil {
			log.Printf("⚠️ Failed to fetch idea %d/%d: %v", i+1, neededIdeas, err)
			// Add a small delay before retrying
			time.Sleep(1 * time.Second)
			i-- // Retry this iteration
			continue
		}
		// Small delay between fetches to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}
}

// fetchIdea fetches a single idea from the external API and stores it
func fetchIdea() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch idea from API
	resp, err := http.Get("https://itsthisforthat.com/api.php?text")
	if err != nil {
		log.Println("❌ Failed to fetch idea:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("❌ Failed to read response body:", err)
		return err
	}

	// Parse "This for That"
	ideaText := string(body)
	log.Println("📥 Raw idea from API:", ideaText)
	parts := strings.SplitN(ideaText, " for ", 2)
	if len(parts) != 2 {
		log.Println("⚠️ Invalid idea format:", ideaText)
		return err
	}

	// Create the idea object with generated ObjectID
	idea := models.Idea{
		ID:        primitive.NewObjectID(),
		This:      strings.TrimSpace(parts[0]),
		That:      strings.TrimSpace(parts[1]),
		CreatedAt: time.Now(),
	}

	// Check for duplicates
	var existingIdea models.Idea
	err = db.IdeaCollection.FindOne(ctx, bson.M{
		"this": idea.This,
		"that": idea.That,
	}).Decode(&existingIdea)

	if err == nil {
		log.Println("⚠️ Duplicate idea found, skipping:", idea.This, "for", idea.That)
		return fetchIdea() // Try again recursively
	}

	// Insert the idea
	_, err = db.IdeaCollection.InsertOne(ctx, idea)
	if err != nil {
		log.Println("❌ Failed to insert idea:", err)
		return err
	}

	log.Printf("💡 New idea added: %s for %s", idea.This, idea.That)
	return nil
}
