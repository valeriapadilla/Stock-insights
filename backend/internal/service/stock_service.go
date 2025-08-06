package service

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/errors"
	"github.com/valeriapadilla/stock-insights/internal/model"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/service/interfaces"
)

type StockService struct {
	stockRepo repoInterfaces.StockRepository
	logger    *logrus.Logger
}

var _ interfaces.StockServiceInterface = (*StockService)(nil)

func NewStockService(stockRepo repoInterfaces.StockRepository, logger *logrus.Logger) *StockService {
	return &StockService{
		stockRepo: stockRepo,
		logger:    logger,
	}
}

func (s *StockService) ListStocks(limit, offset int, sort, order string) ([]*model.Stock, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	params := repoInterfaces.GetStocksParams{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Order:  order,
	}

	stocks, err := s.stockRepo.GetStocks(params)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get stocks from repository")
		return nil, 0, errors.NewDatabaseError("failed to retrieve stocks", err)
	}

	total, err := s.stockRepo.GetStocksCount(params)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get stocks count from repository")
		return nil, 0, errors.NewDatabaseError("failed to get stocks count", err)
	}

	return stocks, total, nil
}

func (s *StockService) GetStock(ticket string) (*model.Stock, error) {
	if ticket == "" {
		return nil, errors.NewValidationError("ticket is required", nil)
	}

	stock, err := s.stockRepo.GetStockByTicket(ticket)
	if err != nil {
		s.logger.WithError(err).WithField("ticket", ticket).Error("Failed to get stock from repository")
		return nil, errors.NewDatabaseError("failed to retrieve stock", err)
	}

	if stock == nil {
		return nil, errors.NewNotFoundError("stock not found", nil)
	}

	return stock, nil
}

func (s *StockService) SearchStocks(params interfaces.StockSearchParams) ([]*model.Stock, int, error) {
	if params.Limit <= 0 {
		params.Limit = 50
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	var dateFrom, dateTo *time.Time
	if params.DateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", params.DateFrom); err == nil {
			dateFrom = &parsed
		}
	}

	if params.DateTo != "" {
		if parsed, err := time.Parse("2006-01-02", params.DateTo); err == nil {
			endOfDay := parsed.Add(24*time.Hour - time.Nanosecond)
			dateTo = &endOfDay
		}
	}

	repoParams := repoInterfaces.GetStocksParams{
		Limit:  params.Limit,
		Offset: params.Offset,
		Sort:   "time",
		Order:  "desc",
		Search: &repoInterfaces.StockSearchFilters{
			Ticket:   params.Ticket,
			DateFrom: dateFrom,
			DateTo:   dateTo,
			MinPrice: params.MinPrice,
			MaxPrice: params.MaxPrice,
		},
	}

	stocks, err := s.stockRepo.GetStocks(repoParams)
	if err != nil {
		s.logger.WithError(err).Error("Failed to search stocks from repository")
		return nil, 0, errors.NewDatabaseError("failed to search stocks", err)
	}

	total, err := s.stockRepo.GetStocksCount(repoParams)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get stocks count from repository")
		return nil, 0, errors.NewDatabaseError("failed to get stocks count", err)
	}

	return stocks, total, nil
}
