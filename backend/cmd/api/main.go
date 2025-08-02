package main

import (
	"log"

	"github.com/valeriapadilla/stock-insights/internal/bootstrap"
	"github.com/valeriapadilla/stock-insights/internal/config"
)

func main() {
	cfg := config.Load()
	application := bootstrap.Bootstrap(cfg)

	if err := application.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
