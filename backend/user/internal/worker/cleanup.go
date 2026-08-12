package worker

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type (
	TokenCleaner interface {
		DeleteExpiredTokens(ctx context.Context) error
	}

	EventsCleaner interface {
		DeleteDeliveredEvents(ctx context.Context) error
	}
)

func StartTokenCleaner(ctx context.Context, logger *zap.Logger, cleaner TokenCleaner, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleaner.DeleteExpiredTokens(ctx); err != nil {
				logger.Error("failed to clean expired tokens:", zap.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}
}

func StartEventsCleaner(ctx context.Context, logger *zap.Logger, cleaner EventsCleaner, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleaner.DeleteDeliveredEvents(ctx); err != nil {
				logger.Error("failed to clean delivered events:", zap.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}
}
