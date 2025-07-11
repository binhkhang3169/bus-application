package handler

import (
	"go-bigquery-dashboard/internal/repository"
	"go-bigquery-dashboard/internal/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	repo *repository.DashboardRepository
}

func NewDashboardHandler(repo *repository.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

func (h *DashboardHandler) GetKPIs(c *gin.Context) {
	startDate, endDate, ok := util.ParseTimeRange(c)
	if !ok {
		return
	}
	data, err := h.repo.GetKPIs(startDate, endDate)
	if err != nil {
		log.Printf("Error getting KPIs data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve KPIs data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *DashboardHandler) GetRevenueOverTime(c *gin.Context) {
	startDate, endDate, ok := util.ParseTimeRange(c)
	if !ok {
		return
	}
	groupBy := c.DefaultQuery("group_by", "day")
	switch groupBy {
	case "day":
		groupBy = "DAY"
	case "week":
		groupBy = "WEEK"
	case "month":
		groupBy = "MONTH"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group_by value. Use 'day', 'week', or 'month'."})
		return
	}
	data, err := h.repo.GetRevenueOverTime(startDate, endDate, groupBy)
	if err != nil {
		log.Printf("Error getting revenue over time data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve revenue data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *DashboardHandler) GetTicketDistribution(c *gin.Context) {
	startDate, endDate, ok := util.ParseTimeRange(c)
	if !ok {
		return
	}
	data, err := h.repo.GetTicketDistribution(startDate, endDate)
	if err != nil {
		log.Printf("Error getting ticket distribution data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ticket distribution data"})
		return
	}
	c.JSON(http.StatusOK, data)
}
