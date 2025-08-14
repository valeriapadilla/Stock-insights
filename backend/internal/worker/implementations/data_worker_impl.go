package implementations

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/client"
	"github.com/valeriapadilla/stock-insights/internal/errors"
	"github.com/valeriapadilla/stock-insights/internal/model"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	workerInterfaces "github.com/valeriapadilla/stock-insights/internal/worker/interfaces"
)

type DataWorkerConfig struct {
	ScheduleInterval time.Duration
	MaxRetries       int
	RetryDelay       time.Duration
}

type DataWorkerImpl struct {
	externalClient *client.ExternalAPIClient
	stockRepo      repoInterfaces.StockRepository
	stockCommand   repoInterfaces.StockCommand
	logger         *logrus.Logger
	config         DataWorkerConfig
}

func NewDataWorker(
	externalClient *client.ExternalAPIClient,
	stockRepo repoInterfaces.StockRepository,
	stockCommand repoInterfaces.StockCommand,
	logger *logrus.Logger,
	config DataWorkerConfig,
) workerInterfaces.DataWorker {
	return &DataWorkerImpl{
		externalClient: externalClient,
		stockRepo:      stockRepo,
		stockCommand:   stockCommand,
		logger:         logger,
		config:         config,
	}
}

func (w *DataWorkerImpl) FetchAndProcessStocks(ctx context.Context) error {
	return w.FetchAndProcessStocksEfficient(ctx)
}

func (w *DataWorkerImpl) FetchAndProcessStocksEfficient(ctx context.Context) error {
	w.logger.Info("Starting stock data fetch and processing (UPSERT strategy)")

	allStocks, err := w.externalClient.GetAllStocks(ctx)
	if err != nil {
		w.logger.WithError(err).Error("Failed to fetch stocks from external API")
		return errors.NewExternalError("failed to fetch stocks from external API", err)
	}

	w.logger.WithFields(logrus.Fields{
		"total_stocks": len(allStocks),
	}).Info("Fetched all stocks from external API")

	if len(allStocks) == 0 {
		w.logger.Info("No stocks found from external API")
		return nil
	}

	if err := w.saveStocksInBatchesOptimized(ctx, allStocks); err != nil {
		w.logger.WithError(err).Error("Failed to save stocks to database")
		return errors.NewDatabaseError("failed to save stocks to database", err)
	}

	w.logger.WithFields(logrus.Fields{
		"stocks_processed": len(allStocks),
		"stocks_saved":     len(allStocks),
	}).Info("Successfully processed and saved all stocks using UPSERT")

	return nil
}

func (w *DataWorkerImpl) processAllStocks(ctx context.Context) error {
	w.logger.Info("Redirecting to efficient processing method")
	return w.FetchAndProcessStocksEfficient(ctx)
}

func (w *DataWorkerImpl) filterStocksByDate(allStocks []model.Stock, since time.Time) []model.Stock {
	var filteredStocks []model.Stock

	for _, stock := range allStocks {
		if stock.Time.After(since) {
			filteredStocks = append(filteredStocks, stock)
		}
	}

	w.logger.WithFields(logrus.Fields{
		"total_stocks":    len(allStocks),
		"filtered_stocks": len(filteredStocks),
		"since":           since.Format(time.RFC3339),
	}).Info("Applied local date filtering (DEPRECATED)")

	return filteredStocks
}

func (w *DataWorkerImpl) saveStocksInBatchesOptimized(ctx context.Context, stocks []model.Stock) error {
	if len(stocks) == 0 {
		w.logger.WithContext(ctx).Warn("No stocks to save")
		return nil
	}

	const batchSize = 500
	totalStocks := len(stocks)
	now := time.Now()

	w.logger.WithFields(logrus.Fields{
		"total_stocks": totalStocks,
		"batch_size":   batchSize,
	}).Info("Starting batch processing of stocks")

	for i := 0; i < totalStocks; i += batchSize {
		end := i + batchSize
		if end > totalStocks {
			end = totalStocks
		}

		batch := stocks[i:end]
		stockPtrs := make([]*model.Stock, len(batch))
		for j := range batch {
			stockPtrs[j] = &batch[j]
			if stockPtrs[j].CreatedAt.IsZero() {
				stockPtrs[j].CreatedAt = now
			}
			stockPtrs[j].UpdatedAt = now
		}

		if err := w.stockCommand.BulkUpsert(stockPtrs); err != nil {
			w.logger.WithError(err).WithFields(logrus.Fields{
				"batch_start": i,
				"batch_end":   end,
			}).Error("Failed to save batch of stocks")
			return err
		}

		progress := (end * 100) / totalStocks
		w.logger.WithFields(logrus.Fields{
			"batch_start": i,
			"batch_end":   end,
			"progress":    progress,
		}).Info("Saved batch of stocks")
	}

	w.logger.WithField("stocks_count", totalStocks).Info("Successfully saved all stocks to database")
	return nil
}

func (w *DataWorkerImpl) HealthCheck(ctx context.Context) error {
	w.logger.Debug("Performing external API health check")

	if err := w.externalClient.HealthCheck(ctx); err != nil {
		w.logger.WithError(err).Error("External API health check failed")
		return errors.NewExternalError("external API health check failed", err)
	}

	w.logger.Info("External API health check passed")
	return nil
}

func (w *DataWorkerImpl) GetLastRunTime(ctx context.Context) (*time.Time, error) {
	return w.stockRepo.GetLastUpdateTime()
}

func (w *DataWorkerImpl) ShouldRun(ctx context.Context) (bool, error) {
	lastRun, err := w.GetLastRunTime(ctx)
	if err != nil {
		return false, errors.NewInternalError("failed to get last run time", err)
	}

	if lastRun == nil {
		w.logger.Info("No previous run found, should run")
		return true, nil
	}

	timeSinceLastRun := time.Since(*lastRun)
	shouldRun := timeSinceLastRun >= w.config.ScheduleInterval

	w.logger.WithFields(logrus.Fields{
		"last_run":            lastRun,
		"time_since_last_run": timeSinceLastRun,
		"schedule_interval":   w.config.ScheduleInterval,
		"should_run":          shouldRun,
	}).Debug("Checking if worker should run")

	return shouldRun, nil
}
