package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/adapter/httpclient"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewPetClient(addr string, httpClient *http.Client) (*client, error) {
	addr = httpclient.NormalizeAddr(addr)
	if addr == "" {
		return nil, errors.New("pet service address is empty")
	}

	baseURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse pet service address: %w", err)
	}

	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

func (c *client) SendDailyBonus(ctx context.Context, userID string, streak, coins int32) error {
	payload := dailyBonusRequest{
		UserID: userID,
		Streak: streak,
		Coins:  coins,
	}

	endpoint := c.baseURL.ResolveReference(
		&url.URL{
			Path: "/internal/daily-bonus",
		},
	)

	_, err := httpclient.WithRequest[dailyBonusRequest, emptyResponse](
		ctx, c.httpClient, endpoint.String(), http.MethodPost, payload)
	if err != nil {
		return fmt.Errorf("send daily bonus: %w", err)
	}

	return nil
}

func (c *client) GetPetDailyStat(ctx context.Context, userID string) (*entity.PetStat, error) {
	payload := petDailyStatRequest{UserID: userID}

	endpoint := c.baseURL.ResolveReference(
		&url.URL{
			Path: "/internal/daily-gained-xp",
		},
	)

	resp, err := httpclient.WithRequest[petDailyStatRequest, entity.PetStat](
		ctx, c.httpClient, endpoint.String(), http.MethodGet, payload)
	if err != nil {
		return nil, fmt.Errorf("get pet stat: %w", err)
	}

	return resp, nil
}
