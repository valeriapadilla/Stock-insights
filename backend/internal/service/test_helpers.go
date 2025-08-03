package service

import (
	"database/sql"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/valeriapadilla/stock-insights/internal/model"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
)

// Mock implementations
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) GetStocks(params repoInterfaces.GetStocksParams) ([]*model.Stock, error) {
	args := m.Called(params)
	return args.Get(0).([]*model.Stock), args.Error(1)
}

func (m *MockStockRepository) GetStocksCount(params repoInterfaces.GetStocksParams) (int, error) {
	args := m.Called(params)
	return args.Int(0), args.Error(1)
}

func (m *MockStockRepository) GetLastUpdateTime() (*time.Time, error) {
	args := m.Called()
	return args.Get(0).(*time.Time), args.Error(1)
}

func (m *MockStockRepository) ExistsByTicker(ticker string) (bool, error) {
	args := m.Called(ticker)
	return args.Bool(0), args.Error(1)
}

func (m *MockStockRepository) GetStockByTicket(ticket string) (*model.Stock, error) {
	args := m.Called(ticket)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Stock), args.Error(1)
}

func (m *MockStockRepository) SearchStocks(filters repoInterfaces.StockSearchFilters) ([]*model.Stock, error) {
	args := m.Called(filters)
	return args.Get(0).([]*model.Stock), args.Error(1)
}

func (m *MockStockRepository) GetDB() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}

type MockRecommendationRepository struct {
	mock.Mock
}

func (m *MockRecommendationRepository) GetLatest(limit int) ([]*model.Recommendation, error) {
	args := m.Called(limit)
	return args.Get(0).([]*model.Recommendation), args.Error(1)
}

func (m *MockRecommendationRepository) GetRecommendationsByDate(date time.Time) ([]*model.Recommendation, error) {
	args := m.Called(date)
	return args.Get(0).([]*model.Recommendation), args.Error(1)
}

func (m *MockRecommendationRepository) GetLatestRunAt() (*time.Time, error) {
	args := m.Called()
	return args.Get(0).(*time.Time), args.Error(1)
}

func (m *MockRecommendationRepository) DeleteOldRecommendations(days int) error {
	args := m.Called(days)
	return args.Error(0)
}

func (m *MockRecommendationRepository) GetRecommendationCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockRecommendationRepository) CreateRecommendation(recommendation *model.Recommendation) error {
	args := m.Called(recommendation)
	return args.Error(0)
}

func (m *MockRecommendationRepository) GetDB() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}

type MockRecommendationCommand struct {
	mock.Mock
}

func (m *MockRecommendationCommand) BulkCreate(recommendations []*model.Recommendation) error {
	args := m.Called(recommendations)
	return args.Error(0)
}

func (m *MockRecommendationCommand) DeleteAllRecommendations() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRecommendationCommand) GetDB() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}
