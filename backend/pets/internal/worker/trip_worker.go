package worker

import (
	"context"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"
)

const TripWorkerCooldown = 20 * time.Second

type TripWorker struct {
	tripRepository *repository.TripRepository
	tripService    *service.TripService
}

func NewTripWorker(tripRepository *repository.TripRepository, tripService *service.TripService) *TripWorker {
	return &TripWorker{tripRepository: tripRepository, tripService: tripService}
}

func (tw *TripWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(TripWorkerCooldown)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			tw.processExpiredTrips(ctx)
		}
	}
}

func (tw *TripWorker) processExpiredTrips(ctx context.Context) {
	trips, err := tw.tripRepository.GetFinishedTrips(ctx)
	if err != nil {
		return
	}

	for _, trip := range trips {
		tw.tripService.SendToClient(&trip)
	}
}
