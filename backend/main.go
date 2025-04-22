package main

import (
	"fluffy-enigma/db" // Import the db package		
	"fluffy-enigma/api" // Import the api package
)

func main() { 
	db.Connect() // Connect to MongoDB
	api.Resp() // Call the API response function
}
