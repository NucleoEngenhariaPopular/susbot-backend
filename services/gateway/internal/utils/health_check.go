package utils

import (
	"fmt"
	"gateway/internal/config"
	"log"
	"net/http"
	"time"
)

func CheckAPIConnections(cfg *config.Config) {
	apis := map[string]string{
		"User API":         fmt.Sprintf("http://%s:%s/users/1", cfg.UserAPIHost, cfg.UserAPIPort),
		"Address API":      fmt.Sprintf("http://%s:%s/health", cfg.AddressAPIHost, cfg.AddressAPIPort),
		"Conversation API": fmt.Sprintf("http://%s:%s/health", cfg.ConversationAPIHost, cfg.ConversationAPIPort),
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for name, url := range apis {
		log.Printf("Checking connection to %s...", name)
		_, err := client.Get(url)
		if err != nil {
			log.Printf("⚠️ Warning: Could not connect to %s: %v", name, err)
		} else {
			log.Printf("✅ Successfully connected to %s", name)
		}
	}
}

