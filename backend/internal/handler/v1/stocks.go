package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/service/interfaces"
)

type StocksHandler struct {
	stockService interfaces.StockServiceInterface
	logger       *logrus.Logger
}

func NewStocksHandler(stockService interfaces.StockServiceInterface, logger *logrus.Logger) *StocksHandler {
	return &StocksHandler{
		stockService: stockService,
		logger:       logger,
	}
}

func (h *StocksHandler) ListStocks(c *gin.Context) {
	limit, offset, sort, order := parsePaginationParams(c)

	stocks, total, err := h.stockService.ListStocks(limit, offset, sort, order)
	if err != nil {
		handleError(c, err, "retrieve stocks", h.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stocks": stocks,
		"pagination": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_next": offset+limit < total,
		},
	})
}

func (h *StocksHandler) GetStock(c *gin.Context) {
	ticket := c.Param("ticket")
	if ticket == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Ticket parameter is required",
		})
		return
	}

	stock, err := h.stockService.GetStock(ticket)
	if err != nil {
		handleError(c, err, "retrieve stock", h.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stock": stock,
	})
}

func (h *StocksHandler) SearchStocks(c *gin.Context) {
	limit, offset, _, _ := parsePaginationParams(c)
	minPrice, maxPrice := parsePriceParams(c)

	ticket := c.Query("ticket")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	stocks, total, err := h.stockService.SearchStocks(interfaces.StockSearchParams{
		Ticket:   ticket,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Limit:    limit,
		Offset:   offset,
	})

	if err != nil {
		handleError(c, err, "search stocks", h.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stocks": stocks,
		"pagination": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_next": offset+limit < total,
		},
		"filters_applied": gin.H{
			"ticket":    ticket,
			"date_from": dateFrom,
			"date_to":   dateTo,
			"min_price": minPrice,
			"max_price": maxPrice,
		},
	})
}
