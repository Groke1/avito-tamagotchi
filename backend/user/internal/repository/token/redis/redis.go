package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	redis "github.com/redis/go-redis/v9"
)

const consumeRefreshTokenScript = `
local value = redis.call("GET", KEYS[1])
if not value then
	return false
end

redis.call("DEL", KEYS[1])
return value
`

type redisRepository struct {
	client redis.UniversalClient
}

func NewRedisRepository(client redis.UniversalClient) *redisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) AddToken(ctx context.Context, userID string, token entity.RefreshToken) error {
	token.UserID = userID
	payload, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("add refresh token: marshal: %w", err)
	}

	ttl := max(0, time.Until(token.ExpiresAt))

	ok, err := r.client.SetNX(ctx, r.sessionKey(token.TokenHash), payload, ttl).Result()
	if err != nil {
		return fmt.Errorf("add refresh token: %w", err)
	}
	if !ok {
		return errors.New("add refresh token: token hash already exists")
	}

	return nil
}

func (r *redisRepository) ConsumeRefreshToken(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	value, err := r.client.Eval(ctx, consumeRefreshTokenScript, []string{r.sessionKey(hash)}).Text()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, entity.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("consume refresh token: %w", err)
	}

	var token entity.RefreshToken
	if err = json.Unmarshal([]byte(value), &token); err != nil {
		return nil, fmt.Errorf("consume refresh token: unmarshal: %w", err)
	}

	return &token, nil
}

func (r *redisRepository) DeleteSession(ctx context.Context, userID, tokenHash string) error {
	key := r.sessionKey(tokenHash)

	value, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("delete session: %w", err)
	}

	var token entity.RefreshToken
	if err = json.Unmarshal(value, &token); err != nil {
		return fmt.Errorf("delete session: unmarshal: %w", err)
	}

	if token.UserID != userID {
		return nil
	}

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (r *redisRepository) sessionKey(tokenHash string) string {
	return "user:refresh_token:" + tokenHash
}
