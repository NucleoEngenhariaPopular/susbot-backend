package main

import (
	"log"
	"net/http"

	"gateway/internal/config"
	"gateway/internal/handlers"
	"gateway/internal/utils"
)

func main() {
	// Configurando os logs
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Carregando as configs
	cfg := config.Load()

	// Check API connections
	utils.CheckAPIConnections(cfg)

	// Criando o Handler com os clients
	handler := handlers.NewHandler(cfg)

	mux := http.NewServeMux()

	// Apenas uma rota para lidar com as mensagens recebidas da twilio
	mux.HandleFunc("POST /", handler.HandlePost)

	// Iniciando o server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

