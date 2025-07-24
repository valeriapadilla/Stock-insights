package app

import (
	"database/sql"
	"github.com/valeriapadilla/stock-insights/internal/service"
	"github.com/valeriapadilla/stock-insights/internal/repository"
)

func SyncStockReports(db *sql.DB) error{
	reports, err := service.FetchAllStockReports()
	if err != nil{
		return err
	}

	err = repository.SaveReportsDB(db, reports)
	if err != nil{
		return err
	}

	return nil
}