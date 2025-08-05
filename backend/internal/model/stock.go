package model

import (
	"strings"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/utils"
)

type Stock struct {
	Ticker        string    `json:"ticker" db:"ticker"`
	Company       string    `json:"company" db:"company"`
	TargetFrom    string    `json:"target_from" db:"target_from"`
	TargetTo      string    `json:"target_to" db:"target_to"`
	RatingFrom    string    `json:"rating_from" db:"rating_from"`
	RatingTo      string    `json:"rating_to" db:"rating_to"`
	Action        string    `json:"action" db:"action"`
	Brokerage     string    `json:"brokerage" db:"brokerage"`
	Time          time.Time `json:"time" db:"time"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	ChangePercent string    `json:"change_percent,omitempty"` // Calculado dinámicamente
}

func (s *Stock) GetRating() string {
	if s.RatingTo != "" {
		return strings.ToLower(s.RatingTo)
	}
	return strings.ToLower(s.RatingFrom)
}

func (s *Stock) GetTargetFromPrice() float64 {
	return utils.ParsePrice(s.TargetFrom)
}

func (s *Stock) GetTargetToPrice() float64 {
	return utils.ParsePrice(s.TargetTo)
}

func (s *Stock) GetChangePercentage() float64 {
	fromPrice := s.GetTargetFromPrice()
	toPrice := s.GetTargetToPrice()
	return utils.CalculateChangePercentage(fromPrice, toPrice)
}

type ExternalAPIResponse struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page,omitempty"`
}
