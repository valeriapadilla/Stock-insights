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
	"github.com/valeriapadilla/stock-insights/internal/job"
	"github.com/valeriapadilla/stock-insights/internal/repository"
	"github.com/valeriapadilla/stock-insights/internal/service"
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

	jobManager := job.NewJobManager(3, logger)

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

func runScheduler(ctx context.Context, ingestionService *service.IngestionService, jobManager *job.JobManager, logger *logrus.Logger) error {
	logger.Info("Starting single ingestion execution")

	if err := runScheduledIngestion(ctx, ingestionService, jobManager, logger); err != nil {
		logger.WithError(err).Error("Ingestion failed")
		return err
	}

	logger.Info("Ingestion completed successfully")
	return nil
}

func runScheduledIngestion(ctx context.Context, ingestionService *service.IngestionService, jobManager *job.JobManager, logger *logrus.Logger) error {
	logger.WithContext(ctx).Info("Starting scheduled ingestion...")

	job, err := jobManager.CreateJob()
	if err != nil {
		logger.WithError(err).Error("Failed to create job")
		return err
	}

	if err := jobManager.RunJobAsync(job.ID, func(ctx context.Context) error {
		return ingestionService.TriggerIngestionAsync(ctx)
	}); err != nil {
		logger.WithError(err).Error("Failed to start ingestion job")
		return err
	}

	logger.WithField("job_id", job.ID).Info("Scheduled ingestion job started")
	return nil
}
