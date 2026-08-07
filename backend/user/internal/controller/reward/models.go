package reward

import (
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type userRewardResponse struct {
	RewardID    string     `json:"reward_id"`
	PromoCode   string     `json:"promo_code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at"`
	RedeemedAt  *time.Time `json:"redeemed_at"`
}

type userListRewardResponse struct {
	Items []userRewardResponse `json:"items"`
}

type redeemRewardRequest struct {
	PromoCode string `json:"promo_code"`
}

type getDefinitionResponse struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type grantRewardRequest struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
}

func toUserRewardResponse(userReward entity.UserReward) userRewardResponse {
	return userRewardResponse{
		RewardID:    userReward.ID,
		PromoCode:   userReward.PromoCode,
		Name:        userReward.Definition.Name,
		Description: userReward.Definition.Description,
		Status:      string(userReward.Status),
		ExpiresAt:   userReward.ExpiresAt,
		RedeemedAt:  userReward.RedeemedAt,
	}
}
