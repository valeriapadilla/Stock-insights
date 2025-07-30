package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/app"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/server"
)

func main() {
	cfg := config.Load()
	app.SetupLogging(cfg)
	srv := server.NewServer(cfg)

	logrus.Infof("Starting server on port %s", cfg.Port)
	if err := srv.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
