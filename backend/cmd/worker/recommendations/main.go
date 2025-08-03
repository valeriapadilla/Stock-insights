package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/app"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/repository"
	"github.com/valeriapadilla/stock-insights/internal/service"
	"github.com/valeriapadilla/stock-insights/internal/worker/implementations"
)

func main() {
	cfg := config.Load()

	app.SetupLogging(cfg)
	logger := logrus.StandardLogger()

	logger.Info("Starting Recommendation Worker...")

	if err := database.Connect(); err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close()

	stockRepo := repository.NewStockRepository(database.DB)
	recommendationRepo := repository.NewRecommendationRepository(database.DB)
	recommendationCmd := repository.NewRecommendationCommand(database.DB, stockRepo)

	recommendationService := service.NewRecommendationService(stockRepo, recommendationRepo, recommendationCmd, logger)

	recommendationWorker := implementations.NewRecommendationWorker(recommendationService, stockRepo, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Running daily recommendation process...")

		if err := recommendationWorker.RunDailyRecommendations(ctx); err != nil {
			logger.WithError(err).Error("Failed to run daily recommendation process")
			os.Exit(1)
		}

		logger.Info("Daily recommendation process completed successfully")
		os.Exit(0)
	}()

	select {
	case sig := <-sigChan:
		logger.WithField("signal", sig).Info("Received shutdown signal")
		cancel()
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	logger.Info("Shutting down Recommendation Worker...")
	time.Sleep(5 * time.Second)
}
