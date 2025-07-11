package handler

import (
	"go-bigquery-dashboard/internal/repository"
	"go-bigquery-dashboard/internal/util"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TripSearchHandler struct {
	repo *repository.TripSearchRepository
}

func NewTripSearchHandler(repo *repository.TripSearchRepository) *TripSearchHandler {
	return &TripSearchHandler{repo: repo}
}

func (h *TripSearchHandler) GetTopRoutes(c *gin.Context) {
	startDate, endDate, ok := util.ParseTimeRange(c)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := h.repo.GetTopRoutes(startDate, endDate, limit)
	if err != nil {
		log.Printf("Error getting top routes data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top routes data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *TripSearchHandler) GetTopProvinces(c *gin.Context) {
	startDate, endDate, ok := util.ParseTimeRange(c)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := h.repo.GetTopProvinces(startDate, endDate, limit)
	if err != nil {
		log.Printf("Error getting top provinces data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top provinces data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *TripSearchHandler) GetSearchesByHour(c *gin.Context) {
	startDate, endDate, ok := util.ParseTimeRange(c)
	if !ok {
		return
	}
	data, err := h.repo.GetSearchesByHourOfDay(startDate, endDate)
	if err != nil {
		log.Printf("Error getting searches by hour data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve searches by hour data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *TripSearchHandler) GetSearchesOverTime(c *gin.Context) {
	startDate, endDate, ok := util.ParseTimeRange(c)
	if !ok {
		return
	}
	data, err := h.repo.GetSearchesOverTime(startDate, endDate)
	if err != nil {
		log.Printf("Error getting searches over time data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve searches over time data"})
		return
	}
	c.JSON(http.StatusOK, data)
}
