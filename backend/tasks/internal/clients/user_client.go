package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/config"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
)

var (
	ErrUserNotFound           = errors.New("user-service: user not found")
	ErrBadRequest             = errors.New("user-service: bad request")
	ErrUnavailable            = errors.New("user-service: unavailable")
	maxResponseBodySize int64 = 1 << 20
)

type UserServiceClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: config.HTTPClientTimeout,
		},
	}
}

func (c *UserServiceClient) UpdateCoins(ctx context.Context, reqBody controller.UpdateCoinsRequest) (*controller.UpdateCoinsResponse, error) {
	url := c.baseURL + "/internal/update-coins"

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	limitedBody := io.LimitReader(resp.Body, maxResponseBodySize)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(limitedBody)
		return nil, fmt.Errorf("%w: status %d, body: %s",
			classifyStatus(resp.StatusCode), resp.StatusCode, string(bodyBytes))
	}

	var respBody controller.UpdateCoinsResponse
	if err := json.NewDecoder(limitedBody).Decode(&respBody); err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &respBody, nil
}

func (c *UserServiceClient) NotifyActionDone(ctx context.Context, reqBody controller.NotifyActionRequest) error {
	url := c.baseURL + "/internal/action"

	body, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: status %d, body: %s",
			classifyStatus(resp.StatusCode), resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func classifyStatus(status int) error {
	switch {
	case status == http.StatusNotFound:
		return ErrUserNotFound
	case status >= 400 && status < 500:
		return ErrBadRequest
	default:
		return ErrUnavailable
	}
}
