package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/config"
)

func SetupLogging(cfg *config.Config) {
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	if cfg.Environment == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	logrus.SetOutput(os.Stdout)
}

func GetLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"service": "stock-insights",
		"version": "1.0.0",
	})
}
