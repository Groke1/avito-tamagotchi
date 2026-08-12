package event

import (
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type listEventsResponse struct {
	UserID string          `json:"user_id"`
	Items  []eventResponse `json:"items"`
}

type eventResponse struct {
	Type      string          `json:"type"`
	XP        *int32          `json:"xp,omitempty"`
	Coins     *int32          `json:"coins,omitempty"`
	Streak    *int32          `json:"streak,omitempty"`
	Reward    *rewardResponse `json:"reward,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type rewardResponse struct {
	RewardID     string     `json:"reward_id"`
	PromoCode    string     `json:"promo_code"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	EarnedReason string     `json:"earned_reason"`
	ExpiresAt    *time.Time `json:"expires_at"`
	RedeemedAt   *time.Time `json:"redeemed_at"`
}

func toEventResponse(event *entity.UserEventDetails) eventResponse {
	resp := eventResponse{
		Type:      string(event.Type),
		Streak:    event.Streak,
		CreatedAt: event.CreatedAt,
	}

	switch event.Type {
	case entity.StreakReward:
		resp.XP = &event.XP

		if event.Coins > 0 {
			resp.Coins = &event.Coins
		}

	case entity.NewReward:
		if event.Reward != nil {
			resp.Reward = &rewardResponse{
				RewardID:     event.Reward.ID,
				PromoCode:    event.Reward.PromoCode,
				Name:         event.Reward.Definition.Name,
				Description:  event.Reward.Definition.Description,
				Status:       string(event.Reward.Status),
				EarnedReason: event.Reward.EarnedReason,
				ExpiresAt:    event.Reward.ExpiresAt,
				RedeemedAt:   event.Reward.RedeemedAt,
			}
		}
	}

	return resp
}
