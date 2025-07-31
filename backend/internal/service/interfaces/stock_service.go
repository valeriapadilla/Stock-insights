package interfaces

import (
	"github.com/valeriapadilla/stock-insights/internal/model"
)

type StockService interface {
	GetStocksWithFilters(limit, offset int, filters map[string]string) ([]*model.Stock, error)

	GetStocksCount(filters map[string]string) (int, error)

	ValidateFilters(filters map[string]string) error
}
