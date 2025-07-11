package controllers

import (
	"net/http"
	"ticket-service/internal/services"

	"github.com/gin-gonic/gin"
)

type ManagerTicketController struct {
	managerTicketService services.IManagerTicketService
}

func NewManagerTicketController(managerTicketService services.IManagerTicketService) *ManagerTicketController {
	return &ManagerTicketController{
		managerTicketService: managerTicketService,
	}
}

type TripRequest struct {
	TripID string `json:"trip_id" binding:"required"`
}

type UpdateTripRequest struct {
	TicketID   string `json:"ticket_id" binding:"required"`
	StatusCode string `json:"status_code" binding:"required"`
}

// Create seat for tripID
func (m *ManagerTicketController) CreateManagerTicketHandler(c *gin.Context) {
	var req TripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "trip_id is required",
			"data":    nil,
		})
		return
	}

	tickets, err := m.managerTicketService.CreateManagerTicketsForTrip(c.Request.Context(), req.TripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "Manager tickets created",
		"data":    tickets,
	})
}

// Update status payment seats
func (m *ManagerTicketController) UpdateManagerTicketHandler(c *gin.Context) {
	var req UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "lack of data", // Corrected typo from "lake off data"
			"data":    nil,
		})
		return
	}

	err := m.managerTicketService.UpdateStatusByTicketID(c.Request.Context(), req.TicketID, req.StatusCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "failed to update status",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Status updated successfully",
		"data":    nil,
	})
}
