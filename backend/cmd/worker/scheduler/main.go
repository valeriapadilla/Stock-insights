package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/app"
	"github.com/valeriapadilla/stock-insights/internal/client"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/repository"
	"github.com/valeriapadilla/stock-insights/internal/service"
	"github.com/valeriapadilla/stock-insights/internal/worker"
	"github.com/valeriapadilla/stock-insights/internal/worker/implementations"
)

func main() {
	cfg := config.Load()
	app.SetupLogging(cfg)
	logger := logrus.StandardLogger()

	logger.Info("Starting Scheduled Data Worker...")

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	externalClient := client.NewExternalAPIClient(client.ExternalAPIConfig{
		BaseURL: cfg.ExternalAPIURL,
		APIKey:  cfg.ExternalAPIKey,
	}, logger)

	stockRepo := repository.NewStockRepository(database.DB)
	stockCmd := repository.NewStockCommand(database.DB)

	dataWorker := implementations.NewDataWorker(
		externalClient,
		stockRepo,
		stockCmd,
		logger,
		implementations.DataWorkerConfig{
			ScheduleInterval: 24 * time.Hour,
			MaxRetries:       3,
			RetryDelay:       5 * time.Second,
		},
	)

	ingestionService := service.NewIngestionService(dataWorker, logger)

	jobManager := worker.NewJobManager(3, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutting down Scheduled Data Worker...")
		cancel()
	}()

	if err := runScheduler(ctx, ingestionService, jobManager, logger); err != nil {
		log.Fatal("Scheduler failed:", err)
	}
}

func runScheduler(ctx context.Context, ingestionService *service.IngestionService, jobManager *worker.JobManager, logger *logrus.Logger) error {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	logger.Info("Starting scheduler with 24-hour interval")

	if err := runScheduledIngestion(ctx, ingestionService, jobManager, logger); err != nil {
		logger.WithError(err).Error("Initial scheduled ingestion failed")
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Scheduler stopped")
			return nil
		case <-ticker.C:
			if err := runScheduledIngestion(ctx, ingestionService, jobManager, logger); err != nil {
				logger.WithError(err).Error("Scheduled ingestion failed")
			}
		}
	}
}

func runScheduledIngestion(ctx context.Context, ingestionService *service.IngestionService, jobManager *worker.JobManager, logger *logrus.Logger) error {
	logger.WithContext(ctx).Info("Starting scheduled ingestion...")

	job, err := jobManager.CreateJob()
	if err != nil {
		logger.WithContext(ctx).WithError(err).Error("Failed to create scheduled job")
		return err
	}

	logger.WithContext(ctx).WithField("job_id", job.ID).Info("Created scheduled ingestion job")

	if err := jobManager.RunJobAsync(job.ID, func(ctx context.Context) error {
		return ingestionService.TriggerIngestionAsync(ctx)
	}); err != nil {
		logger.WithContext(ctx).WithError(err).Error("Failed to start scheduled ingestion job")
		return err
	}

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"job_id": job.ID,
		"type":   "scheduled",
	}).Info("Scheduled ingestion job started successfully")

	return nil
}
