package interfaces

import (
	"context"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type ExternalAPIClient interface {
	GetAllStocks(ctx context.Context) ([]model.Stock, error)
	HealthCheck(ctx context.Context) error
}
