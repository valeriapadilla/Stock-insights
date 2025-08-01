package interfaces

import (
	"github.com/valeriapadilla/stock-insights/internal/model"
)

type StockCommand interface {
	Create(stock *model.Stock) error

	BulkCreate(stocks []*model.Stock) error

	Upsert(stock *model.Stock) error

	BulkUpsert(stocks []*model.Stock) error
}
