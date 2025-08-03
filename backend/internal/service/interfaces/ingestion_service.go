package interfaces

import (
	"context"
)

type IngestionServiceInterface interface {
	TriggerIngestionAsync(ctx context.Context) error
}
