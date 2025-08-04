package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valeriapadilla/stock-insights/internal/job"
	"github.com/valeriapadilla/stock-insights/internal/service/interfaces"
)

// Mock implementations
type MockIngestionService struct {
	mock.Mock
}

func (m *MockIngestionService) TriggerIngestionAsync(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockJobManager struct {
	mock.Mock
}

func (m *MockJobManager) CreateJob() (*job.Job, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*job.Job), args.Error(1)
}

func (m *MockJobManager) GetJob(jobID string) (*job.Job, bool) {
	args := m.Called(jobID)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*job.Job), args.Bool(1)
}

func (m *MockJobManager) UpdateJob(jobID string, status job.JobStatus, progress int, message string) {
	m.Called(jobID, status, progress, message)
}

func (m *MockJobManager) SetJobError(jobID string, err error) {
	m.Called(jobID, err)
}

func (m *MockJobManager) RunJobAsync(jobID string, workFunc func(context.Context) error) error {
	args := m.Called(jobID, workFunc)
	return args.Error(0)
}

func (m *MockJobManager) CleanupOldJobs(maxAge time.Duration) {
	m.Called(maxAge)
}

// Tests
func TestStocksIngestionHandler_TriggerIngestion(t *testing.T) {
	tests := []struct {
		name           string
		mockJob        *job.Job
		expectedStatus int
		setupMocks     func(interfaces.IngestionServiceInterface, *MockJobManager)
	}{
		{
			name: "successful trigger",
			mockJob: &job.Job{
				ID:        "test-job-id",
				Status:    job.JobStatusRunning,
				CreatedAt: time.Now(),
				StartedAt: &time.Time{},
				Progress:  0,
				Message:   "Starting job...",
			},
			expectedStatus: http.StatusAccepted,
			setupMocks: func(ingestionService interfaces.IngestionServiceInterface, jobManager *MockJobManager) {
				jobManager.On("CreateJob").Return(&job.Job{
					ID:        "test-job-id",
					Status:    job.JobStatusRunning,
					CreatedAt: time.Now(),
					StartedAt: &time.Time{},
					Progress:  0,
					Message:   "Starting job...",
				}, nil)
				jobManager.On("RunJobAsync", "test-job-id", mock.Anything).Return(nil)
			},
		},
		{
			name:           "job creation error",
			mockJob:        nil,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(ingestionService interfaces.IngestionServiceInterface, jobManager *MockJobManager) {
				jobManager.On("CreateJob").Return(nil, assert.AnError)
			},
		},
		{
			name:           "job start error",
			mockJob:        nil,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(ingestionService interfaces.IngestionServiceInterface, jobManager *MockJobManager) {
				jobManager.On("CreateJob").Return(&job.Job{
					ID:        "test-job-id",
					Status:    job.JobStatusPending,
					CreatedAt: time.Now(),
				}, nil)
				jobManager.On("RunJobAsync", "test-job-id", mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			mockIngestionService := &MockIngestionService{}
			mockJobManager := &MockJobManager{}
			tt.setupMocks(mockIngestionService, mockJobManager)

			handler := &StocksIngestionHandler{
				ingestionService: mockIngestionService,
				jobManager:       mockJobManager,
				logger:           logrus.New(),
			}

			// Create request
			req, _ := http.NewRequest("POST", "/api/v1/admin/ingest/stocks", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.TriggerIngestion(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mocks
			mockIngestionService.AssertExpectations(t)
			mockJobManager.AssertExpectations(t)
		})
	}
}

func TestStocksIngestionHandler_GetJobStatus(t *testing.T) {
	tests := []struct {
		name           string
		jobID          string
		mockJob        *job.Job
		expectedStatus int
		setupMocks     func(*MockJobManager)
	}{
		{
			name:  "successful get status",
			jobID: "test-job-id",
			mockJob: &job.Job{
				ID:        "test-job-id",
				Status:    job.JobStatusCompleted,
				CreatedAt: time.Now(),
				StartedAt: &time.Time{},
				Progress:  100,
				Message:   "Job completed successfully",
			},
			expectedStatus: http.StatusOK,
			setupMocks: func(jobManager *MockJobManager) {
				jobManager.On("GetJob", "test-job-id").Return(&job.Job{
					ID:        "test-job-id",
					Status:    job.JobStatusCompleted,
					CreatedAt: time.Now(),
					StartedAt: &time.Time{},
					Progress:  100,
					Message:   "Job completed successfully",
				}, true)
			},
		},
		{
			name:           "job not found",
			jobID:          "invalid-job-id",
			mockJob:        nil,
			expectedStatus: http.StatusNotFound,
			setupMocks: func(jobManager *MockJobManager) {
				jobManager.On("GetJob", "invalid-job-id").Return(nil, false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockJobManager := &MockJobManager{}
			tt.setupMocks(mockJobManager)

			handler := &StocksIngestionHandler{
				ingestionService: &MockIngestionService{},
				jobManager:       mockJobManager,
				logger:           logrus.New(),
			}

			req, _ := http.NewRequest("GET", "/api/v1/admin/jobs/"+tt.jobID, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "jobId", Value: tt.jobID}}

			handler.GetJobStatus(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockJobManager.AssertExpectations(t)
		})
	}
}
