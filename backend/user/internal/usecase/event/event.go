package event

import (
	"context"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type (
	eventRepository interface {
		GetUserEventsAndMarkDelivered(ctx context.Context, userID string) ([]entity.UserEvent, error)
	}

	rewardRepository interface {
		GetRewardByUserIDAndRewardID(ctx context.Context, userID, rewardID string) (*entity.UserReward, error)
	}
)

type eventService struct {
	eventRepository  eventRepository
	rewardRepository rewardRepository
}

func NewEventService(
	eventRepository eventRepository,
	rewardRepository rewardRepository,
) *eventService {
	return &eventService{
		rewardRepository: rewardRepository,
		eventRepository:  eventRepository,
	}
}

func (s *eventService) GetNewEvents(ctx context.Context, userID string) ([]entity.UserEventDetails, error) {
	events, err := s.eventRepository.GetUserEventsAndMarkDelivered(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get new rewards: get user events: %w", err)
	}

	result := make([]entity.UserEventDetails, 0, len(events))

	for _, event := range events {
		details := entity.UserEventDetails{
			Type:      event.Type,
			XP:        event.XP,
			Coins:     event.Coins,
			Streak:    event.Streak,
			CreatedAt: event.CreatedAt,
		}

		if event.UserRewardID != nil {
			reward, err := s.rewardRepository.GetRewardByUserIDAndRewardID(
				ctx, userID, *event.UserRewardID,
			)
			if err != nil {
				return nil, fmt.Errorf("get event reward: %w", err)
			}

			details.Reward = reward
		}

		result = append(result, details)
	}

	return result, nil
}
