package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/handler"
	v1 "github.com/valeriapadilla/stock-insights/internal/handler/v1"
	"github.com/valeriapadilla/stock-insights/internal/middleware"
	"github.com/valeriapadilla/stock-insights/internal/service"
	"github.com/valeriapadilla/stock-insights/internal/worker"
)

type Server struct {
	router           *gin.Engine
	config           *config.Config
	ingestionService *service.IngestionService
	jobManager       *worker.JobManager
	logger           *logrus.Logger
}

func NewServer(cfg *config.Config, ingestionService *service.IngestionService, logger *logrus.Logger) *Server {
	server := &Server{
		config:           cfg,
		router:           gin.New(),
		ingestionService: ingestionService,
		jobManager:       worker.NewJobManager(5, logger),
		logger:           logger,
	}

	server.setupMiddleware()
	server.setupRoutes()
	server.startCleanupRoutine()

	return server
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestIDMiddleware())
	s.router.Use(middleware.LoggingMiddleware())
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.RateLimitMiddleware(s.config.RateLimit, 10))
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", handler.HealthCheck)
	s.router.HEAD("/health", handler.HealthCheckHead)

	v1API := s.router.Group("/api/v1")
	{
		publicV1 := v1API.Group("/public")
		publicV1.GET("/health", handler.HealthCheck)
		// publicV1.GET("/stocks", stockHandler.ListStocks)
		// publicV1.GET("/stocks/:id", stockHandler.GetStock)
		// publicV1.GET("/stocks/search", stockHandler.SearchStocks)

		adminV1 := v1API.Group("/admin")
		adminV1.Use(middleware.AuthMiddleware())
		{
			stocksIngestionHandler := v1.NewStocksIngestionHandler(s.ingestionService, s.jobManager, s.logger)
			adminV1.POST("/ingest/stocks", stocksIngestionHandler.TriggerIngestion)
			adminV1.GET("/jobs/:jobId", stocksIngestionHandler.GetJobStatus)
		}
	}
}

func (s *Server) startCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			s.jobManager.CleanupOldJobs(24 * time.Hour)
		}
	}()
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}
