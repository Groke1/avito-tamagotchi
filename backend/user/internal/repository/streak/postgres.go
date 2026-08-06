package streak

import (
	"context"
	"errors"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/repository/converter"
	sqlcstreak "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/streak/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type streakRepository struct {
	db      sqlcstreak.DBTX
	queries *sqlcstreak.Queries
}

func NewStreakRepository(qdb sqlcstreak.DBTX) *streakRepository {
	return &streakRepository{
		db:      qdb,
		queries: sqlcstreak.New(qdb),
	}
}

func (r *streakRepository) GetStreakByUserIDForUpdate(ctx context.Context, userID string) (*entity.Streak, error) {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return &entity.Streak{}, fmt.Errorf("get streak by user id: parse user id: %w", err)
	}
	userStreak, err := r.getQueries(ctx).GetStreakByUserIDForUpdate(ctx, userUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}

		return nil, fmt.Errorf("get streak token by user_id: %w", err)
	}

	return &entity.Streak{
		UserID:         userStreak.UserID.String(),
		CurrentStreak:  userStreak.CurrentStreak,
		LastActiveDate: userStreak.LastActiveDate.Time,
	}, nil
}

func (r *streakRepository) UpdateStreak(ctx context.Context, streak *entity.Streak) error {
	userUUID, err := converter.StringToUUID(streak.UserID)
	if err != nil {
		return fmt.Errorf("update streak: parse user id: %w", err)
	}
	err = r.getQueries(ctx).UpdateStreak(ctx, sqlcstreak.UpdateStreakParams{
		UserID:        userUUID,
		CurrentStreak: streak.CurrentStreak,
		LastActiveDate: pgtype.Date{
			Time:             streak.LastActiveDate,
			Valid:            true,
			InfinityModifier: pgtype.Finite},
	})

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == "23503" && pgErr.ConstraintName == "user_streaks_user_id_key" {
			return entity.ErrUserNotFound
		}
	}
	if err != nil {
		return fmt.Errorf("update streak: %w", err)
	}

	return nil
}

func (r *streakRepository) getQueries(ctx context.Context) *sqlcstreak.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return sqlcstreak.New(tx)
}
