package interfaces

import (
	"context"
	"time"
)

type DataWorker interface {
	FetchAndProcessStocks(ctx context.Context) error
	FetchAndProcessStocksSmart(ctx context.Context) error
	FetchAndProcessStocksIncremental(ctx context.Context) error
	HealthCheck(ctx context.Context) error
	GetLastRunTime(ctx context.Context) (*time.Time, error)
	ShouldRun(ctx context.Context) (bool, error)
}
