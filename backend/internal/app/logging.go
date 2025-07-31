package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/config"
)

func SetupLogging(cfg *config.Config) {
	setupLogLevel(cfg)
	setupProductionMode(cfg)
}

func setupLogLevel(cfg *config.Config) {
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Warn("Invalid log level, using INFO as default")
	} else {
		logrus.SetLevel(level)
	}
}

func setupProductionMode(cfg *config.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		logrus.Info("Running in production mode")
	} else {
		logrus.Info("Running in development mode")
	}
}
