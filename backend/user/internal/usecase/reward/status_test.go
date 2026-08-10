package reward

import (
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func TestUpdateStatus(t *testing.T) {
	now := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name   string
		reward entity.UserReward
		want   entity.Status
	}{
		{
			name:   "already redeemed stays redeemed even if expired",
			reward: entity.UserReward{Status: entity.StatusRedeemed, ExpiresAt: &past},
			want:   entity.StatusRedeemed,
		},
		{
			name:   "expires_at in the past becomes expired",
			reward: entity.UserReward{Status: entity.StatusActive, ExpiresAt: &past},
			want:   entity.StatusExpired,
		},
		{
			name:   "expires_at exactly now becomes expired",
			reward: entity.UserReward{Status: entity.StatusActive, ExpiresAt: &now},
			want:   entity.StatusExpired,
		},
		{
			name:   "expires_at in the future stays active",
			reward: entity.UserReward{Status: entity.StatusActive, ExpiresAt: &future},
			want:   entity.StatusActive,
		},
		{
			name:   "no expiration (nil) stays active",
			reward: entity.UserReward{Status: entity.StatusActive, ExpiresAt: nil},
			want:   entity.StatusActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateStatus(&tt.reward, now)
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
