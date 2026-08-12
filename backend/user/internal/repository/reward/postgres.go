package reward

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/repository/converter"
	sqlcreward "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/reward/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type rewardRepository struct {
	db      sqlcreward.DBTX
	queries *sqlcreward.Queries
}

func NewRewardRepository(qdb sqlcreward.DBTX) *rewardRepository {
	return &rewardRepository{
		db:      qdb,
		queries: sqlcreward.New(qdb),
	}
}

func (r *rewardRepository) GetUserRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error) {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("get user rewards: %w", err)
	}
	rows, err := r.getQueries(ctx).GetUserRewardsByUserID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user rewards: %w", err)
	}

	result := make([]entity.UserReward, 0, len(rows))

	for _, row := range rows {
		result = append(result, entity.UserReward{
			ID:        row.ID.String(),
			UserID:    row.UserID.String(),
			PromoCode: row.PromoCode,
			Status:    entity.Status(row.Status),
			Definition: entity.RewardDefinition{
				Name:        row.Name,
				Description: row.Description,
			},
			EarnedReason: row.EarnedReason,
			ExpiresAt:    converter.TimestamptzToTime(row.ExpiresAt),
			RedeemedAt:   converter.TimestamptzToTime(row.RedeemedAt),
		})
	}

	return result, nil
}

func (r *rewardRepository) GetActiveRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error) {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("get user rewards: %w", err)
	}
	rows, err := r.getQueries(ctx).GetActiveRewardsByUserID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user rewards: %w", err)
	}

	result := make([]entity.UserReward, 0, len(rows))

	for _, row := range rows {
		result = append(result, entity.UserReward{
			ID:        row.ID.String(),
			UserID:    row.UserID.String(),
			PromoCode: row.PromoCode,
			Status:    entity.Status(row.Status),
			Definition: entity.RewardDefinition{
				Name:        row.Name,
				Description: row.Description,
			},
			EarnedReason: row.EarnedReason,
			ExpiresAt:    converter.TimestamptzToTime(row.ExpiresAt),
			RedeemedAt:   converter.TimestamptzToTime(row.RedeemedAt),
		})
	}

	return result, nil
}

func (r *rewardRepository) GetRewardsByUserIDAndPeriod(
	ctx context.Context, userID string, from, to time.Time,
) ([]entity.UserReward, error) {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("convert user id to uuid: %w", err)
	}

	rows, err := r.queries.GetRewardsByUserIDAndPeriod(
		ctx, sqlcreward.GetRewardsByUserIDAndPeriodParams{
			UserID:   userUUID,
			FromTime: converter.TimeToTimestamptz(&from),
			ToTime:   converter.TimeToTimestamptz(&to),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get rewards by user id and period: %w", err)
	}

	rewards := make([]entity.UserReward, 0, len(rows))

	for _, row := range rows {
		rewards = append(rewards, entity.UserReward{
			ID:        row.ID.String(),
			UserID:    row.UserID.String(),
			PromoCode: row.PromoCode,
			Status:    entity.Status(row.Status),

			Definition: entity.RewardDefinition{
				Name:                row.Name,
				Description:         row.Description,
				EarnedDescription:   row.EarnedDescription,
				RedeemedDescription: row.RedeemedDescription,
			},
			CreatedAt: *converter.TimestamptzToTime(row.CreatedAt),

			RedeemedAt: converter.TimestamptzToTime(row.RedeemedAt),
			ExpiresAt:  converter.TimestamptzToTime(row.ExpiresAt),
		})
	}

	return rewards, nil
}

func (r *rewardRepository) GetRewardByUserIDAndRewardID(ctx context.Context, userID, rewardID string) (*entity.UserReward, error) {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	rewardUUID, err := converter.StringToUUID(rewardID)
	if err != nil {
		return nil, fmt.Errorf("get reward by id: %w", err)
	}
	reward, err := r.getQueries(ctx).GetRewardByUserIDAndRewardID(ctx, sqlcreward.GetRewardByUserIDAndRewardIDParams{
		UserID:   userUUID,
		RewardID: rewardUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrRewardNotFound
		}

		return nil, fmt.Errorf("get user rewards: %w", err)
	}

	return &entity.UserReward{
		ID:        reward.ID.String(),
		UserID:    reward.UserID.String(),
		PromoCode: reward.PromoCode,
		Status:    entity.Status(reward.Status),
		Definition: entity.RewardDefinition{
			Name:        reward.Name,
			Description: reward.Description,
		},
		EarnedReason: reward.EarnedReason,
		ExpiresAt:    converter.TimestamptzToTime(reward.ExpiresAt),
		RedeemedAt:   converter.TimestamptzToTime(reward.RedeemedAt),
	}, nil
}

func (r *rewardRepository) GetRewardDefinitionByCode(ctx context.Context, code string) (*entity.RewardDefinition, error) {
	definitionRaw, err := r.getQueries(ctx).GetRewardDefinitionByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrRewardDefinitionNotFound
		}
		return nil, fmt.Errorf("get reward definition by code: %w", err)
	}

	return &entity.RewardDefinition{
		ID:          definitionRaw.ID,
		Code:        definitionRaw.Code,
		Name:        definitionRaw.Name,
		Description: definitionRaw.Description,
	}, nil
}

func (r *rewardRepository) GetRewardDefinitions(ctx context.Context) ([]entity.RewardDefinition, error) {
	definitionsRaw, err := r.getQueries(ctx).GetRewardDefinitions(ctx)
	if err != nil {
		return nil, fmt.Errorf("get reward definitions: %w", err)
	}
	definitions := make([]entity.RewardDefinition, 0, len(definitionsRaw))
	for _, definitionRaw := range definitionsRaw {
		definitions = append(definitions, entity.RewardDefinition{
			ID:          definitionRaw.ID,
			Code:        definitionRaw.Code,
			Name:        definitionRaw.Name,
			Description: definitionRaw.Description,
		})
	}
	return definitions, nil
}

func (r *rewardRepository) AddUserReward(ctx context.Context, userID, promoCode, earnedReason string,
	rewardID int32, expiresAt *time.Time) (*entity.UserReward, error) {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("add user rewards: %w", err)
	}

	reward, err := r.getQueries(ctx).AddUserReward(ctx, sqlcreward.AddUserRewardParams{
		UserID:       userUUID,
		RewardID:     rewardID,
		PromoCode:    promoCode,
		EarnedReason: earnedReason,
		ExpiresAt:    converter.TimeToTimestamptz(expiresAt),
	})
	if err == nil {
		return &entity.UserReward{
			ID:           reward.ID.String(),
			UserID:       reward.UserID.String(),
			PromoCode:    reward.PromoCode,
			Status:       entity.Status(reward.Status),
			EarnedReason: reward.EarnedReason,
			ExpiresAt:    converter.TimestamptzToTime(reward.ExpiresAt),
			RedeemedAt:   converter.TimestamptzToTime(reward.RedeemedAt),
		}, nil
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == "23505" &&
			pgErr.ConstraintName == "user_rewards_promo_code_key" {
			return nil, entity.ErrPromoCodeAlreadyExists
		}
	}

	return nil, fmt.Errorf("add user reward: %w", err)
}

func (r *rewardRepository) RedeemUserReward(ctx context.Context, userID, promoCode string) error {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return fmt.Errorf("redeem user rewards: %w", err)
	}

	_, err = r.getQueries(ctx).RedeemUserReward(ctx, sqlcreward.RedeemUserRewardParams{
		UserID:    userUUID,
		PromoCode: promoCode,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ErrRewardUnavailable
		}
		return fmt.Errorf("redeem user rewards: %w", err)
	}
	return nil
}

func (r *rewardRepository) getQueries(ctx context.Context) *sqlcreward.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return sqlcreward.New(tx)
}
