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
	ErrPetBadRequest  = errors.New("pet-service: bad request")
	ErrPetUnavailable = errors.New("pet-service: unavailable")
)

type PetServiceClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewPetServiceClient(baseURL string) *PetServiceClient {
	return &PetServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: config.HTTPClientTimeout,
		},
	}
}

func (c *PetServiceClient) UpdateXP(
	ctx context.Context,
	reqBody controller.UpdateXPRequest,
) error {
	url := c.baseURL + "/internal/update-xp"

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("marshal update xp request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewReader(jsonData),
	)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("create update xp request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("%w: %v", ErrPetUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(
			io.LimitReader(resp.Body, maxResponseBodySize),
		)

		return fmt.Errorf(
			"%w: status %d, body: %s",
			classifyPetStatus(resp.StatusCode),
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	return nil
}

func classifyPetStatus(status int) error {
	switch {
	case status >= 400 && status < 500:
		return ErrPetBadRequest
	default:
		return ErrPetUnavailable
	}
}
