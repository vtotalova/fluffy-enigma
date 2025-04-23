package api

import (
	"context"
	"fluffy-enigma/db"
	"fluffy-enigma/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

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
