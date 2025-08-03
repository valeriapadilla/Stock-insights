package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/errors"
	"github.com/valeriapadilla/stock-insights/internal/model"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/service/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type ScoringConfig struct {
	// Action Scores (0-40 points)
	TargetRaisedScore     int
	UpgradedScore         int
	InitiatedScore        int
	TargetMaintainedScore int

	// Rating Scores (0-25 points)
	BuyScore           int
	OverweightScore    int
	SectorPerformScore int
	EqualWeightScore   int
	NeutralScore       int

	// Target Change Scores (0-20 points)
	HighTargetChangeScore   int
	MediumTargetChangeScore int
	LowTargetChangeScore    int
	MinTargetChangeScore    int

	// Freshness Scores (0-15 points)
	TodayScore     int
	YesterdayScore int
	ThreeDaysScore int
	WeekScore      int

	// Thresholds
	HighTargetChangePercent   float64
	MediumTargetChangePercent float64
	LowTargetChangePercent    float64
	MinTargetChangePercent    float64
}

func DefaultScoringConfig() *ScoringConfig {
	return &ScoringConfig{
		TargetRaisedScore:     40,
		UpgradedScore:         35,
		InitiatedScore:        30,
		TargetMaintainedScore: 20,

		BuyScore:           25,
		OverweightScore:    20,
		SectorPerformScore: 15,
		EqualWeightScore:   10,
		NeutralScore:       5,

		HighTargetChangeScore:   20,
		MediumTargetChangeScore: 15,
		LowTargetChangeScore:    10,
		MinTargetChangeScore:    5,

		TodayScore:     15,
		YesterdayScore: 12,
		ThreeDaysScore: 10,
		WeekScore:      5,

		HighTargetChangePercent:   50.0,
		MediumTargetChangePercent: 25.0,
		LowTargetChangePercent:    10.0,
		MinTargetChangePercent:    5.0,
	}
}

type RecommendationService struct {
	stockRepo          repoInterfaces.StockRepository
	recommendationRepo repoInterfaces.RecommendationRepository
	recommendationCmd  repoInterfaces.RecommendationCommand
	logger             *logrus.Logger
	scoringConfig      *ScoringConfig
	validator          *validator.RecommendationValidator
}

var _ interfaces.RecommendationServiceInterface = (*RecommendationService)(nil)

type StockScore struct {
	Stock       *model.Stock
	Score       int
	Explanation string
}

func NewRecommendationService(
	stockRepo repoInterfaces.StockRepository,
	recommendationRepo repoInterfaces.RecommendationRepository,
	recommendationCmd repoInterfaces.RecommendationCommand,
	logger *logrus.Logger,
) *RecommendationService {
	return &RecommendationService{
		stockRepo:          stockRepo,
		recommendationRepo: recommendationRepo,
		recommendationCmd:  recommendationCmd,
		logger:             logger,
		scoringConfig:      DefaultScoringConfig(),
		validator:          validator.NewRecommendationValidator(),
	}
}

func (s *RecommendationService) CalculateRecommendations(params validator.RecommendationParams) ([]*model.Recommendation, error) {
	validatedParams := s.validator.ValidateRecommendationParams(params)

	if err := s.recommendationCmd.DeleteAllRecommendations(); err != nil {
		s.logger.WithError(err).Error("Failed to delete existing recommendations")
		return nil, errors.NewDatabaseError("failed to delete existing recommendations", err)
	}

	s.logger.Info("Deleted all existing recommendations, calculating new ones...")

	stocks, err := s.getStocksForRecommendations(validatedParams.DaysBack)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get stocks for recommendations")
		return nil, errors.NewDatabaseError("failed to get stocks for recommendations", err)
	}

	stockScores := s.calculateStockScores(stocks)
	filteredScores := s.filterAndSortScores(stockScores, validatedParams.MinScore, validatedParams.MaxResults)
	recommendations := s.convertToRecommendations(filteredScores)

	s.logRecommendationCalculation(stocks, stockScores, filteredScores, recommendations, validatedParams)

	return recommendations, nil
}

func (s *RecommendationService) GetLatestRecommendations(limit int) ([]*model.Recommendation, error) {
	validatedLimit := s.validator.ValidateLimit(limit, 10)

	recommendations, err := s.recommendationRepo.GetLatest(validatedLimit)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get latest recommendations")
		return nil, errors.NewDatabaseError("failed to get latest recommendations", err)
	}

	return recommendations, nil
}

func (s *RecommendationService) SaveRecommendations(recommendations []*model.Recommendation) error {
	runAt := time.Now()

	for i, rec := range recommendations {
		if rec.ID == "" {
			rec.ID = uuid.New().String()
		}
		rec.RunAt = runAt
		rec.Rank = i + 1

		if err := s.recommendationRepo.CreateRecommendation(rec); err != nil {
			s.logger.WithError(err).WithField("ticker", rec.Ticker).Error("Failed to save recommendation")
			return errors.NewDatabaseError("failed to save recommendation", err)
		}
	}

	s.logger.WithField("count", len(recommendations)).Info("Recommendations saved successfully")
	return nil
}

func (s *RecommendationService) getStocksForRecommendations(daysBack int) ([]*model.Stock, error) {
	cutoffDate := time.Now().AddDate(0, 0, -daysBack)

	countParams := repoInterfaces.GetStocksParams{
		Search: &repoInterfaces.StockSearchFilters{
			DateFrom: &cutoffDate,
		},
	}

	totalCount, err := s.stockRepo.GetStocksCount(countParams)
	if err != nil {
		return nil, err
	}
	limit := totalCount + 100

	params := repoInterfaces.GetStocksParams{
		Limit: limit,
		Search: &repoInterfaces.StockSearchFilters{
			DateFrom: &cutoffDate,
		},
	}

	stocks, err := s.stockRepo.GetStocks(params)
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"days_back":    daysBack,
		"total_stocks": totalCount,
		"limit_used":   limit,
		"stocks_found": len(stocks),
	}).Info("Retrieved stocks for recommendations")

	return s.filterPositiveStocks(stocks), nil
}

func (s *RecommendationService) filterPositiveStocks(stocks []*model.Stock) []*model.Stock {
	var positiveStocks []*model.Stock
	for _, stock := range stocks {
		if s.isPositiveAction(stock.Action) && s.isPositiveRating(stock.RatingTo) {
			positiveStocks = append(positiveStocks, stock)
		}
	}
	return positiveStocks
}

func (s *RecommendationService) calculateStockScores(stocks []*model.Stock) []StockScore {
	var stockScores []StockScore

	for _, stock := range stocks {
		score := s.calculateStockScore(stock)
		explanation := s.generateExplanation(stock, score)

		stockScores = append(stockScores, StockScore{
			Stock:       stock,
			Score:       score,
			Explanation: explanation,
		})
	}

	return stockScores
}

func (s *RecommendationService) calculateStockScore(stock *model.Stock) int {
	score := 0

	score += s.getActionScore(stock.Action)
	score += s.getRatingScore(stock.RatingTo)
	score += s.getTargetChangeScore(stock.TargetFrom, stock.TargetTo)
	score += s.getFreshnessScore(stock.Time)

	return score
}

func (s *RecommendationService) getActionScore(action string) int {
	switch strings.ToLower(action) {
	case "target raised by":
		return s.scoringConfig.TargetRaisedScore
	case "upgraded by":
		return s.scoringConfig.UpgradedScore
	case "initiated by":
		return s.scoringConfig.InitiatedScore
	case "target maintained by":
		return s.scoringConfig.TargetMaintainedScore
	default:
		return 0
	}
}

func (s *RecommendationService) getRatingScore(rating string) int {
	switch strings.ToLower(rating) {
	case "buy":
		return s.scoringConfig.BuyScore
	case "overweight":
		return s.scoringConfig.OverweightScore
	case "sector perform":
		return s.scoringConfig.SectorPerformScore
	case "equal weight":
		return s.scoringConfig.EqualWeightScore
	case "neutral":
		return s.scoringConfig.NeutralScore
	default:
		return 0
	}
}

func (s *RecommendationService) getTargetChangeScore(targetFrom, targetTo string) int {
	fromPrice, err1 := s.extractPrice(targetFrom)
	toPrice, err2 := s.extractPrice(targetTo)

	if err1 != nil || err2 != nil {
		return 0
	}

	if fromPrice <= 0 || toPrice <= 0 {
		return 0
	}

	changePercent := ((toPrice - fromPrice) / fromPrice) * 100

	switch {
	case changePercent > s.scoringConfig.HighTargetChangePercent:
		return s.scoringConfig.HighTargetChangeScore
	case changePercent > s.scoringConfig.MediumTargetChangePercent:
		return s.scoringConfig.MediumTargetChangeScore
	case changePercent > s.scoringConfig.LowTargetChangePercent:
		return s.scoringConfig.LowTargetChangeScore
	case changePercent > s.scoringConfig.MinTargetChangePercent:
		return s.scoringConfig.MinTargetChangeScore
	case changePercent > 0:
		return 2
	default:
		return 0
	}
}

// getFreshnessScore returns score based on how recent the data is (0-15 points)
func (s *RecommendationService) getFreshnessScore(stockTime time.Time) int {
	daysSince := int(time.Since(stockTime).Hours() / 24)

	switch {
	case daysSince == 0:
		return s.scoringConfig.TodayScore
	case daysSince == 1:
		return s.scoringConfig.YesterdayScore
	case daysSince <= 3:
		return s.scoringConfig.ThreeDaysScore
	case daysSince <= 7:
		return s.scoringConfig.WeekScore
	default:
		return 0
	}
}

func (s *RecommendationService) extractPrice(priceStr string) (float64, error) {
	cleanPrice := strings.TrimSpace(strings.Replace(priceStr, "$", "", -1))
	return strconv.ParseFloat(cleanPrice, 64)
}

func (s *RecommendationService) isPositiveAction(action string) bool {
	positiveActions := []string{
		"target raised by",
		"upgraded by",
		"initiated by",
		"target maintained by",
	}

	actionLower := strings.ToLower(action)
	for _, positive := range positiveActions {
		if actionLower == positive {
			return true
		}
	}
	return false
}

func (s *RecommendationService) isPositiveRating(rating string) bool {
	positiveRatings := []string{
		"buy",
		"overweight",
		"sector perform",
	}

	ratingLower := strings.ToLower(rating)
	for _, positive := range positiveRatings {
		if ratingLower == positive {
			return true
		}
	}
	return false
}

func (s *RecommendationService) filterAndSortScores(scores []StockScore, minScore, maxResults int) []StockScore {
	var filtered []StockScore

	for _, score := range scores {
		if score.Score >= minScore {
			filtered = append(filtered, score)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Score > filtered[j].Score
	})

	if len(filtered) > maxResults {
		filtered = filtered[:maxResults]
	}

	return filtered
}

func (s *RecommendationService) convertToRecommendations(scores []StockScore) []*model.Recommendation {
	var recommendations []*model.Recommendation

	for i, score := range scores {
		recommendation := &model.Recommendation{
			ID:          uuid.New().String(),
			Ticker:      score.Stock.Ticker,
			Score:       float64(score.Score),
			Explanation: score.Explanation,
			RunAt:       time.Now(),
			Rank:        i + 1,
		}
		recommendations = append(recommendations, recommendation)
	}

	return recommendations
}

func (s *RecommendationService) generateExplanation(stock *model.Stock, score int) string {
	var reasons []string

	reasons = append(reasons, s.getActionExplanation(stock.Action))
	reasons = append(reasons, s.getRatingExplanation(stock.RatingTo))
	reasons = append(reasons, s.getTargetChangeExplanation(stock.TargetFrom, stock.TargetTo))
	reasons = append(reasons, s.getBrokerageExplanation(stock.Brokerage))
	reasons = append(reasons, fmt.Sprintf("Score: %d/100", score))

	return strings.Join(reasons, ", ")
}

func (s *RecommendationService) getActionExplanation(action string) string {
	switch strings.ToLower(action) {
	case "target raised by":
		return "Target price raised"
	case "upgraded by":
		return "Rating upgraded"
	case "initiated by":
		return "New coverage initiated"
	case "target maintained by":
		return "Target price maintained"
	default:
		return ""
	}
}

func (s *RecommendationService) getRatingExplanation(rating string) string {
	switch strings.ToLower(rating) {
	case "buy":
		return "Buy rating"
	case "overweight":
		return "Overweight rating"
	default:
		return ""
	}
}

func (s *RecommendationService) getTargetChangeExplanation(targetFrom, targetTo string) string {
	if fromPrice, err1 := s.extractPrice(targetFrom); err1 == nil {
		if toPrice, err2 := s.extractPrice(targetTo); err2 == nil {
			if toPrice > fromPrice {
				changePercent := ((toPrice - fromPrice) / fromPrice) * 100
				return fmt.Sprintf("Target raised by %.1f%%", changePercent)
			}
		}
	}
	return ""
}

func (s *RecommendationService) getBrokerageExplanation(brokerage string) string {
	if brokerage != "" {
		return fmt.Sprintf("by %s", brokerage)
	}
	return ""
}

func (s *RecommendationService) logRecommendationCalculation(
	stocks []*model.Stock,
	stockScores []StockScore,
	filteredScores []StockScore,
	recommendations []*model.Recommendation,
	params validator.RecommendationParams,
) {
	s.logger.WithFields(logrus.Fields{
		"total_stocks":    len(stocks),
		"scored_stocks":   len(stockScores),
		"filtered_stocks": len(filteredScores),
		"recommendations": len(recommendations),
		"min_score":       params.MinScore,
		"max_results":     params.MaxResults,
	}).Info("Recommendations calculated successfully")
}
