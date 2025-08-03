package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valeriapadilla/stock-insights/internal/errors"
	"github.com/valeriapadilla/stock-insights/internal/model"
	serviceInterfaces "github.com/valeriapadilla/stock-insights/internal/service/interfaces"
)

// Mock implementations
type MockStockService struct {
	mock.Mock
}

func (m *MockStockService) ListStocks(limit, offset int, sort, order string) ([]*model.Stock, int, error) {
	args := m.Called(limit, offset, sort, order)
	return args.Get(0).([]*model.Stock), args.Int(1), args.Error(2)
}

func (m *MockStockService) GetStock(ticket string) (*model.Stock, error) {
	args := m.Called(ticket)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Stock), args.Error(1)
}

func (m *MockStockService) SearchStocks(params serviceInterfaces.StockSearchParams) ([]*model.Stock, int, error) {
	args := m.Called(params)
	return args.Get(0).([]*model.Stock), args.Int(1), args.Error(2)
}

// Tests
func TestStocksHandler_ListStocks(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockStocks     []*model.Stock
		mockCount      int
		expectedStatus int
		setupMocks     func(*MockStockService)
	}{
		{
			name:        "successful list",
			queryParams: "?limit=5&offset=0",
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
			mockCount:      2,
			expectedStatus: http.StatusOK,
			setupMocks: func(service *MockStockService) {
				service.On("ListStocks", 5, 0, "time", "desc").Return([]*model.Stock{
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
				}, 2, nil)
			},
		},
		{
			name:           "service error",
			queryParams:    "?limit=5&offset=0",
			mockStocks:     []*model.Stock{},
			mockCount:      0,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(service *MockStockService) {
				service.On("ListStocks", 5, 0, "time", "desc").Return([]*model.Stock{}, 0, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			mockService := &MockStockService{}
			tt.setupMocks(mockService)

			handler := &StocksHandler{
				stockService: mockService,
				logger:       logrus.New(),
			}

			// Create request
			req, _ := http.NewRequest("GET", "/api/v1/public/stocks"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.ListStocks(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mocks
			mockService.AssertExpectations(t)
		})
	}
}

func TestStocksHandler_GetStock(t *testing.T) {
	tests := []struct {
		name           string
		ticket         string
		mockStock      *model.Stock
		expectedStatus int
		setupMocks     func(*MockStockService)
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
			expectedStatus: http.StatusOK,
			setupMocks: func(service *MockStockService) {
				service.On("GetStock", "AAPL").Return(&model.Stock{
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
			name:           "stock not found",
			ticket:         "INVALID",
			mockStock:      nil,
			expectedStatus: http.StatusNotFound,
			setupMocks: func(service *MockStockService) {
				service.On("GetStock", "INVALID").Return(nil, errors.NewNotFoundError("stock not found", nil))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			mockService := &MockStockService{}
			tt.setupMocks(mockService)

			handler := &StocksHandler{
				stockService: mockService,
				logger:       logrus.New(),
			}

			// Create request
			req, _ := http.NewRequest("GET", "/api/v1/public/stocks/"+tt.ticket, nil)
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "ticket", Value: tt.ticket}}

			// Execute
			handler.GetStock(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mocks
			mockService.AssertExpectations(t)
		})
	}
}

func TestStocksHandler_SearchStocks(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockStocks     []*model.Stock
		mockCount      int
		expectedStatus int
		setupMocks     func(*MockStockService)
	}{
		{
			name:        "successful search",
			queryParams: "?ticket=AAPL&limit=5",
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
			mockCount:      1,
			expectedStatus: http.StatusOK,
			setupMocks: func(service *MockStockService) {
				service.On("SearchStocks", mock.Anything).Return([]*model.Stock{
					{
						Ticker:     "AAPL",
						Company:    "Apple Inc",
						Action:     "target raised by",
						RatingTo:   "Buy",
						TargetTo:   "$200.00",
						TargetFrom: "$150.00",
						Time:       time.Now(),
					},
				}, 1, nil)
			},
		},
		{
			name:           "service error",
			queryParams:    "?ticket=INVALID&limit=5",
			mockStocks:     []*model.Stock{},
			mockCount:      0,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(service *MockStockService) {
				service.On("SearchStocks", mock.Anything).Return([]*model.Stock{}, 0, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			mockService := &MockStockService{}
			tt.setupMocks(mockService)

			handler := &StocksHandler{
				stockService: mockService,
				logger:       logrus.New(),
			}

			// Create request
			req, _ := http.NewRequest("GET", "/api/v1/public/stocks/search"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.SearchStocks(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mocks
			mockService.AssertExpectations(t)
		})
	}
}
