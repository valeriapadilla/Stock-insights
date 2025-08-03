package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/service/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type RecommendationsHandler struct {
	recommendationService interfaces.RecommendationServiceInterface
	logger                *logrus.Logger
}

func NewRecommendationsHandler(
	recommendationService interfaces.RecommendationServiceInterface,
	logger *logrus.Logger,
) *RecommendationsHandler {
	return &RecommendationsHandler{
		recommendationService: recommendationService,
		logger:                logger,
	}
}

func (h *RecommendationsHandler) GetRecommendations(c *gin.Context) {
	limit, err := h.parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	recommendations, err := h.recommendationService.GetLatestRecommendations(limit)
	if err != nil {
		handleError(c, err, "retrieve recommendations", h.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"total":           len(recommendations),
		"limit":           limit,
	})
}

func (h *RecommendationsHandler) CalculateRecommendations(c *gin.Context) {
	params := h.parseRecommendationParams(c)

	recommendations, err := h.recommendationService.CalculateRecommendations(params)
	if err != nil {
		handleError(c, err, "calculate recommendations", h.logger)
		return
	}

	if err := h.recommendationService.SaveRecommendations(recommendations); err != nil {
		handleError(c, err, "save recommendations", h.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Recommendations calculated and saved successfully",
		"recommendations": recommendations,
		"total":           len(recommendations),
		"run_at":          recommendations[0].RunAt,
	})
}

func (h *RecommendationsHandler) parseLimitParam(c *gin.Context) (int, error) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, err
	}
	return limit, nil
}

func (h *RecommendationsHandler) parseRecommendationParams(c *gin.Context) validator.RecommendationParams {
	daysBackStr := c.DefaultQuery("days_back", "7")
	maxResultsStr := c.DefaultQuery("max_results", "30")
	minScoreStr := c.DefaultQuery("min_score", "80")

	daysBack, _ := strconv.Atoi(daysBackStr)
	maxResults, _ := strconv.Atoi(maxResultsStr)
	minScore, _ := strconv.Atoi(minScoreStr)

	return validator.RecommendationParams{
		DaysBack:   daysBack,
		MaxResults: maxResults,
		MinScore:   minScore,
	}
}
