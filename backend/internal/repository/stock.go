package repository

import (
	"context"
	"database/sql"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

func SaveReportsDB(db *sql.DB, reports []model.StockReport) error{
	query := `INSERT INTO STOCK_REPORTS(
		ID, TICKER, COMPANY, BROKERAGE, ACTION, RATING_FROM,
		RATING_TO, TARGET_FROM, TARGET_TO, TIME
	) VALUES(
	 	gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9
	)`

	for _, report := range reports {
		_, err := db.ExecContext(context.Background(), query, 
		 	report.Ticker, report.Company, report.Brokerage,
		 	report.Action, report.RatingFrom, report.RatingTo,
			report.TargetFrom, report.TargetTo, report.Time,
		)
		if err != nil{
			return err
		}
	}
	return nil

}

