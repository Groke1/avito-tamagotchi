package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/domain"
)

type UserClient struct {
	baseUrl string
	client  *http.Client
}

func NewUserClient(baseUrl string) *UserClient {
	return &UserClient{
		baseUrl: baseUrl,
		client:  &http.Client{Timeout: 3 * time.Second}}
}

func (uc *UserClient) WithdrawCoins(ctx context.Context, userID string, amount int) error {
	url := fmt.Sprintf("%s/update-coins", uc.baseUrl)

	body, _ := json.Marshal(UpdateCoinsRequest{UserID: userID, Delta: -amount})
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application-json")

	resp, err := uc.client.Do(req)
	if err != nil {
		return fmt.Errorf("user service unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var apiErr APIError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return fmt.Errorf("failed to decode api error: %v (status code %d)", err, resp.StatusCode)
	}

	switch resp.StatusCode {

	case http.StatusNotFound:
		return domain.ErrUserNotFound
	case http.StatusConflict:
		return domain.ErrNotEnoguhCoins
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("bad request to user service: %s", apiErr.Message)
	case http.StatusInternalServerError:
		return fmt.Errorf("user service internal error: %v", apiErr.Message)
	default:
		return fmt.Errorf("unexpected status %d from user service: %v", resp.StatusCode, apiErr.Message)
	}
}
