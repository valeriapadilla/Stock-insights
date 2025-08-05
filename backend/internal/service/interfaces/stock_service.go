package interfaces

import (
	"github.com/valeriapadilla/stock-insights/internal/model"
)

type StockSearchParams struct {
	Ticket   string   `json:"ticket"`
	Rating   string   `json:"rating"`
	DateFrom string   `json:"date_from"`
	DateTo   string   `json:"date_to"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
}

type StockServiceInterface interface {
	ListStocks(limit, offset int, sort, order string) ([]*model.Stock, int, error)
	GetStock(ticket string) (*model.Stock, error)
	SearchStocks(params StockSearchParams) ([]*model.Stock, int, error)
}
