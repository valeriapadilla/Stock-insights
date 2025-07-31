package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/valeriapadilla/stock-insights/internal/config"
)

var DB *sql.DB

func Connect() error {
	var databaseURL string

	if os.Getenv("TESTING") == "true" {
		testCfg := config.LoadTestConfig()
		if !testCfg.HasTestDatabase() {
			return fmt.Errorf("DATABASE_URL_TEST is not set")
		}
		databaseURL = testCfg.GetTestDatabaseURL()
	} else {
		cfg := config.Load()
		if cfg.DatabaseURL == "" {
			return fmt.Errorf("DATABASE_URL is not set")
		}
		databaseURL = cfg.DatabaseURL
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Hour)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	if os.Getenv("LOG_LEVEL") != "error" {
		log.Println("Database connection established")
	}
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}
	return DB.Ping()
}

func GetStats() sql.DBStats {
	if DB == nil {
		return sql.DBStats{}
	}
	return DB.Stats()
}
