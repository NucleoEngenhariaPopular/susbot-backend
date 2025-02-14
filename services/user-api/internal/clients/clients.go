package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"user-api/internal/config"
	"user-api/internal/models"
)

type AddressClient struct {
	baseURL string
}

func NewAddressClient(cfg *config.Config) *AddressClient {
	return &AddressClient{
		baseURL: fmt.Sprintf("http://%s:%s", cfg.AddressAPIHost, cfg.AddressAPIPort),
	}
}

type addressAPIResponse struct {
	Success bool       `json:"success"`
	Data    searchData `json:"data"`
	Error   string     `json:"error"`
}

type searchData struct {
	Team struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		UBS  struct {
			Name string `json:"name"`
		} `json:"ubs"`
	} `json:"team"`
}

func (c *AddressClient) GetTeamInfo(streetName, number, city, state string) (*models.TeamInfo, error) {
	// Build query URL
	query := url.Values{}
	query.Add("street", streetName)
	query.Add("number", number)
	query.Add("city", city)
	query.Add("state", state)

	url := fmt.Sprintf("%s/streets/search?%s", c.baseURL, query.Encode())

	// Make request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to address API: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No team found for this address
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("address API returned status %d", resp.StatusCode)
	}

	// Parse response
	var apiResp addressAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("address API error: %s", apiResp.Error)
	}

	// Map to TeamInfo
	teamInfo := &models.TeamInfo{
		ID:      apiResp.Data.Team.ID,
		Name:    apiResp.Data.Team.Name,
		UBSName: apiResp.Data.Team.UBS.Name,
	}

	return teamInfo, nil
}

