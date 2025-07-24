package main

import (
	"log"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/app"
)

func main() {
	godotenv.Load()
	db, err := database.Connect()
	if err != nil {
		fmt.Printf("database connection error: %v\n", err)
	}
	defer db.Close()
	if err := app.SyncStockReports(db); err != nil {
	  fmt.Printf("sync error: %v\n", err)
	}
}
