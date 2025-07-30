package server

import (
	"github.com/gin-gonic/gin"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/handler"
	"github.com/valeriapadilla/stock-insights/internal/middleware"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	server := &Server{
		config: cfg,
		router: gin.Default(),
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.RateLimitMiddleware(s.config.RateLimit, 10))
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", handler.HealthCheck)
	s.router.HEAD("/health", handler.HealthCheckHead)
	
	publicAPI := s.router.Group("/api/public")
	{
		publicAPI.GET("/health", handler.HealthCheck)
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
} 