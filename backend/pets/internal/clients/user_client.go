package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

type UserClient struct {
	baseURL string
	client  *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		//nolint:mnd // таймаут на весь запрос 3 секунды
		client: &http.Client{Timeout: 3 * time.Second}}
}

func (uc *UserClient) WithdrawCoins(ctx context.Context, userID string, amount int) error {
	if amount == 0 {
		return nil
	}

	url := uc.baseURL + "/update-coins"

	body, _ := json.Marshal(UpdateCoinsRequest{UserID: userID, Delta: -amount})
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application-json")

	resp, err := uc.client.Do(req)
	if err != nil {
		return fmt.Errorf("user service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var apiErr APIError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return fmt.Errorf("failed to decode api error: %w (status code %d)", err, resp.StatusCode)
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

func (uc *UserClient) GetUsernamesByIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	url := uc.baseURL + "/usernames"

	body, _ := json.Marshal(struct {
		UserIDs []string `json:"user_ids"`
	}{
		UserIDs: userIDs,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userNames UserNamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&userNames); err != nil {
		return nil, err
	}

	userMap := make(map[string]string, len(userNames.Users))
	for _, user := range userNames.Users {
		userMap[user.ID] = user.UserName
	}

	return userMap, nil
}
