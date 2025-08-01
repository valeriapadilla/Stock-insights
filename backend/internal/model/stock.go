package model

import "time"

type Stock struct {
	Ticker     string    `json:"ticker" db:"ticker"`
	Company    string    `json:"company" db:"company"`
	TargetFrom string    `json:"target_from" db:"target_from"`
	TargetTo   string    `json:"target_to" db:"target_to"`
	RatingFrom string    `json:"rating_from" db:"rating_from"`
	RatingTo   string    `json:"rating_to" db:"rating_to"`
	Action     string    `json:"action" db:"action"`
	Brokerage  string    `json:"brokerage" db:"brokerage"`
	Time       time.Time `json:"time" db:"time"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ExternalAPIResponse struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page,omitempty"`
}
