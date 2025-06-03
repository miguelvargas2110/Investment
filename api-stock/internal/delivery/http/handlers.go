package http

import (
	"api-stock/internal/domain"
	"api-stock/pkg/errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService          domain.StockService
	recommendationService domain.RecommendationService
}

func NewStockHandler(
	stockService domain.StockService,
	recommendationService domain.RecommendationService,
) *StockHandler {
	return &StockHandler{
		stockService:          stockService,
		recommendationService: recommendationService,
	}
}

// GetRecommendations godoc
// @Summary Get stock recommendations
// @Description Get paginated list of stock recommendations, filterable by ticker
// @Tags recommendations
// @Accept json
// @Produce json
// @Param ticker query string false "Stock ticker to filter by"
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(50) minimum(1) maximum(100)
// @Success 200 {object} map[string]interface{} "Returns recommendations and pagination info"
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /http/v1/recommendations [get]
func (h *StockHandler) GetRecommendations(c *gin.Context) {
	ticker := c.Query("ticker")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	recommendations, total, err := h.stockService.GetRecommendations(c.Request.Context(), ticker, page, limit)
	if err != nil {
		c.Error(errors.NewAppError(http.StatusInternalServerError, "Failed to get recommendations", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": recommendations,
		"pagination": gin.H{
			"page":      page,
			"limit":     limit,
			"total":     total,
			"last_page": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// GetAvailableTickers godoc
// @Summary Get available stock tickers
// @Description Get list of all available stock tickers with recommendations
// @Tags recommendations
// @Accept json
// @Produce json
// @Success 200 {array} string "List of tickers"
// @Failure 500 {object} errors.AppError "Internal server error"
// @Router /http/v1/recommendations/tickers [get]
func (h *StockHandler) GetAvailableTickers(c *gin.Context) {
	tickers, err := h.stockService.GetAvailableTickers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickers)
}

// GetBestRecommendations godoc
// @Summary Get best stock recommendations
// @Description Get top stock recommendations based on scoring model
// @Tags recommendations
// @Accept json
// @Produce json
// @Param limit query int false "Number of recommendations to return" default(5) minimum(1) maximum(20)
// @Success 200 {object} gin.H "Returns best recommendations and generation timestamp"
// @Failure 500 {object} errors.AppError "Internal server error"
// @Router /http/v1/recommendations/best [get]
func (h *StockHandler) GetBestRecommendations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	best, err := h.recommendationService.GetBestStocks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"best_recommendations": best,
		"generated_at":         time.Now().Format(time.RFC3339),
	})
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} gin.H "Service status and version"
// @Failure 503 {object} errors.AppError "Service unavailable"
// @Router /http/v1/health [get]
func (h *StockHandler) HealthCheck(c *gin.Context) {
	if err := h.stockService.HealthCheck(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
	})
}
