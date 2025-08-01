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
	externalClient client.ExternalAPIClient
	stockRepo      repoInterfaces.StockRepository
	stockCommand   repoInterfaces.StockCommand
	logger         *logrus.Logger
	config         DataWorkerConfig
}

func NewDataWorker(
	externalClient client.ExternalAPIClient,
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
	w.logger.Info("Starting stock data fetch and processing")

	stocks, err := w.externalClient.GetAllStocks(ctx)
	if err != nil {
		w.logger.WithError(err).Error("Failed to fetch stocks from external API")
		return errors.NewExternalError("failed to fetch stocks from external API", err)
	}

	w.logger.WithField("stocks_count", len(stocks)).Info("Successfully fetched stocks from external API")

	if err := w.saveStocks(ctx, stocks); err != nil {
		w.logger.WithError(err).Error("Failed to save stocks to database")
		return errors.NewDatabaseError("failed to save stocks to database", err)
	}

	w.logger.WithFields(logrus.Fields{
		"stocks_processed": len(stocks),
		"stocks_saved":     len(stocks),
	}).Info("Successfully processed and saved stocks")

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

func (w *DataWorkerImpl) saveStocks(ctx context.Context, stocks []model.Stock) error {
	if len(stocks) == 0 {
		w.logger.WithContext(ctx).Warn("No stocks to save")
		return nil
	}

	stockPtrs := make([]*model.Stock, len(stocks))
	for i := range stocks {
		stockPtrs[i] = &stocks[i]
	}

	if err := w.stockCommand.BulkUpsert(stockPtrs); err != nil {
		w.logger.WithError(err).Error("Failed to save stocks to database")
		return err
	}

	w.logger.WithField("stocks_count", len(stocks)).Info("Successfully saved stocks to database")
	return nil
}
