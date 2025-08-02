package interfaces

import (
	"database/sql"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type StockSearchFilters struct {
	Ticket   string     `json:"ticket"`
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
	MinPrice *float64   `json:"min_price"`
	MaxPrice *float64   `json:"max_price"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

type GetStocksParams struct {
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	Sort    string              `json:"sort"`
	Order   string              `json:"order"`
	Filters map[string]string   `json:"filters"`
	Search  *StockSearchFilters `json:"search"`
}

type StockRepository interface {
	GetStocks(params GetStocksParams) ([]*model.Stock, error)

	GetStocksCount(params GetStocksParams) (int, error)

	GetLastUpdateTime() (*time.Time, error)

	ExistsByTicker(ticker string) (bool, error)

	GetStockByTicket(ticket string) (*model.Stock, error)

	SearchStocks(filters StockSearchFilters) ([]*model.Stock, error)

	GetDB() *sql.DB
}
