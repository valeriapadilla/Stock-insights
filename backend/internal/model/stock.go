package model
import "time"

type StockReport struct {
	Ticker      string    `json:"ticker"`
	Company     string    `json:"company"`
	Brokerage   string    `json:"brokerage"`
	Action      string    `json:"action"`
	RatingFrom  string    `json:"rating_from"`
	RatingTo    string    `json:"rating_to"`
	TargetFrom  string   `json:"target_from"`
	TargetTo    string   `json:"target_to"`
	Time        time.Time `json:"time"`
}

type StockApiResponse struct {
	Items []StockReport `json:"items"`
	NextPage string 	`json:"next_page"`
}
