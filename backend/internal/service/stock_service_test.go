package service

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valeriapadilla/stock-insights/internal/model"
	serviceInterfaces "github.com/valeriapadilla/stock-insights/internal/service/interfaces"
)

func TestStockService_ListStocks(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		offset        int
		sort          string
		order         string
		mockStocks    []*model.Stock
		mockCount     int
		expectedCount int
		expectedError bool
		setupMocks    func(*MockStockRepository)
	}{
		{
			name:   "successful list",
			limit:  10,
			offset: 0,
			sort:   "time",
			order:  "desc",
			mockStocks: []*model.Stock{
				{
					Ticker:     "AAPL",
					Company:    "Apple Inc",
					Action:     "target raised by",
					RatingTo:   "Buy",
					TargetTo:   "$200.00",
					TargetFrom: "$150.00",
					Time:       time.Now(),
				},
				{
					Ticker:     "GOOGL",
					Company:    "Alphabet Inc",
					Action:     "target raised by",
					RatingTo:   "Overweight",
					TargetTo:   "$300.00",
					TargetFrom: "$250.00",
					Time:       time.Now(),
				},
			},
			mockCount:     2,
			expectedCount: 2,
			expectedError: false,
			setupMocks: func(stockRepo *MockStockRepository) {
				stockRepo.On("GetStocks", mock.Anything).Return([]*model.Stock{
					{
						Ticker:     "AAPL",
						Company:    "Apple Inc",
						Action:     "target raised by",
						RatingTo:   "Buy",
						TargetTo:   "$200.00",
						TargetFrom: "$150.00",
						Time:       time.Now(),
					},
					{
						Ticker:     "GOOGL",
						Company:    "Alphabet Inc",
						Action:     "target raised by",
						RatingTo:   "Overweight",
						TargetTo:   "$300.00",
						TargetFrom: "$250.00",
						Time:       time.Now(),
					},
				}, nil)
				stockRepo.On("GetStocksCount", mock.Anything).Return(2, nil)
			},
		},
		{
			name:          "repository error",
			limit:         10,
			offset:        0,
			sort:          "time",
			order:         "desc",
			mockStocks:    []*model.Stock{},
			mockCount:     0,
			expectedCount: 0,
			expectedError: true,
			setupMocks: func(stockRepo *MockStockRepository) {
				stockRepo.On("GetStocks", mock.Anything).Return([]*model.Stock{}, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStockRepo := &MockStockRepository{}

			tt.setupMocks(mockStockRepo)

			service := &StockService{
				stockRepo: mockStockRepo,
				logger:    logrus.New(),
			}

			stocks, count, err := service.ListStocks(tt.limit, tt.offset, tt.sort, tt.order)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, stocks, tt.expectedCount)
				assert.Equal(t, tt.mockCount, count)
			}

			mockStockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_GetStock(t *testing.T) {
	tests := []struct {
		name          string
		ticket        string
		mockStock     *model.Stock
		expectedError bool
		setupMocks    func(*MockStockRepository)
	}{
		{
			name:   "successful get",
			ticket: "AAPL",
			mockStock: &model.Stock{
				Ticker:     "AAPL",
				Company:    "Apple Inc",
				Action:     "target raised by",
				RatingTo:   "Buy",
				TargetTo:   "$200.00",
				TargetFrom: "$150.00",
				Time:       time.Now(),
			},
			expectedError: false,
			setupMocks: func(stockRepo *MockStockRepository) {
				stockRepo.On("GetStockByTicket", "AAPL").Return(&model.Stock{
					Ticker:     "AAPL",
					Company:    "Apple Inc",
					Action:     "target raised by",
					RatingTo:   "Buy",
					TargetTo:   "$200.00",
					TargetFrom: "$150.00",
					Time:       time.Now(),
				}, nil)
			},
		},
		{
			name:          "stock not found",
			ticket:        "INVALID",
			mockStock:     nil,
			expectedError: true,
			setupMocks: func(stockRepo *MockStockRepository) {
				stockRepo.On("GetStockByTicket", "INVALID").Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStockRepo := &MockStockRepository{}

			tt.setupMocks(mockStockRepo)

			service := &StockService{
				stockRepo: mockStockRepo,
				logger:    logrus.New(),
			}

			stock, err := service.GetStock(tt.ticket)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, stock)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, stock)
				assert.Equal(t, tt.ticket, stock.Ticker)
			}

			mockStockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_SearchStocks(t *testing.T) {
	tests := []struct {
		name          string
		params        serviceInterfaces.StockSearchParams
		mockStocks    []*model.Stock
		mockCount     int
		expectedCount int
		expectedError bool
		setupMocks    func(*MockStockRepository)
	}{
		{
			name: "successful search",
			params: serviceInterfaces.StockSearchParams{
				Ticket:   "AAPL",
				DateFrom: "2025-01-01",
				DateTo:   "2025-12-31",
				Limit:    10,
				Offset:   0,
			},
			mockStocks: []*model.Stock{
				{
					Ticker:     "AAPL",
					Company:    "Apple Inc",
					Action:     "target raised by",
					RatingTo:   "Buy",
					TargetTo:   "$200.00",
					TargetFrom: "$150.00",
					Time:       time.Now(),
				},
			},
			mockCount:     1,
			expectedCount: 1,
			expectedError: false,
			setupMocks: func(stockRepo *MockStockRepository) {
				stockRepo.On("GetStocks", mock.Anything).Return([]*model.Stock{
					{
						Ticker:     "AAPL",
						Company:    "Apple Inc",
						Action:     "target raised by",
						RatingTo:   "Buy",
						TargetTo:   "$200.00",
						TargetFrom: "$150.00",
						Time:       time.Now(),
					},
				}, nil)
				stockRepo.On("GetStocksCount", mock.Anything).Return(1, nil)
			},
		},
		{
			name: "search error",
			params: serviceInterfaces.StockSearchParams{
				Ticket: "INVALID",
				Limit:  10,
				Offset: 0,
			},
			mockStocks:    []*model.Stock{},
			mockCount:     0,
			expectedCount: 0,
			expectedError: true,
			setupMocks: func(stockRepo *MockStockRepository) {
				stockRepo.On("GetStocks", mock.Anything).Return([]*model.Stock{}, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStockRepo := &MockStockRepository{}

			tt.setupMocks(mockStockRepo)

			service := &StockService{
				stockRepo: mockStockRepo,
				logger:    logrus.New(),
			}

			stocks, count, err := service.SearchStocks(tt.params)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, stocks, tt.expectedCount)
				assert.Equal(t, tt.mockCount, count)
			}

			mockStockRepo.AssertExpectations(t)
		})
	}
}
