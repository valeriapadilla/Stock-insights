package implementations

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDataWorkerImpl_Config(t *testing.T) {
	logger := logrus.New()
	config := DataWorkerConfig{
		ScheduleInterval: 1 * time.Hour,
		MaxRetries:       3,
		RetryDelay:       5 * time.Second,
	}

	worker := &DataWorkerImpl{
		logger: logger,
		config: config,
	}

	assert.NotNil(t, worker)
	assert.Equal(t, 1*time.Hour, worker.config.ScheduleInterval)
	assert.Equal(t, 3, worker.config.MaxRetries)
	assert.Equal(t, 5*time.Second, worker.config.RetryDelay)
}

func TestDataWorkerConfig(t *testing.T) {
	config := DataWorkerConfig{
		ScheduleInterval: 1 * time.Hour,
		MaxRetries:       3,
		RetryDelay:       5 * time.Second,
	}

	assert.Equal(t, 1*time.Hour, config.ScheduleInterval)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 5*time.Second, config.RetryDelay)
}
