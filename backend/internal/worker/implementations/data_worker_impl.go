package implementations

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	clientInterfaces "github.com/valeriapadilla/stock-insights/internal/client/interfaces"
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
	externalClient clientInterfaces.ExternalAPIClient
	stockRepo      repoInterfaces.StockRepository
	stockCommand   repoInterfaces.StockCommand
	logger         *logrus.Logger
	config         DataWorkerConfig
}

func NewDataWorker(
	externalClient clientInterfaces.ExternalAPIClient,
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
	w.logger.Info("Starting incremental stock data fetch and processing")

	lastUpdate, err := w.stockRepo.GetLastUpdateTime()
	if err != nil {
		w.logger.WithError(err).Warn("Failed to get last update time, falling back to full fetch")
		return w.FetchAndProcessStocksIncremental(ctx)
	}

	if lastUpdate == nil {
		w.logger.Info("No previous update found, fetching all stocks")
		return w.FetchAndProcessStocksIncremental(ctx)
	}

	allStocks, err := w.externalClient.GetAllStocks(ctx)
	if err != nil {
		w.logger.WithError(err).Error("Failed to fetch stocks from external API")
		return errors.NewExternalError("failed to fetch stocks from external API", err)
	}

	w.logger.WithFields(logrus.Fields{
		"total_stocks": len(allStocks),
		"since":        lastUpdate.Format(time.RFC3339),
	}).Info("Fetched stocks from external API")

	newStocks := w.filterStocksByDateOptimized(allStocks, *lastUpdate)

	w.logger.WithFields(logrus.Fields{
		"total_stocks": len(allStocks),
		"new_stocks":   len(newStocks),
		"efficiency":   float64(len(newStocks)) / float64(len(allStocks)) * 100,
	}).Info("Filtered stocks by date")

	if len(newStocks) == 0 {
		w.logger.Info("No new stocks found")
		return nil
	}

	if err := w.saveStocksInBatchesOptimized(ctx, newStocks); err != nil {
		w.logger.WithError(err).Error("Failed to save stocks to database")
		return errors.NewDatabaseError("failed to save stocks to database", err)
	}

	w.logger.WithFields(logrus.Fields{
		"stocks_processed": len(newStocks),
		"stocks_saved":     len(newStocks),
	}).Info("Successfully processed and saved stocks")

	return nil
}

func (w *DataWorkerImpl) FetchAndProcessStocksIncremental(ctx context.Context) error {
	w.logger.Info("Starting incremental stock data fetch and processing")

	lastUpdate, err := w.stockRepo.GetLastUpdateTime()
	if err != nil {
		w.logger.WithError(err).Warn("Failed to get last update time, will fetch all stocks")
		lastUpdate = nil
	}

	var stocks []model.Stock
	var isIncremental bool

	if lastUpdate != nil {
		w.logger.WithField("since", lastUpdate.Format(time.RFC3339)).Info("Fetching all stocks and applying local date filtering")

		allStocks, err := w.externalClient.GetAllStocks(ctx)
		if err != nil {
			w.logger.WithError(err).Error("Failed to fetch stocks from external API")
			return errors.NewExternalError("failed to fetch stocks from external API", err)
		}

		stocks = w.filterStocksByDate(allStocks, *lastUpdate)
		isIncremental = true
	} else {
		w.logger.Info("No previous update found, fetching all stocks")
		stocks, err = w.externalClient.GetAllStocks(ctx)
		if err != nil {
			w.logger.WithError(err).Error("Failed to fetch stocks from external API")
			return errors.NewExternalError("failed to fetch stocks from external API", err)
		}
		isIncremental = false
	}

	w.logger.WithFields(logrus.Fields{
		"stocks_count": len(stocks),
		"incremental":  isIncremental,
		"since":        lastUpdate,
	}).Info("Successfully fetched stocks")

	if len(stocks) == 0 {
		w.logger.Info("No new stocks to process")
		return nil
	}

	if err := w.saveStocksInBatches(ctx, stocks); err != nil {
		w.logger.WithError(err).Error("Failed to save stocks to database")
		return errors.NewDatabaseError("failed to save stocks to database", err)
	}

	w.logger.WithFields(logrus.Fields{
		"stocks_processed": len(stocks),
		"stocks_saved":     len(stocks),
		"incremental":      isIncremental,
	}).Info("Successfully processed and saved stocks")

	return nil
}

func (w *DataWorkerImpl) FetchAndProcessStocksSmart(ctx context.Context) error {
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
	}).Info("Applied local date filtering")

	return filteredStocks
}

func (w *DataWorkerImpl) filterStocksByDateOptimized(allStocks []model.Stock, since time.Time) []model.Stock {
	var filteredStocks []model.Stock
	filteredStocks = make([]model.Stock, 0, len(allStocks)/3)

	for _, stock := range allStocks {
		if stock.Time.After(since) {
			filteredStocks = append(filteredStocks, stock)
		}
	}

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

func (w *DataWorkerImpl) saveStocksInBatches(ctx context.Context, stocks []model.Stock) error {
	if len(stocks) == 0 {
		w.logger.WithContext(ctx).Warn("No stocks to save")
		return nil
	}

	const batchSize = 100
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
		}).Info("Successfully saved batch of stocks")
	}

	w.logger.WithField("stocks_count", totalStocks).Info("Successfully saved all stocks to database")
	return nil
}
