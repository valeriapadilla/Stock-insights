package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/config"
)

// SetupLogging configura el sistema de logging
func SetupLogging(cfg *config.Config) {
	// Configurar nivel de log
	setupLogLevel(cfg)
	
	// Configurar modo de producción
	setupProductionMode(cfg)
}

// setupLogLevel configura el nivel de logging
func setupLogLevel(cfg *config.Config) {
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Warn("Invalid log level, using INFO as default")
	} else {
		logrus.SetLevel(level)
	}
}

// setupProductionMode configura el modo de producción
func setupProductionMode(cfg *config.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		logrus.Info("Running in production mode")
	} else {
		logrus.Info("Running in development mode")
	}
} 