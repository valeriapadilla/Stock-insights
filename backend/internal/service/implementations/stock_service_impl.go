package implementations

import (
	"fmt"

	"github.com/valeriapadilla/stock-insights/internal/model"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	serviceInterfaces "github.com/valeriapadilla/stock-insights/internal/service/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type StockServiceImpl struct {
	stockRepo repoInterfaces.StockRepository
	validator *validator.CommonValidator
}

var _ serviceInterfaces.StockService = (*StockServiceImpl)(nil)

func NewStockService(stockRepo repoInterfaces.StockRepository) *StockServiceImpl {
	return &StockServiceImpl{
		stockRepo: stockRepo,
		validator: validator.NewCommonValidator(),
	}
}

func (s *StockServiceImpl) GetStocksWithFilters(limit, offset int, filters map[string]string) ([]*model.Stock, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, fmt.Errorf("invalid filters: %w", err)
	}

	limit = s.applyPaginationLimits(limit)
	offset = s.applyOffsetLimits(offset)

	stocks, err := s.stockRepo.GetAll(limit, offset, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks: %w", err)
	}

	return s.applyBusinessTransformations(stocks), nil
}

func (s *StockServiceImpl) GetStocksCount(filters map[string]string) (int, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return 0, fmt.Errorf("invalid filters: %w", err)
	}

	count, err := s.stockRepo.Count(filters)
	if err != nil {
		return 0, fmt.Errorf("failed to count stocks: %w", err)
	}

	return count, nil
}

func (s *StockServiceImpl) ValidateFilters(filters map[string]string) error {
	if filters == nil {
		return nil
	}

	for key, value := range filters {
		if err := s.validator.ValidateFilter(key, value); err != nil {
			return fmt.Errorf("filter '%s': %w", key, err)
		}
	}

	return nil
}

func (s *StockServiceImpl) applyPaginationLimits(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func (s *StockServiceImpl) applyOffsetLimits(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func (s *StockServiceImpl) applyBusinessTransformations(stocks []*model.Stock) []*model.Stock {
	if stocks == nil {
		return make([]*model.Stock, 0)
	}
	return stocks
}
