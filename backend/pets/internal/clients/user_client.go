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

	var errResponse ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
		return fmt.Errorf("failed to decode api error: %w (status code %d)", err, resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return domain.ErrUserNotFound

	case http.StatusConflict:
		return domain.ErrNotEnoguhCoins

	case http.StatusUnprocessableEntity:
		return fmt.Errorf("bad request to user service: %s", errResponse.Message)

	case http.StatusInternalServerError:
		return fmt.Errorf("user service internal error: %v", errResponse.Message)

	default:
		return fmt.Errorf("unexpected status %d from user service: %v", resp.StatusCode, errResponse.Message)
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

func (uc *UserClient) ClaimReward(ctx context.Context, userID string, code string) (*domain.Reward, error) {
	url := uc.baseURL + "/rewards"

	body, _ := json.Marshal(struct {
		UserID string `json:"user_id"`
		Code   string `json:"code"`
	}{
		UserID: userID,
		Code:   code,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rewardResponse RewardResponse
	if err := json.NewDecoder(resp.Body).Decode(&rewardResponse); err != nil {
		return nil, err
	}

	reward := &domain.Reward{
		ID:           rewardResponse.RewardID,
		PromoCode:    rewardResponse.PromoCode,
		Name:         rewardResponse.Name,
		Description:  rewardResponse.Description,
		Status:       rewardResponse.Status,
		ExpiresAt:    rewardResponse.ExpiresAt,
		EarnedReason: rewardResponse.EarnedReason,
		RedeemedAt:   rewardResponse.RedeemedAt,
	}

	return reward, nil
}

func (uc *UserClient) GetRewardDescription(ctx context.Context, code string) (*domain.Reward, error) {
	url := fmt.Sprintf("%s/rewards/%s", uc.baseURL, code)

	body, _ := json.Marshal(struct{}{})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rewardResponse RewardDescriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&rewardResponse); err != nil {
		return nil, err
	}

	reward := &domain.Reward{
		PromoCode:   rewardResponse.Code,
		Name:        rewardResponse.Name,
		Description: rewardResponse.Description,
	}

	return reward, nil
}

func (uc *UserClient) NotifyActionDone(ctx context.Context, userID string, occuredAt time.Time) error {
	url := uc.baseURL + "/action"

	body, _ := json.Marshal(struct {
		UserID    string    `json:"user_id"`
		OccuredAt time.Time `json:"occurred_at"`
	}{
		UserID:    userID,
		OccuredAt: occuredAt,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	var errResponse ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
		return fmt.Errorf("failed to decode api error: %w (status code %d)", err, resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("bad request to user service: %s", errResponse.Message)

	case http.StatusInternalServerError:
		return fmt.Errorf("user service internal error: %v", errResponse.Message)

	default:
		return fmt.Errorf("unexpected status %d from user service: %v", resp.StatusCode, errResponse.Message)
	}
}
