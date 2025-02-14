package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/internal/config"
	"net/http"
)

type AddressClient struct {
	baseURL string
}

// Client da api de endereco
// mesmo que nao usemos agora, pode ser util no futuro
func NewAddressClient(cfg *config.Config) *AddressClient {
	return &AddressClient{
		baseURL: fmt.Sprintf("http://%s:%s", cfg.AddressAPIHost, cfg.AddressAPIPort),
	}
}

func (c *AddressClient) SaveAddress(address interface{}) error {
	jsonData, err := json.Marshal(address)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.baseURL+"/addresses/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create address, status: %d", resp.StatusCode)
	}

	return nil
}
