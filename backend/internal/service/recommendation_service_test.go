package service

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

// Tests
func TestRecommendationService_CalculateRecommendations(t *testing.T) {
	tests := []struct {
		name          string
		params        validator.RecommendationParams
		mockStocks    []*model.Stock
		expectedCount int
		expectedError bool
		setupMocks    func(*MockStockRepository, *MockRecommendationCommand)
	}{
		{
			name: "successful calculation with valid stocks",
			params: validator.RecommendationParams{
				DaysBack:   7,
				MaxResults: 10,
				MinScore:   80,
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
			expectedCount: 2,
			expectedError: false,
			setupMocks: func(stockRepo *MockStockRepository, recCmd *MockRecommendationCommand) {
				recCmd.On("DeleteAllRecommendations").Return(nil)
				stockRepo.On("GetStocksCount", mock.Anything).Return(2, nil)
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
			},
		},
		{
			name: "no stocks found",
			params: validator.RecommendationParams{
				DaysBack:   7,
				MaxResults: 10,
				MinScore:   80,
			},
			mockStocks:    []*model.Stock{},
			expectedCount: 0,
			expectedError: false,
			setupMocks: func(stockRepo *MockStockRepository, recCmd *MockRecommendationCommand) {
				recCmd.On("DeleteAllRecommendations").Return(nil)
				stockRepo.On("GetStocksCount", mock.Anything).Return(0, nil)
				stockRepo.On("GetStocks", mock.Anything).Return([]*model.Stock{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockStockRepo := &MockStockRepository{}
			mockRecRepo := &MockRecommendationRepository{}
			mockRecCmd := &MockRecommendationCommand{}

			tt.setupMocks(mockStockRepo, mockRecCmd)

			// Create service with real validator and logger
			service := &RecommendationService{
				stockRepo:          mockStockRepo,
				recommendationRepo: mockRecRepo,
				recommendationCmd:  mockRecCmd,
				validator:          validator.NewRecommendationValidator(),
				logger:             logrus.New(),
				scoringConfig:      DefaultScoringConfig(),
			}

			// Execute
			recommendations, err := service.CalculateRecommendations(tt.params)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, recommendations, tt.expectedCount)
			}

			// Verify mocks
			mockStockRepo.AssertExpectations(t)
			mockRecCmd.AssertExpectations(t)
		})
	}
}

func TestRecommendationService_GetLatestRecommendations(t *testing.T) {
	tests := []struct {
		name                string
		limit               int
		mockRecommendations []*model.Recommendation
		expectedCount       int
		expectedError       bool
		setupMocks          func(*MockRecommendationRepository)
	}{
		{
			name:  "successful retrieval",
			limit: 5,
			mockRecommendations: []*model.Recommendation{
				{Ticker: "AAPL", Score: 95},
				{Ticker: "GOOGL", Score: 90},
			},
			expectedCount: 2,
			expectedError: false,
			setupMocks: func(recRepo *MockRecommendationRepository) {
				recRepo.On("GetLatest", 5).Return([]*model.Recommendation{
					{Ticker: "AAPL", Score: 95},
					{Ticker: "GOOGL", Score: 90},
				}, nil)
			},
		},
		{
			name:                "repository error",
			limit:               5,
			mockRecommendations: []*model.Recommendation{},
			expectedCount:       0,
			expectedError:       true,
			setupMocks: func(recRepo *MockRecommendationRepository) {
				recRepo.On("GetLatest", 5).Return([]*model.Recommendation{}, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockStockRepo := &MockStockRepository{}
			mockRecRepo := &MockRecommendationRepository{}
			mockRecCmd := &MockRecommendationCommand{}

			tt.setupMocks(mockRecRepo)

			// Create service
			service := &RecommendationService{
				stockRepo:          mockStockRepo,
				recommendationRepo: mockRecRepo,
				recommendationCmd:  mockRecCmd,
				validator:          validator.NewRecommendationValidator(),
				logger:             logrus.New(),
				scoringConfig:      DefaultScoringConfig(),
			}

			// Execute
			recommendations, err := service.GetLatestRecommendations(tt.limit)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, recommendations, tt.expectedCount)
			}

			// Verify mocks
			mockRecRepo.AssertExpectations(t)
		})
	}
}

func TestRecommendationService_SaveRecommendations(t *testing.T) {
	tests := []struct {
		name            string
		recommendations []*model.Recommendation
		expectedError   bool
		setupMocks      func(*MockRecommendationRepository)
	}{
		{
			name: "successful save",
			recommendations: []*model.Recommendation{
				{Ticker: "AAPL", Score: 95},
				{Ticker: "GOOGL", Score: 90},
			},
			expectedError: false,
			setupMocks: func(recRepo *MockRecommendationRepository) {
				recRepo.On("CreateRecommendation", mock.Anything).Return(nil).Times(2)
			},
		},
		{
			name: "repository error",
			recommendations: []*model.Recommendation{
				{Ticker: "AAPL", Score: 95},
			},
			expectedError: true,
			setupMocks: func(recRepo *MockRecommendationRepository) {
				recRepo.On("CreateRecommendation", mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockStockRepo := &MockStockRepository{}
			mockRecRepo := &MockRecommendationRepository{}
			mockRecCmd := &MockRecommendationCommand{}

			tt.setupMocks(mockRecRepo)

			// Create service
			service := &RecommendationService{
				stockRepo:          mockStockRepo,
				recommendationRepo: mockRecRepo,
				recommendationCmd:  mockRecCmd,
				validator:          validator.NewRecommendationValidator(),
				logger:             logrus.New(),
				scoringConfig:      DefaultScoringConfig(),
			}

			// Execute
			err := service.SaveRecommendations(tt.recommendations)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockRecRepo.AssertExpectations(t)
		})
	}
}
