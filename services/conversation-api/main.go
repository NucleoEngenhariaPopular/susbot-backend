package main

import (
	"conversation-api/internal/config"
	"conversation-api/internal/database"
	"conversation-api/internal/handlers"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cfg := config.Load()

	client, err := database.InitMongoDB(
		cfg.MongoURI,
		cfg.MongoDBName,
		cfg.MongoCollection,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	mux := http.NewServeMux()

	mux.HandleFunc("/conversations/", handlers.HandleConversations)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
