package domain

type TripEventTone string

const (
	NegativeTone = "negative"
	PositiveTone = "positive"
)

type TripEvent struct {
	ID          int32
	Description string
	IsNegative  bool
}
