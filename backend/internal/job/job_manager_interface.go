package job

import (
	"context"
	"time"
)

type JobManagerInterface interface {
	CreateJob() (*Job, error)
	GetJob(jobID string) (*Job, bool)
	UpdateJob(jobID string, status JobStatus, progress int, message string)
	SetJobError(jobID string, err error)
	RunJobAsync(jobID string, workFunc func(context.Context) error) error
	CleanupOldJobs(maxAge time.Duration)
}
