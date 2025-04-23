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

	r.HandleFunc("/api/ideas", getIdeasHandler).Methods("GET")
	r.HandleFunc("/api/ideas/{id}", deleteIdeaHandler).Methods("DELETE")

	log.Println("🌐 API Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("❌ HTTP server error: %v", err)
	}
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
}

// StartIdeaFetcher starts a goroutine to fetch ideas from the external API
func StartIdeaFetcher() {
	ticker := time.NewTicker(10 * time.Minute)

	// Run once at startup
	go fetchIdea()

	go func() {
		for range ticker.C {
			fetchIdea()
		}
	}()
}

// fetchIdea fetches a new idea from the external API and stores it in the database
func fetchIdea() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check count
	count, err := db.IdeaCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Println("⚠️ Failed to count documents:", err)
		return
	}
	if count >= 10 {
		log.Println("🚫 10 ideas already in DB. Skipping fetch.")
		return
	}

	// Fetch idea from API
	resp, err := http.Get("https://itsthisforthat.com/api.php?text")
	if err != nil {
		log.Println("❌ Failed to fetch idea:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("❌ Failed to read response body:", err)
		return
	}

	// Parse "This for That"
	parts := strings.SplitN(string(body), " for ", 2)
	if len(parts) != 2 {
		log.Println("⚠️ Invalid idea format:", string(body))
		return
	}

	idea := models.Idea{
		This:      strings.TrimSpace(parts[0]),
		That:      strings.TrimSpace(parts[1]),
		CreatedAt: time.Now(),
	}

	_, err = db.IdeaCollection.InsertOne(ctx, idea)
	if err != nil {
		log.Println("❌ Failed to insert idea:", err)
		return
	}

	log.Println("💡 New idea added:", idea.This, "for", idea.That)
}
