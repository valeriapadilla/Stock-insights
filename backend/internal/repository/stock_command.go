package repository

import (
	"database/sql"
	"fmt"

	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type StockCommandImpl struct {
	*BaseRepository
	validator *validator.CommonValidator
}

var _ interfaces.StockCommand = (*StockCommandImpl)(nil)

func NewStockCommand(db *sql.DB) *StockCommandImpl {
	return &StockCommandImpl{
		BaseRepository: NewBaseRepository(db),
		validator:      validator.NewCommonValidator(),
	}
}

func (c *StockCommandImpl) Create(stock *model.Stock) error {
	if err := c.validateStock(stock); err != nil {
		return fmt.Errorf("stock validation failed: %w", err)
	}

	query := `
		INSERT INTO stocks (ticker, company, target_from, target_to, rating_from, rating_to, action, brokerage, time, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := c.GetDB().Exec(query,
		stock.Ticker, stock.Company, stock.TargetFrom, stock.TargetTo,
		stock.RatingFrom, stock.RatingTo, stock.Action, stock.Brokerage, stock.Time,
		stock.CreatedAt, stock.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}

	return nil
}

func (c *StockCommandImpl) BulkCreate(stocks []*model.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := c.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO stocks (ticker, company, target_from, target_to, rating_from, rating_to, action, brokerage, time, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, stock := range stocks {
		if err := c.validateStock(stock); err != nil {
			return fmt.Errorf("stock validation failed for %s: %w", stock.Ticker, err)
		}

		_, err := stmt.Exec(
			stock.Ticker, stock.Company, stock.TargetFrom, stock.TargetTo,
			stock.RatingFrom, stock.RatingTo, stock.Action, stock.Brokerage, stock.Time,
			stock.CreatedAt, stock.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert stock %s: %w", stock.Ticker, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *StockCommandImpl) Upsert(stock *model.Stock) error {
	if err := c.validateStock(stock); err != nil {
		return fmt.Errorf("stock validation failed: %w", err)
	}

	query := `
		INSERT INTO stocks (ticker, company, target_from, target_to, rating_from, rating_to, action, brokerage, time, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (ticker, time) DO UPDATE SET
			company = EXCLUDED.company,
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			action = EXCLUDED.action,
			brokerage = EXCLUDED.brokerage,
			updated_at = EXCLUDED.updated_at
	`

	_, err := c.GetDB().Exec(query,
		stock.Ticker, stock.Company, stock.TargetFrom, stock.TargetTo,
		stock.RatingFrom, stock.RatingTo, stock.Action, stock.Brokerage, stock.Time,
		stock.CreatedAt, stock.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert stock: %w", err)
	}

	return nil
}

func (c *StockCommandImpl) BulkUpsert(stocks []*model.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := c.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO stocks (ticker, company, target_from, target_to, rating_from, rating_to, action, brokerage, time, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (ticker, time) DO UPDATE SET
			company = EXCLUDED.company,
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			action = EXCLUDED.action,
			brokerage = EXCLUDED.brokerage,
			updated_at = EXCLUDED.updated_at
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, stock := range stocks {
		if err := c.validateStock(stock); err != nil {
			return fmt.Errorf("stock validation failed for %s: %w", stock.Ticker, err)
		}

		_, err := stmt.Exec(
			stock.Ticker, stock.Company, stock.TargetFrom, stock.TargetTo,
			stock.RatingFrom, stock.RatingTo, stock.Action, stock.Brokerage, stock.Time,
			stock.CreatedAt, stock.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to upsert stock %s: %w", stock.Ticker, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *StockCommandImpl) validateStock(stock *model.Stock) error {
	if stock == nil {
		return fmt.Errorf("stock cannot be nil")
	}

	if stock.Ticker == "" {
		return fmt.Errorf("ticker is required")
	}
	if stock.Company == "" {
		return fmt.Errorf("company is required")
	}

	stock.Ticker = c.validator.SanitizeString(stock.Ticker)
	stock.Company = c.validator.SanitizeString(stock.Company)
	stock.TargetFrom = c.validator.SanitizeString(stock.TargetFrom)
	stock.TargetTo = c.validator.SanitizeString(stock.TargetTo)
	stock.RatingFrom = c.validator.SanitizeString(stock.RatingFrom)
	stock.RatingTo = c.validator.SanitizeString(stock.RatingTo)
	stock.Action = c.validator.SanitizeString(stock.Action)
	stock.Brokerage = c.validator.SanitizeString(stock.Brokerage)

	if len(stock.Ticker) > 10 {
		return fmt.Errorf("ticker too long (max 10 characters)")
	}

	return nil
}
