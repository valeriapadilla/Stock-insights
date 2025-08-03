package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

// Mock implementations
type MockRecommendationService struct {
	mock.Mock
}

func (m *MockRecommendationService) CalculateRecommendations(params validator.RecommendationParams) ([]*model.Recommendation, error) {
	args := m.Called(params)
	return args.Get(0).([]*model.Recommendation), args.Error(1)
}

func (m *MockRecommendationService) GetLatestRecommendations(limit int) ([]*model.Recommendation, error) {
	args := m.Called(limit)
	return args.Get(0).([]*model.Recommendation), args.Error(1)
}

func (m *MockRecommendationService) SaveRecommendations(recommendations []*model.Recommendation) error {
	args := m.Called(recommendations)
	return args.Error(0)
}

// Tests
func TestRecommendationsHandler_GetRecommendations(t *testing.T) {
	tests := []struct {
		name                string
		queryParams         string
		mockRecommendations []*model.Recommendation
		expectedStatus      int
		expectedBody        map[string]interface{}
		setupMocks          func(*MockRecommendationService)
	}{
		{
			name:        "successful retrieval",
			queryParams: "?limit=5",
			mockRecommendations: []*model.Recommendation{
				{
					ID:          "1",
					Ticker:      "AAPL",
					Score:       95,
					Explanation: "Test explanation",
					RunAt:       time.Now(),
					Rank:        1,
				},
				{
					ID:          "2",
					Ticker:      "GOOGL",
					Score:       90,
					Explanation: "Test explanation 2",
					RunAt:       time.Now(),
					Rank:        2,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"limit":           5,
				"recommendations": []interface{}{},
				"total":           2,
			},
			setupMocks: func(service *MockRecommendationService) {
				service.On("GetLatestRecommendations", 5).Return([]*model.Recommendation{
					{
						ID:          "1",
						Ticker:      "AAPL",
						Score:       95,
						Explanation: "Test explanation",
						RunAt:       time.Now(),
						Rank:        1,
					},
					{
						ID:          "2",
						Ticker:      "GOOGL",
						Score:       90,
						Explanation: "Test explanation 2",
						RunAt:       time.Now(),
						Rank:        2,
					},
				}, nil)
			},
		},
		{
			name:                "service error",
			queryParams:         "?limit=5",
			mockRecommendations: []*model.Recommendation{},
			expectedStatus:      http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error":   "Internal server error",
				"message": "Failed to retrieve recommendations",
			},
			setupMocks: func(service *MockRecommendationService) {
				service.On("GetLatestRecommendations", 5).Return([]*model.Recommendation{}, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			mockService := &MockRecommendationService{}
			tt.setupMocks(mockService)

			handler := &RecommendationsHandler{
				recommendationService: mockService,
				logger:                logrus.New(),
			}

			// Create request
			req, _ := http.NewRequest("GET", "/api/v1/public/recommendations"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.GetRecommendations(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verify mocks
			mockService.AssertExpectations(t)
		})
	}
}

func TestRecommendationsHandler_CalculateRecommendations(t *testing.T) {
	tests := []struct {
		name                string
		requestBody         map[string]interface{}
		mockRecommendations []*model.Recommendation
		expectedStatus      int
		expectedBody        map[string]interface{}
		setupMocks          func(*MockRecommendationService)
	}{
		{
			name: "successful calculation",
			requestBody: map[string]interface{}{
				"days_back":   7,
				"max_results": 10,
				"min_score":   80,
			},
			mockRecommendations: []*model.Recommendation{
				{
					ID:          "1",
					Ticker:      "AAPL",
					Score:       95,
					Explanation: "Test explanation",
					RunAt:       time.Now(),
					Rank:        1,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message":         "Recommendations calculated and saved successfully",
				"recommendations": []interface{}{},
				"run_at":          "",
				"total":           1,
			},
			setupMocks: func(service *MockRecommendationService) {
				service.On("CalculateRecommendations", mock.Anything).Return([]*model.Recommendation{
					{
						ID:          "1",
						Ticker:      "AAPL",
						Score:       95,
						Explanation: "Test explanation",
						RunAt:       time.Now(),
						Rank:        1,
					},
				}, nil)
				service.On("SaveRecommendations", mock.Anything).Return(nil)
			},
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"days_back":   7,
				"max_results": 10,
				"min_score":   80,
			},
			mockRecommendations: []*model.Recommendation{},
			expectedStatus:      http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error":   "Internal server error",
				"message": "Failed to calculate recommendations",
			},
			setupMocks: func(service *MockRecommendationService) {
				service.On("CalculateRecommendations", mock.Anything).Return([]*model.Recommendation{}, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			mockService := &MockRecommendationService{}
			tt.setupMocks(mockService)

			handler := &RecommendationsHandler{
				recommendationService: mockService,
				logger:                logrus.New(),
			}

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/admin/recommendations/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.CalculateRecommendations(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verify mocks
			mockService.AssertExpectations(t)
		})
	}
}
