package worker

import (
	"context"
	"log"
)

type GoroutineWorkerInterface interface {
	Execute(task func(ctx context.Context) error) error
}

type goroutineWorker struct{}

func NewGoroutineWorker() GoroutineWorkerInterface {
	return &goroutineWorker{}
}

func (w *goroutineWorker) Execute(task func(ctx context.Context) error) error {
	go func() {
		ctx := context.Background()
		if err := task(ctx); err != nil {
			log.Printf("background task failed: %v", err)
		}
	}()
	return nil
}
