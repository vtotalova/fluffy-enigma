package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Idea represents a startup idea with unique ID
type Idea struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	This      string             `bson:"this" json:"this"`
	That      string             `bson:"that" json:"that"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// APIResponse represents the response from itsthisforthat.com
type APIResponse struct {
	This string `json:"this"`
	That string `json:"that"`
}
