package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/config"
)

func main() {
	cfg := config.Load()

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"message": "Stock Insights API is running",
		})
	})

	logrus.Infof("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
