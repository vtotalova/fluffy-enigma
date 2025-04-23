package main

import (
	"fluffy-enigma/api" // Import the api package
	"fluffy-enigma/db"  // Import the db package
)

func main() {
	db.Connect() // Connect to MongoDB
	api.Resp()   // Call the API response function
}
