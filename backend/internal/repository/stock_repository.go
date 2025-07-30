package repository

import (
	"database/sql"
	"fmt"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type StockRepository struct {
	*BaseRepository
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *StockRepository) GetAll(limit, offset int, filters map[string]string) ([]*model.Stock, error) {
	qb := NewQueryBuilder().
		Select("ticker", "company", "target_from", "target_to",
			"rating_from", "rating_to", "action", "brokerage", "time",
			"created_at", "updated_at").
		From("stocks").
		OrderBy("time", "DESC").
		Limit(limit).
		Offset(offset)

	qb.Where("brokerage = ?", filters["brokerage"])
	qb.Where("rating_to = ?", filters["rating"])
	qb.Where("action = ?", filters["action"])

	query, args := qb.Build()

	rows, err := r.GetDB().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks: %w", err)
	}
	defer rows.Close()

	var stocks []*model.Stock
	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(
			&stock.Ticker, &stock.Company, &stock.TargetFrom, &stock.TargetTo,
			&stock.RatingFrom, &stock.RatingTo, &stock.Action, &stock.Brokerage, &stock.Time,
			&stock.CreatedAt, &stock.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, &stock)
	}

	if stocks == nil {
		stocks = make([]*model.Stock, 0)
	}

	return stocks, nil
}

func (r *StockRepository) Count(filters map[string]string) (int, error) {
	qb := NewQueryBuilder().From("stocks")

	qb.Where("brokerage = ?", filters["brokerage"])
	qb.Where("rating_to = ?", filters["rating"])
	qb.Where("action = ?", filters["action"])

	query, args := qb.CountQuery()

	var count int
	err := r.GetDB().QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count stocks: %w", err)
	}

	return count, nil
}
