package main

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/app"
	"github.com/valeriapadilla/stock-insights/internal/client"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/repository"
	"github.com/valeriapadilla/stock-insights/internal/worker/implementations"
	workerInterfaces "github.com/valeriapadilla/stock-insights/internal/worker/interfaces"
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
		implementations.DataWorkerConfig{},
	)

	ctx := context.Background()

	if err := runIngestion(ctx, dataWorker, logger); err != nil {
		log.Fatal("Ingestion failed:", err)
	}

	logger.Info("Ingestion completed successfully")
}

func runIngestion(ctx context.Context, dataWorker workerInterfaces.DataWorker, logger *logrus.Logger) error {
	logger.Info("Starting scheduled ingestion...")

	if err := dataWorker.FetchAndProcessStocks(ctx); err != nil {
		logger.WithError(err).Error("Failed to complete ingestion")
		return err
	}

	logger.Info("Scheduled ingestion completed successfully")
	return nil
}
