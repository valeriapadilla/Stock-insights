package implementations

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationWorkerImpl_NewRecommendationWorker(t *testing.T) {
	logger := logrus.New()

	worker := &RecommendationWorkerImpl{
		logger: logger,
	}

	assert.NotNil(t, worker)
	assert.NotNil(t, worker.logger)
	assert.False(t, worker.isRunning)
}

func TestRecommendationWorkerImpl_StructInitialization(t *testing.T) {
	worker := &RecommendationWorkerImpl{
		logger:    logrus.New(),
		isRunning: false,
	}

	assert.NotNil(t, worker)
	assert.NotNil(t, worker.logger)
	assert.False(t, worker.isRunning)
}

func TestRecommendationWorkerImpl_IsRunning(t *testing.T) {
	worker := &RecommendationWorkerImpl{
		logger:    logrus.New(),
		isRunning: false,
	}

	assert.False(t, worker.IsRunning())

	worker.isRunning = true
	assert.True(t, worker.IsRunning())
}

func TestRecommendationWorkerImpl_GetLastRunTime(t *testing.T) {
	worker := &RecommendationWorkerImpl{
		logger: logrus.New(),
	}

	lastRun, err := worker.GetLastRunTime()
	assert.NoError(t, err)
	assert.Nil(t, lastRun)

	now := time.Now()
	worker.lastRunTime = &now
	lastRun, err = worker.GetLastRunTime()
	assert.NoError(t, err)
	assert.Equal(t, now, *lastRun)
}
