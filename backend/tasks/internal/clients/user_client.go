package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UpdateCoinsRequest struct {
	UserID string `json:"user_id"`
	Delta  int64  `json:"delta"`
}

type UpdateCoinsResponse struct {
	UserID string `json:"user_id"`
	Coins  int64  `json:"coins"`
}

type UserServiceClient struct {
	httpClient *http.Client
	baseURL    string // Например: "http://user-service:8080" (из docker-compose или DNS)
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *UserServiceClient) UpdateCoins(ctx context.Context, userID string, delta int64) (*UpdateCoinsResponse, error) {
	url := fmt.Sprintf("%s/internal/update-coins", c.baseURL)

	reqBody := UpdateCoinsRequest{
		UserID: userID,
		Delta:  delta,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Здесь можно распарсить ErrorBody из OpenAPI, если сервис вернул ошибку
		return nil, fmt.Errorf("user-service returned bad status: %d", resp.StatusCode)
	}

	var respBody UpdateCoinsResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &respBody, nil
}
