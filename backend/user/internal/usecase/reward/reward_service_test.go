package reward

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func TestGetAllRewards_RecomputesStatusPerItem(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	repo := &fakeRewardRepository{
		getUserRewardsByUserIDFn: func(ctx context.Context, userID string) ([]entity.UserReward, error) {
			return []entity.UserReward{
				{ID: "1", Status: entity.StatusActive, ExpiresAt: &past},
				{ID: "2", Status: entity.StatusRedeemed, ExpiresAt: &past},
				{ID: "3", Status: entity.StatusActive, ExpiresAt: nil},
			}, nil
		},
	}
	s := NewRewardService(repo)

	rewards, err := s.GetAllRewards(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rewards) != 3 {
		t.Fatalf("expected 3 rewards, got %d", len(rewards))
	}
	if rewards[0].Status != entity.StatusExpired {
		t.Errorf("expected reward 1 to be expired, got %v", rewards[0].Status)
	}
	if rewards[1].Status != entity.StatusRedeemed {
		t.Errorf("expected reward 2 to stay redeemed, got %v", rewards[1].Status)
	}
	if rewards[2].Status != entity.StatusActive {
		t.Errorf("expected reward 3 to stay active, got %v", rewards[2].Status)
	}
}

func TestGetAllRewards_RepositoryError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &fakeRewardRepository{
		getUserRewardsByUserIDFn: func(ctx context.Context, userID string) ([]entity.UserReward, error) {
			return nil, repoErr
		},
	}
	s := NewRewardService(repo)

	_, err := s.GetAllRewards(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestGetActiveRewards_PassesThrough(t *testing.T) {
	want := []entity.UserReward{{ID: "1"}}
	repo := &fakeRewardRepository{
		getActiveRewardsByUserIDFn: func(ctx context.Context, userID string) ([]entity.UserReward, error) { return want, nil },
	}
	s := NewRewardService(repo)

	got, err := s.GetActiveRewards(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestGetActiveRewards_RepositoryError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &fakeRewardRepository{
		getActiveRewardsByUserIDFn: func(ctx context.Context, userID string) ([]entity.UserReward, error) { return nil, repoErr },
	}
	s := NewRewardService(repo)

	_, err := s.GetActiveRewards(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestGetReward_RecomputesStatus(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	repo := &fakeRewardRepository{
		getRewardByUserIDAndRewardIDFn: func(ctx context.Context, userID, rewardID string) (*entity.UserReward, error) {
			return &entity.UserReward{ID: rewardID, Status: entity.StatusActive, ExpiresAt: &past}, nil
		},
	}
	s := NewRewardService(repo)

	got, err := s.GetReward(context.Background(), "user-1", "reward-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != entity.StatusExpired {
		t.Errorf("expected expired status, got %v", got.Status)
	}
}

func TestGetReward_NotFound(t *testing.T) {
	repo := &fakeRewardRepository{
		getRewardByUserIDAndRewardIDFn: func(ctx context.Context, userID, rewardID string) (*entity.UserReward, error) {
			return nil, entity.ErrRewardNotFound
		},
	}
	s := NewRewardService(repo)

	_, err := s.GetReward(context.Background(), "user-1", "missing")
	if !errors.Is(err, entity.ErrRewardNotFound) {
		t.Fatalf("expected ErrRewardNotFound, got %v", err)
	}
}

func TestRedeemReward_Success(t *testing.T) {
	repo := &fakeRewardRepository{
		redeemUserRewardFn: func(ctx context.Context, userID, promoCode string) error { return nil },
	}
	s := NewRewardService(repo)

	if err := s.RedeemReward(context.Background(), "user-1", "PROMO-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRedeemReward_Unavailable(t *testing.T) {
	repo := &fakeRewardRepository{
		redeemUserRewardFn: func(ctx context.Context, userID, promoCode string) error { return entity.ErrRewardUnavailable },
	}
	s := NewRewardService(repo)

	err := s.RedeemReward(context.Background(), "user-1", "PROMO-1")
	if !errors.Is(err, entity.ErrRewardUnavailable) {
		t.Fatalf("expected ErrRewardUnavailable, got %v", err)
	}
}

func TestGetDefinition_Success(t *testing.T) {
	repo := &fakeRewardRepository{
		getRewardDefinitionByCodeFn: func(ctx context.Context, code string) (*entity.RewardDefinition, error) {
			return &entity.RewardDefinition{Code: code, Name: "Скидка"}, nil
		},
	}
	s := NewRewardService(repo)

	def, err := s.GetDefinition(context.Background(), "DELIVERY_DISCOUNT_10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if def.Code != "DELIVERY_DISCOUNT_10" {
		t.Fatalf("unexpected code: %q", def.Code)
	}
}

func TestGetDefinition_NotFound(t *testing.T) {
	repo := &fakeRewardRepository{
		getRewardDefinitionByCodeFn: func(ctx context.Context, code string) (*entity.RewardDefinition, error) {
			return nil, entity.ErrRewardDefinitionNotFound
		},
	}
	s := NewRewardService(repo)

	_, err := s.GetDefinition(context.Background(), "UNKNOWN")
	if !errors.Is(err, entity.ErrRewardDefinitionNotFound) {
		t.Fatalf("expected ErrRewardDefinitionNotFound, got %v", err)
	}
}

func TestGrantReward_SuccessOnFirstAttempt(t *testing.T) {
	repo := &fakeRewardRepository{
		getRewardDefinitionByCodeFn: func(ctx context.Context, code string) (*entity.RewardDefinition, error) {
			return &entity.RewardDefinition{ID: 42, Code: code}, nil
		},
		addUserRewardFn: func(ctx context.Context, userID, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error) {
			return &entity.UserReward{ID: "reward-1", PromoCode: promoCode}, nil
		},
	}
	s := NewRewardService(repo)

	reward, err := s.GrantReward(context.Background(), "user-1", "DELIVERY_DISCOUNT_10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reward.Definition.Code != "DELIVERY_DISCOUNT_10" {
		t.Errorf("expected definition to be attached to reward, got %+v", reward.Definition)
	}
	if len(repo.addUserRewardCalls) != 1 {
		t.Fatalf("expected exactly one AddUserReward attempt, got %d", len(repo.addUserRewardCalls))
	}
}

func TestGrantReward_RetriesOnPromoCodeCollisionThenSucceeds(t *testing.T) {
	attempts := 0
	repo := &fakeRewardRepository{
		getRewardDefinitionByCodeFn: func(ctx context.Context, code string) (*entity.RewardDefinition, error) {
			return &entity.RewardDefinition{ID: 1, Code: code}, nil
		},
		addUserRewardFn: func(ctx context.Context, userID, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error) {
			attempts++
			if attempts < 3 {
				return nil, entity.ErrPromoCodeAlreadyExists
			}
			return &entity.UserReward{ID: "reward-1", PromoCode: promoCode}, nil
		},
	}
	s := NewRewardService(repo)

	reward, err := s.GrantReward(context.Background(), "user-1", "CODE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reward == nil {
		t.Fatalf("expected a reward to be returned")
	}
	if attempts != 3 {
		t.Fatalf("expected exactly 3 attempts, got %d", attempts)
	}
}

func TestGrantReward_ExhaustsRetriesOnPersistentCollision(t *testing.T) {
	repo := &fakeRewardRepository{
		getRewardDefinitionByCodeFn: func(ctx context.Context, code string) (*entity.RewardDefinition, error) {
			return &entity.RewardDefinition{ID: 1, Code: code}, nil
		},
		addUserRewardFn: func(ctx context.Context, userID, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error) {
			return nil, entity.ErrPromoCodeAlreadyExists
		},
	}
	s := NewRewardService(repo)

	_, err := s.GrantReward(context.Background(), "user-1", "CODE")
	if err == nil {
		t.Fatalf("expected an error after exhausting all attempts")
	}
	if len(repo.addUserRewardCalls) != maxPromoCodeGenerationAttempts {
		t.Fatalf("expected %d attempts, got %d", maxPromoCodeGenerationAttempts, len(repo.addUserRewardCalls))
	}
}

func TestGrantReward_NonCollisionError_FailsImmediately(t *testing.T) {
	repoErr := errors.New("constraint violation")
	repo := &fakeRewardRepository{
		getRewardDefinitionByCodeFn: func(ctx context.Context, code string) (*entity.RewardDefinition, error) {
			return &entity.RewardDefinition{ID: 1, Code: code}, nil
		},
		addUserRewardFn: func(ctx context.Context, userID, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error) {
			return nil, repoErr
		},
	}
	s := NewRewardService(repo)

	_, err := s.GrantReward(context.Background(), "user-1", "CODE")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
	if len(repo.addUserRewardCalls) != 1 {
		t.Fatalf("expected exactly 1 attempt for a non-collision error, got %d", len(repo.addUserRewardCalls))
	}
}

func TestGrantReward_DefinitionNotFound_NoAttempts(t *testing.T) {
	repo := &fakeRewardRepository{
		getRewardDefinitionByCodeFn: func(ctx context.Context, code string) (*entity.RewardDefinition, error) {
			return nil, entity.ErrRewardDefinitionNotFound
		},
		addUserRewardFn: func(ctx context.Context, userID, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error) {
			t.Fatalf("AddUserReward should not be called when definition lookup fails")
			return nil, nil
		},
	}
	s := NewRewardService(repo)

	_, err := s.GrantReward(context.Background(), "user-1", "UNKNOWN")
	if !errors.Is(err, entity.ErrRewardDefinitionNotFound) {
		t.Fatalf("expected ErrRewardDefinitionNotFound, got %v", err)
	}
}
