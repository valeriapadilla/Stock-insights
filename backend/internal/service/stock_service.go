package service

import (
	"fmt"
	"sort"
	"strings"
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

	s.calculateChangePercentForStocks(stocks)

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

	s.calculateChangePercentForStock(stock)

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

	limitForDB := params.Limit
	if params.Rating != "" {
		limitForDB = params.Limit * 10
	}

	repoParams := repoInterfaces.GetStocksParams{
		Limit:  limitForDB,
		Offset: params.Offset,
		Sort:   "time",
		Order:  "desc",
	}

	// If rating is provided, use Filters approach
	if params.Rating != "" {
		repoParams.Filters = map[string]string{
			"rating": params.Rating,
		}
	} else {
		// Use Search approach for other filters
		repoParams.Search = &repoInterfaces.StockSearchFilters{
			Ticket:   params.Ticket,
			DateFrom: dateFrom,
			DateTo:   dateTo,
			MinPrice: params.MinPrice,
			MaxPrice: params.MaxPrice,
		}
	}

	stocks, err := s.stockRepo.GetStocks(repoParams)
	if err != nil {
		s.logger.WithError(err).Error("Failed to search stocks from repository")
		return nil, 0, errors.NewDatabaseError("failed to search stocks", err)
	}

	filteredStocks := s.applyFilters(stocks, params)
	sortedStocks := s.sortStocks(filteredStocks, params.SortBy, params.Order)

	if len(sortedStocks) > params.Limit {
		sortedStocks = sortedStocks[:params.Limit]
	}

	total, err := s.stockRepo.GetStocksCount(repoParams)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get stocks count from repository")
		return nil, 0, errors.NewDatabaseError("failed to get stocks count", err)
	}

	return sortedStocks, total, nil
}

func (s *StockService) applyFilters(stocks []*model.Stock, params interfaces.StockSearchParams) []*model.Stock {
	if params.Rating == "" {
		return stocks
	}

	filtered := make([]*model.Stock, 0)
	ratingLower := strings.ToLower(params.Rating)

	for _, stock := range stocks {
		if stock.GetRating() == ratingLower {
			filtered = append(filtered, stock)
		}
	}

	return filtered
}

func (s *StockService) sortStocks(stocks []*model.Stock, sortBy, order string) []*model.Stock {
	if len(stocks) == 0 {
		return stocks
	}

	s.calculateChangePercentForStocks(stocks)

	sort.Slice(stocks, func(i, j int) bool {
		var result bool

		switch sortBy {
		case "ticker":
			result = stocks[i].Ticker < stocks[j].Ticker
		case "change_percent":
			result = stocks[i].GetChangePercentage() < stocks[j].GetChangePercentage()
		case "time":
			result = stocks[i].Time.Before(stocks[j].Time)
		default:
			result = stocks[i].Time.Before(stocks[j].Time)
		}

		if order == "desc" {
			result = !result
		}

		return result
	})

	return stocks
}

func (s *StockService) calculateChangePercentForStocks(stocks []*model.Stock) {
	for _, stock := range stocks {
		s.calculateChangePercentForStock(stock)
	}
}

func (s *StockService) calculateChangePercentForStock(stock *model.Stock) {
	change := stock.GetChangePercentage()
	stock.ChangePercent = s.formatChangePercent(change)
}

func (s *StockService) formatChangePercent(change float64) string {
	if change == 0 {
		return "0.0%"
	}

	sign := ""
	if change > 0 {
		sign = "+"
	}

	return fmt.Sprintf("%s%.1f%%", sign, change)
}
