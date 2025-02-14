package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/internal/config"
	"net/http"
)

type UserClient struct {
	baseURL string
}

// Client dos users
func NewUserClient(cfg *config.Config) *UserClient {
	return &UserClient{
		baseURL: fmt.Sprintf("http://%s:%s", cfg.UserAPIHost, cfg.UserAPIPort),
	}
}

// Envia para a user-api
func (c *UserClient) SaveUser(user interface{}) error {
	// Cria o pacote json
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Faz a requisicao para a URL certa
	resp, err := http.Post(c.baseURL+"/users/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create user, status: %d", resp.StatusCode)
	}

	return nil
}

// Encontrar o usuario por telefone
func (c *UserClient) GetUserByPhone(phone string) (interface{}, error) {
	resp, err := http.Get(c.baseURL + "/users/phone/" + phone)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var result struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
