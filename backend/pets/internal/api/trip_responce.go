package api

type TripStatus string

const (
	TripStatusStarted   TripStatus = "Trip started"
	TripStatusCompleted TripStatus = "Trip completed"
)

type TripResponse struct {
	Status TripStatus `json:"status"`
}
