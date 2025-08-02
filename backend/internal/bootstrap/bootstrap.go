package bootstrap

import (
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/app"
	"github.com/valeriapadilla/stock-insights/internal/client"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/repository"
	"github.com/valeriapadilla/stock-insights/internal/server"
	"github.com/valeriapadilla/stock-insights/internal/service"
	"github.com/valeriapadilla/stock-insights/internal/worker/implementations"
)

type App struct {
	config *config.Config
	server *server.Server
	logger *logrus.Logger
}

func Bootstrap(cfg *config.Config) *App {
	app.SetupLogging(cfg)
	logger := logrus.StandardLogger()

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
	srv := server.NewServer(cfg, ingestionService, logger)

	return &App{
		config: cfg,
		server: srv,
		logger: logger,
	}
}

func (a *App) Run() error {
	a.logger.Infof("Starting server on port %s", a.config.Port)
	return a.server.Run()
}
