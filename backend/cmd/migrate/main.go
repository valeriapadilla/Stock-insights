package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/database"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Info("Starting database migrations...")

	err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	mm := database.NewMigrationManager(database.DB)
	if err := mm.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	logrus.Info("Migrations completed successfully")
}
