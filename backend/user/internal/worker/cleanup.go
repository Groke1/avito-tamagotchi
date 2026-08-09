package worker

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Cleaner interface {
	DeleteExpiredTokens(ctx context.Context) error
}

func StartTokenCleaner(ctx context.Context, logger *zap.Logger, cleaner Cleaner, interval time.Duration) {
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
