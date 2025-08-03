package implementations

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/service/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
	workerInterfaces "github.com/valeriapadilla/stock-insights/internal/worker/interfaces"
)

type RecommendationWorkerImpl struct {
	recommendationService interfaces.RecommendationServiceInterface
	stockRepo             repoInterfaces.StockRepository
	logger                *logrus.Logger
	isRunning             bool
	mutex                 sync.RWMutex
	lastRunTime           *time.Time
}

func NewRecommendationWorker(
	recommendationService interfaces.RecommendationServiceInterface,
	stockRepo repoInterfaces.StockRepository,
	logger *logrus.Logger,
) workerInterfaces.RecommendationWorker {
	return &RecommendationWorkerImpl{
		recommendationService: recommendationService,
		stockRepo:             stockRepo,
		logger:                logger,
		isRunning:             false,
	}
}

func (w *RecommendationWorkerImpl) RunDailyRecommendations(ctx context.Context) error {
	w.setRunning(true)
	defer w.setRunning(false)

	w.logger.Info("Starting daily recommendation process...")

	if err := w.checkDataFreshness(ctx); err != nil {
		w.logger.WithError(err).Error("Data freshness check failed")
		return err
	}

	if err := w.calculateRecommendations(ctx); err != nil {
		w.logger.WithError(err).Error("Failed to calculate recommendations")
		return err
	}

	w.setLastRunTime(time.Now())
	w.logger.Info("Daily recommendation process completed successfully")
	return nil
}

func (w *RecommendationWorkerImpl) GetLastRunTime() (*time.Time, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.lastRunTime, nil
}

func (w *RecommendationWorkerImpl) IsRunning() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.isRunning
}

func (w *RecommendationWorkerImpl) checkDataFreshness(ctx context.Context) error {
	w.logger.WithContext(ctx).Info("Checking data freshness for recommendations...")

	lastUpdate, err := w.stockRepo.GetLastUpdateTime()
	if err != nil {
		w.logger.WithError(err).Error("Failed to get last update time")
		return fmt.Errorf("failed to get last update time: %w", err)
	}

	if lastUpdate == nil {
		w.logger.Error("No stock data found in database")
		return fmt.Errorf("no stock data found in database")
	}

	maxAge := 24 * time.Hour
	timeSinceUpdate := time.Since(*lastUpdate)

	w.logger.WithFields(logrus.Fields{
		"last_update":       lastUpdate.Format(time.RFC3339),
		"time_since_update": timeSinceUpdate.String(),
		"max_age":           maxAge.String(),
	}).Info("Data freshness check")

	if timeSinceUpdate > maxAge {
		w.logger.WithFields(logrus.Fields{
			"last_update":       lastUpdate.Format(time.RFC3339),
			"time_since_update": timeSinceUpdate.String(),
			"max_age":           maxAge.String(),
		}).Warn("Data is stale, but proceeding with recommendations")

	}

	w.logger.Info("Data freshness check passed, proceeding with recommendations...")
	return nil
}

func (w *RecommendationWorkerImpl) calculateRecommendations(ctx context.Context) error {
	w.logger.WithContext(ctx).Info("Calculating recommendations with available data...")

	params := validator.RecommendationParams{
		DaysBack:   7,
		MaxResults: 30,
		MinScore:   80,
	}

	recommendations, err := w.recommendationService.CalculateRecommendations(params)
	if err != nil {
		return err
	}

	if err := w.recommendationService.SaveRecommendations(recommendations); err != nil {
		w.logger.WithError(err).Error("Failed to save recommendations")
		return err
	}

	w.logger.Info("Recommendations calculated and saved successfully")
	return nil
}

func (w *RecommendationWorkerImpl) setRunning(running bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.isRunning = running
}

func (w *RecommendationWorkerImpl) setLastRunTime(t time.Time) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.lastRunTime = &t
}
