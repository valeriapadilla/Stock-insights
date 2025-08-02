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

func (r *StockRepository) GetStocks(params interfaces.GetStocksParams) ([]*model.Stock, error) {
	if params.Limit <= 0 {
		params.Limit = 50
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	if params.Sort == "" {
		params.Sort = "time"
	}
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "desc"
	}

	var query string
	var args []interface{}
	argIndex := 1

	if params.Search != nil {
		query = `
			SELECT ticker, company, target_from, target_to, rating_from, rating_to, 
			       action, brokerage, time, created_at, updated_at
			FROM stocks 
			WHERE 1=1
		`

		if params.Search.Ticket != "" {
			query += fmt.Sprintf(" AND ticker ILIKE $%d", argIndex)
			args = append(args, "%"+params.Search.Ticket+"%")
			argIndex++
		}

		if params.Search.DateFrom != nil {
			query += fmt.Sprintf(" AND time >= $%d", argIndex)
			args = append(args, params.Search.DateFrom)
			argIndex++
		}

		if params.Search.DateTo != nil {
			query += fmt.Sprintf(" AND time <= $%d", argIndex)
			args = append(args, params.Search.DateTo)
			argIndex++
		}

		if params.Search.MinPrice != nil {
			query += fmt.Sprintf(" AND CAST(REPLACE(target_to, '$', '') AS DECIMAL) >= $%d", argIndex)
			args = append(args, *params.Search.MinPrice)
			argIndex++
		}

		if params.Search.MaxPrice != nil {
			query += fmt.Sprintf(" AND CAST(REPLACE(target_to, '$', '') AS DECIMAL) <= $%d", argIndex)
			args = append(args, *params.Search.MaxPrice)
			argIndex++
		}
	} else if len(params.Filters) > 0 {
		qb := NewQueryBuilder().
			Select("ticker", "company", "target_from", "target_to",
				"rating_from", "rating_to", "action", "brokerage", "time",
				"created_at", "updated_at").
			From("stocks")

		if params.Filters["brokerage"] != "" {
			qb.Where("brokerage = ?", r.validator.SanitizeString(params.Filters["brokerage"]))
		}
		if params.Filters["rating"] != "" {
			qb.Where("rating_to = ?", r.validator.SanitizeString(params.Filters["rating"]))
		}
		if params.Filters["action"] != "" {
			qb.Where("action = ?", r.validator.SanitizeString(params.Filters["action"]))
		}

		query, args = qb.Build()
	} else {
		query = fmt.Sprintf(`
			SELECT ticker, company, target_from, target_to, rating_from, rating_to, 
			       action, brokerage, time, created_at, updated_at
			FROM stocks 
			ORDER BY %s %s 
			LIMIT %d OFFSET %d
		`, params.Sort, params.Order, params.Limit, params.Offset)
	}

	if params.Search != nil || len(params.Filters) > 0 {
		query += fmt.Sprintf(" ORDER BY %s %s", params.Sort, params.Order)
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, params.Offset)
	}

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

func (r *StockRepository) GetStocksCount(params interfaces.GetStocksParams) (int, error) {
	if params.Search != nil {
		query := "SELECT COUNT(*) FROM stocks WHERE 1=1"
		var args []interface{}
		argIndex := 1

		if params.Search.Ticket != "" {
			query += fmt.Sprintf(" AND ticker ILIKE $%d", argIndex)
			args = append(args, "%"+params.Search.Ticket+"%")
			argIndex++
		}

		if params.Search.DateFrom != nil {
			query += fmt.Sprintf(" AND time >= $%d", argIndex)
			args = append(args, params.Search.DateFrom)
			argIndex++
		}

		if params.Search.DateTo != nil {
			query += fmt.Sprintf(" AND time <= $%d", argIndex)
			args = append(args, params.Search.DateTo)
			argIndex++
		}

		if params.Search.MinPrice != nil {
			query += fmt.Sprintf(" AND CAST(REPLACE(target_to, '$', '') AS DECIMAL) >= $%d", argIndex)
			args = append(args, *params.Search.MinPrice)
			argIndex++
		}

		if params.Search.MaxPrice != nil {
			query += fmt.Sprintf(" AND CAST(REPLACE(target_to, '$', '') AS DECIMAL) <= $%d", argIndex)
			args = append(args, *params.Search.MaxPrice)
			argIndex++
		}

		var count int
		err := r.GetDB().QueryRow(query, args...).Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("failed to count stocks: %w", err)
		}

		return count, nil
	} else if len(params.Filters) > 0 {
		qb := NewQueryBuilder().
			Select("COUNT(*)").
			From("stocks")

		if params.Filters["brokerage"] != "" {
			qb.Where("brokerage = ?", r.validator.SanitizeString(params.Filters["brokerage"]))
		}
		if params.Filters["rating"] != "" {
			qb.Where("rating_to = ?", r.validator.SanitizeString(params.Filters["rating"]))
		}
		if params.Filters["action"] != "" {
			qb.Where("action = ?", r.validator.SanitizeString(params.Filters["action"]))
		}

		query, args := qb.Build()

		var count int
		err := r.GetDB().QueryRow(query, args...).Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("failed to count stocks: %w", err)
		}

		return count, nil
	} else {
		query := "SELECT COUNT(*) FROM stocks"

		var count int
		err := r.GetDB().QueryRow(query).Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("failed to count stocks: %w", err)
		}

		return count, nil
	}
}

func (r *StockRepository) GetStockByTicket(ticket string) (*model.Stock, error) {
	query := `
		SELECT ticker, company, target_from, target_to, rating_from, rating_to, 
		       action, brokerage, time, created_at, updated_at
		FROM stocks 
		WHERE ticker = $1 
		ORDER BY time DESC 
		LIMIT 1
	`

	var stock model.Stock
	err := r.GetDB().QueryRow(query, ticket).Scan(
		&stock.Ticker, &stock.Company, &stock.TargetFrom, &stock.TargetTo,
		&stock.RatingFrom, &stock.RatingTo, &stock.Action, &stock.Brokerage, &stock.Time,
		&stock.CreatedAt, &stock.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock by ticket: %w", err)
	}

	return &stock, nil
}

func (r *StockRepository) SearchStocks(filters interfaces.StockSearchFilters) ([]*model.Stock, error) {
	query := `
		SELECT ticker, company, target_from, target_to, rating_from, rating_to, 
		       action, brokerage, time, created_at, updated_at
		FROM stocks 
		WHERE 1=1
	`
	var args []interface{}
	argIndex := 1

	if filters.Ticket != "" {
		query += fmt.Sprintf(" AND ticker ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Ticket+"%")
		argIndex++
	}

	if filters.DateFrom != nil {
		query += fmt.Sprintf(" AND time >= $%d", argIndex)
		args = append(args, filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		query += fmt.Sprintf(" AND time <= $%d", argIndex)
		args = append(args, filters.DateTo)
		argIndex++
	}

	if filters.MinPrice != nil {
		query += fmt.Sprintf(" AND CAST(REPLACE(target_to, '$', '') AS DECIMAL) >= $%d", argIndex)
		args = append(args, *filters.MinPrice)
		argIndex++
	}

	if filters.MaxPrice != nil {
		query += fmt.Sprintf(" AND CAST(REPLACE(target_to, '$', '') AS DECIMAL) <= $%d", argIndex)
		args = append(args, *filters.MaxPrice)
		argIndex++
	}

	query += " ORDER BY time DESC"
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", filters.Limit, filters.Offset)

	rows, err := r.GetDB().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search stocks: %w", err)
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

func (r *StockRepository) GetLastUpdateTime() (*time.Time, error) {
	query := "SELECT MAX(time) FROM stocks"

	var updatedAt *time.Time
	err := r.GetDB().QueryRow(query).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last update time: %w", err)
	}

	// Check if the result is NULL (zero value) or empty table
	if updatedAt == nil || updatedAt.IsZero() {
		return nil, nil
	}

	return updatedAt, nil
}

func (r *StockRepository) ExistsByTicker(ticker string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM stocks WHERE ticker = $1 LIMIT 1)"

	var exists bool
	err := r.GetDB().QueryRow(query, ticker).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if ticker exists: %w", err)
	}

	return exists, nil
}
