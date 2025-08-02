package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/errors"
	workerInterfaces "github.com/valeriapadilla/stock-insights/internal/worker/interfaces"
)

type IngestionService struct {
	dataWorker workerInterfaces.DataWorker
	logger     *logrus.Logger
}

func NewIngestionService(dataWorker workerInterfaces.DataWorker, logger *logrus.Logger) *IngestionService {
	return &IngestionService{
		dataWorker: dataWorker,
		logger:     logger,
	}
}

func (s *IngestionService) TriggerIngestionAsync(ctx context.Context) error {
	s.logger.Info("Starting async ingestion process")

	if err := s.dataWorker.FetchAndProcessStocks(ctx); err != nil {
		s.logger.WithError(err).Error("Async ingestion failed")
		return errors.NewInternalError("Failed to process stocks", err)
	}

	s.logger.Info("Async ingestion completed successfully")
	return nil
}
