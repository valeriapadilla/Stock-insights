package interfaces

import (
	"database/sql"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type StockRepository interface {
	GetAll(limit, offset int, filters map[string]string) ([]*model.Stock, error)

	Count(filters map[string]string) (int, error)

	GetLastUpdateTime() (*time.Time, error)

	GetDB() *sql.DB
}
