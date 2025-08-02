package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JobManager struct {
	jobs    map[string]*Job
	mutex   sync.RWMutex
	logger  *logrus.Logger
	workers chan struct{}
}

func NewJobManager(maxWorkers int, logger *logrus.Logger) *JobManager {
	return &JobManager{
		jobs:    make(map[string]*Job),
		logger:  logger,
		workers: make(chan struct{}, maxWorkers),
	}
}

func (jm *JobManager) CreateJob() (*Job, error) {
	job := &Job{
		ID:        uuid.New().String(),
		Status:    JobStatusPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	jm.jobs[job.ID] = job

	jm.logger.WithField("job_id", job.ID).Info("Created new job")
	return job, nil
}

func (jm *JobManager) GetJob(jobID string) (*Job, bool) {
	jm.mutex.RLock()
	defer jm.mutex.RUnlock()

	job, exists := jm.jobs[jobID]
	return job, exists
}

func (jm *JobManager) UpdateJob(jobID string, status JobStatus, progress int, message string) {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	if job, exists := jm.jobs[jobID]; exists {
		job.Status = status
		job.Progress = progress
		job.Message = message

		if status == JobStatusRunning && job.StartedAt == nil {
			now := time.Now()
			job.StartedAt = &now
		} else if (status == JobStatusCompleted || status == JobStatusFailed) && job.EndedAt == nil {
			now := time.Now()
			job.EndedAt = &now
		}

		jm.logger.WithFields(logrus.Fields{
			"job_id":   jobID,
			"status":   status,
			"progress": progress,
			"message":  message,
		}).Info("Job updated")
	}
}

func (jm *JobManager) SetJobError(jobID string, err error) {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	if job, exists := jm.jobs[jobID]; exists {
		job.Status = JobStatusFailed
		job.Error = err.Error()
		now := time.Now()
		job.EndedAt = &now

		jm.logger.WithFields(logrus.Fields{
			"job_id": jobID,
			"error":  err.Error(),
		}).Error("Job failed")
	}
}

func (jm *JobManager) RunJobAsync(jobID string, workFunc func(context.Context) error) error {
	jm.mutex.RLock()
	_, exists := jm.jobs[jobID]
	jm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	select {
	case jm.workers <- struct{}{}:
	default:
		return fmt.Errorf("no available workers")
	}

	go func() {
		defer func() {
			<-jm.workers
		}()

		ctx := context.Background()

		jm.UpdateJob(jobID, JobStatusRunning, 0, "Starting job...")

		if err := workFunc(ctx); err != nil {
			jm.SetJobError(jobID, err)
		} else {
			jm.UpdateJob(jobID, JobStatusCompleted, 100, "Job completed successfully")
		}
	}()

	return nil
}

func (jm *JobManager) CleanupOldJobs(maxAge time.Duration) {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	cutoff := time.Now().Add(-maxAge)

	for jobID, job := range jm.jobs {
		if job.CreatedAt.Before(cutoff) {
			delete(jm.jobs, jobID)
			jm.logger.WithField("job_id", jobID).Debug("Cleaned up old job")
		}
	}
}
