package main

import (
	"log"
	"net/http"
	"user-api/internal/clients"
	"user-api/internal/config"
	"user-api/internal/database"
	"user-api/internal/handlers"
)

func main() {
	// Configurando os logs
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Load configuration
	cfg := config.Load()

	// Initialize database
	_, err := database.InitDB(
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize address client
	addressClient := clients.NewAddressClient(cfg)
	handlers.SetAddressClient(addressClient)

	// Initialize router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/users/", handlers.HandleUsers)
	mux.HandleFunc("/users/cpf/", handlers.HandleUsers)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
