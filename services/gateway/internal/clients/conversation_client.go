package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/internal/config"
	"gateway/internal/models"
	"net/http"
)

type ConversationClient struct {
	baseURL string
}

// Client da API de conversas
func NewConversationClient(cfg *config.Config) *ConversationClient {
	return &ConversationClient{
		baseURL: fmt.Sprintf("http://%s:%s", cfg.ConversationAPIHost, cfg.ConversationAPIPort),
	}
}

func (c *ConversationClient) SaveMessage(message models.Message) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.baseURL+"/conversations/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to save message, status: %d", resp.StatusCode)
	}

	return nil
}
