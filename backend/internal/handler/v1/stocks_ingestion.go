package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/job"
	"github.com/valeriapadilla/stock-insights/internal/service/interfaces"
)

type StocksIngestionHandler struct {
	ingestionService interfaces.IngestionServiceInterface
	jobManager       job.JobManagerInterface
	logger           *logrus.Logger
}

func NewStocksIngestionHandler(ingestionService interfaces.IngestionServiceInterface, jobManager job.JobManagerInterface, logger *logrus.Logger) *StocksIngestionHandler {
	return &StocksIngestionHandler{
		ingestionService: ingestionService,
		jobManager:       jobManager,
		logger:           logger,
	}
}

func (h *StocksIngestionHandler) TriggerIngestion(c *gin.Context) {
	userID := c.GetString("user_id")
	h.logger.WithFields(logrus.Fields{
		"endpoint": "/api/v1/admin/ingest/stocks",
		"method":   "POST",
		"user_id":  userID,
	}).Info("Manual stocks ingestion triggered")

	job, err := h.jobManager.CreateJob()
	if err != nil {
		h.logger.WithError(err).Error("Failed to create ingestion job")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create job",
			"error":   err.Error(),
		})
		return
	}

	if err := h.jobManager.RunJobAsync(job.ID, func(ctx context.Context) error {
		return h.ingestionService.TriggerIngestionAsync(ctx)
	}); err != nil {
		h.logger.WithError(err).Error("Failed to start ingestion job")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to start job",
			"error":   err.Error(),
		})
		return
	}

	h.logger.WithField("job_id", job.ID).Info("Ingestion job started")

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"message": "Ingestion job started",
		"job_id":  job.ID,
		"job":     job,
	})
}

func (h *StocksIngestionHandler) GetJobStatus(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Job ID is required",
		})
		return
	}

	job, exists := h.jobManager.GetJob(jobID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Job not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"job":    job,
	})
}
