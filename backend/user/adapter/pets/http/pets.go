package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

type DailyBonusRequest struct {
	UserID    string `json:"user_id"`
	Streak    int32  `json:"streak"`
	BonusDate string `json:"bonus_date"`
}

func NewPetClient(addr string, httpClient *http.Client) (*client, error) {
	addr = normalizeAddr(addr)
	if addr == "" {
		return nil, errors.New("pet service address is empty")
	}

	baseURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse pet service address: %w", err)
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 2 * time.Second,
		}
	}

	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

func (c *client) SendDailyBonus(
	ctx context.Context,
	userID string,
	streak int32,
	bonusDate time.Time,
) error {
	payload := DailyBonusRequest{
		UserID:    userID,
		Streak:    streak,
		BonusDate: bonusDate.Format(time.DateOnly),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal daily bonus request: %w", err)
	}

	endpoint := c.baseURL.ResolveReference(
		&url.URL{Path: "/internal/daily-bonus"},
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint.String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create daily bonus request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	//resp, err := c.httpClient.Do(req)
	//if err != nil {
	//	return fmt.Errorf("send daily bonus request: %w", err)
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode < http.StatusOK ||
	//	resp.StatusCode >= http.StatusMultipleChoices {
	//	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
	//
	//	return fmt.Errorf(
	//		"pet service returned status %d: %s",
	//		resp.StatusCode,
	//		strings.TrimSpace(string(responseBody)),
	//	)
	//}

	return nil
}

func normalizeAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	addr = strings.TrimRight(addr, "/")

	if addr == "" {
		return ""
	}

	if strings.HasPrefix(addr, "http://") ||
		strings.HasPrefix(addr, "https://") {
		return addr
	}

	return "http://" + addr
}
