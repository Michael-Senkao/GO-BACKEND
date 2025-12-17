package main

import (
	"log"

	"task_manager/data"
	"task_manager/router"
)

func main() {


	// -----------------------------
	// Connect to MongoDB
	// -----------------------------
	data.ConnectDB()

	// Initialize collections
	data.InitUserCollection()

	log.Println("Database and collections initialized")

	// -----------------------------
	// Setup Router
	// -----------------------------
	r := router.SetupRouter()

	// -----------------------------
	// Start Server
	// -----------------------------
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
