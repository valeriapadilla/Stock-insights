package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

// StockRepository implements interfaces.StockRepository
type StockRepository struct {
	*BaseRepository
	validator *validator.CommonValidator
}

// Ensure StockRepository implements the interface
var _ interfaces.StockRepository = (*StockRepository)(nil)

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{
		BaseRepository: NewBaseRepository(db),
		validator:      validator.NewCommonValidator(),
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

	if filters["brokerage"] != "" {
		qb.Where("brokerage = ?", r.validator.SanitizeString(filters["brokerage"]))
	}
	if filters["rating"] != "" {
		qb.Where("rating_to = ?", r.validator.SanitizeString(filters["rating"]))
	}
	if filters["action"] != "" {
		qb.Where("action = ?", r.validator.SanitizeString(filters["action"]))
	}

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

	return stocks, nil
}

func (r *StockRepository) Count(filters map[string]string) (int, error) {
	qb := NewQueryBuilder().From("stocks")

	if filters["brokerage"] != "" {
		qb.Where("brokerage = ?", r.validator.SanitizeString(filters["brokerage"]))
	}
	if filters["rating"] != "" {
		qb.Where("rating_to = ?", r.validator.SanitizeString(filters["rating"]))
	}
	if filters["action"] != "" {
		qb.Where("action = ?", r.validator.SanitizeString(filters["action"]))
	}

	query, args := qb.CountQuery()

	var count int
	err := r.GetDB().QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count stocks: %w", err)
	}

	return count, nil
}

func (r *StockRepository) GetLastUpdateTime() (*time.Time, error) {
	query := "SELECT MAX(updated_at) FROM stocks"

	var updatedAt *time.Time
	err := r.GetDB().QueryRow(query).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last update time: %w", err)
	}

	return updatedAt, nil
}
