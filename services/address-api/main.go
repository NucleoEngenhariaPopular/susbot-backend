package main

import (
	"address-api/internal/config"
	"address-api/internal/database"
	"address-api/internal/handlers"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	// Inicializar BD
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

	// Inicializar o router
	mux := http.NewServeMux()

	// Registrando os Handlers
	mux.HandleFunc("/ubs/", handlers.HandleUBS)
	mux.HandleFunc("/teams/", handlers.HandleTeams)
	mux.HandleFunc("/streets/", handlers.HandleStreetSegments)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
