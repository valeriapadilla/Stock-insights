package database

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
	"github.com/valeriapadilla/stock-insights/internal/errors"
)

func Connect() (*sql.DB, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		return nil, fmt.Errorf("%w: empty string",errors.ErrConfig)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDBConnection, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDBPing, err)
	}

	return db, nil
}
