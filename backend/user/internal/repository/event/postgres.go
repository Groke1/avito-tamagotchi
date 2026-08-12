package event

import (
	"context"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/repository/converter"
	sqlcevent "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/event/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type eventRepository struct {
	db      sqlcevent.DBTX
	queries *sqlcevent.Queries
}

func NewEventRepository(qdb sqlcevent.DBTX) *eventRepository {
	return &eventRepository{
		db:      qdb,
		queries: sqlcevent.New(qdb),
	}
}

func (r *eventRepository) GetUserEventsAndMarkDelivered(ctx context.Context, userID string) ([]entity.UserEvent, error) {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("user uuid conversion error: %w", err)
	}
	userEventsRaw, err := r.getQueries(ctx).GetUserEventsAndMarkDelivered(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user events error: %w", err)
	}
	userEvents := make([]entity.UserEvent, 0, len(userEventsRaw))
	for _, userEvent := range userEventsRaw {
		userEvents = append(userEvents, entity.UserEvent{
			ID:           userEvent.ID,
			UserID:       userEvent.UserID.String(),
			Type:         entity.EventType(userEvent.Type),
			XP:           userEvent.Xp,
			Coins:        userEvent.Coins,
			Streak:       converter.PgInt4ToInt32(userEvent.Streak),
			UserRewardID: converter.PgUUIDtoString(userEvent.UserRewardID),
			CreatedAt:    *converter.TimestamptzToTime(userEvent.CreatedAt),
		})
	}
	return userEvents, nil
}

func (r *eventRepository) AddUserEvent(ctx context.Context, event entity.UserEvent) error {
	userUUID, err := converter.StringToUUID(event.UserID)
	if err != nil {
		return fmt.Errorf("user uuid conversion error: %w", err)
	}
	var rewardUserUUID pgtype.UUID
	if event.UserRewardID != nil {
		rewardUserUUID, err = converter.StringToUUID(*event.UserRewardID)
		if err != nil {
			return fmt.Errorf("reward uuid conversion error: %w", err)
		}
	}

	err = r.getQueries(ctx).AddUserEvent(ctx, sqlcevent.AddUserEventParams{
		UserID:       userUUID,
		Type:         sqlcevent.UsersEventType(event.Type),
		Xp:           event.XP,
		Coins:        event.Coins,
		Streak:       converter.Int32ToPgInt4(event.Streak),
		UserRewardID: rewardUserUUID,
	})
	if err != nil {
		return fmt.Errorf("add user event error: %w", err)
	}

	return nil
}

func (r *eventRepository) DeleteDeliveredEvents(ctx context.Context) error {
	if err := r.getQueries(ctx).DeleteDeliveredEvents(ctx); err != nil {
		return fmt.Errorf("delete delivered events: %w", err)
	}
	return nil
}

func (r *eventRepository) getQueries(ctx context.Context) *sqlcevent.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return sqlcevent.New(tx)
}
