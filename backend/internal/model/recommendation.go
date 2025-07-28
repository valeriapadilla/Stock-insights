package model

import "time"

type Recommendation struct {
    ID          string    `json:"id" db:"id"`
    Ticker      string    `json:"ticker" db:"ticker"`
    Score       float64   `json:"score" db:"score"`
    Explanation string    `json:"explanation" db:"explanation"`
    RunAt       time.Time `json:"run_at" db:"run_at"`
    Rank        int       `json:"rank" db:"rank"`
}