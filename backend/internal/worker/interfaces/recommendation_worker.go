package interfaces

import (
	"context"
	"time"
)

type RecommendationWorker interface {
	RunDailyRecommendations(ctx context.Context) error
	GetLastRunTime() (*time.Time, error)
	IsRunning() bool
}
