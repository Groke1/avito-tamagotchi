package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToUUID(userID string) (pgtype.UUID, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("get uuid: %w", err)
	}
	return pgtype.UUID{
		Bytes: parsedUserID,
		Valid: true,
	}, nil
}

func TimestamptzToTime(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
func TimeToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{
			Time:  time.Time{},
			Valid: false,
		}
	}
	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}

func PgInt4ToInt32(pgInt pgtype.Int4) *int32 {
	if !pgInt.Valid {
		return nil
	}
	return &pgInt.Int32
}

func Int32ToPgInt4(val *int32) pgtype.Int4 {
	if val != nil {
		return pgtype.Int4{
			Int32: *val,
			Valid: true,
		}
	}
	return pgtype.Int4{}
}

func PgUUIDtoString(pgUUID pgtype.UUID) *string {
	if !pgUUID.Valid {
		return nil
	}
	return new(pgUUID.String())
}
