package service

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/valeriapadilla/stock-insights/internal/model"
)


func FetchAllStockReports() ([]model.StockReport, error){
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	baseUrl := os.Getenv("BASEURL")
    token := os.Getenv("TOKEN")
	
	client := http.Client{}
	var stockReports []model.StockReport
	nextPage := ""

	for{
		url:= baseUrl
		if nextPage != ""{
			url += "?next_page=" + nextPage
		}

		req, err := http.NewRequest("GET",url, nil)
		if err != nil{
			return nil, err
		}

		req.Header.Set("Authorization",token)
		req.Header.Set("Content-Type","application/json")

		resp, err := client.Do(req)
		if err != nil{
			return nil, err
		}
		defer resp.Body.Close()

		var data model.StockApiResponse

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil{
			return nil, err
		}
		stockReports = append(stockReports, data.Items...)

		if(data.NextPage == ""){
			break
		}
		nextPage = data.NextPage
	}

	fmt.Printf("Fetched %d stock reports\n", len(stockReports))
	return stockReports, nil
}



