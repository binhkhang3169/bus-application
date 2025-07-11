// checkin_controller.go

package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"ticket-service/internal/services"
	"ticket-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CheckinController struct {
	checkinService services.ICheckinService
	logger         utils.Logger
}

func NewCheckinController(checkinService services.ICheckinService, logger utils.Logger) *CheckinController {
	return &CheckinController{
		checkinService: checkinService,
		logger:         logger,
	}
}

// Struct mới cho request check-in
type CheckinRequest struct {
	QRContent string `json:"qr_content" binding:"required"`
	TripID    string `json:"trip_id" binding:"required"`
}

var qrRegex = regexp.MustCompile(`^TICKET:(.*)-SEAT:(\d+)$`)

func (cc *CheckinController) CheckinHandler(c *gin.Context) {
	var req CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.logger.Error("Check-in: Bad request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request data. 'qr_content' and 'trip_id' are required.",
			"data":    nil,
		})
		return
	}

	// Sử dụng Regex để phân tích chuỗi QRContent
	matches := qrRegex.FindStringSubmatch(req.QRContent)

	// matches[0] là toàn bộ chuỗi khớp
	// matches[1] là group bắt giữ đầu tiên (ticketID)
	// matches[2] là group bắt giữ thứ hai (seatTicketID)
	if len(matches) != 3 {
		cc.logger.Error("Check-in: Invalid QR content format: '%s'", req.QRContent)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid QR content format. Expected 'TICKET:{id}-SEAT:{id}'.",
			"data":    nil,
		})
		return
	}

	ticketID := matches[1]
	// Chuyển đổi seatTicketID từ string sang int
	seatTicketID, err := strconv.Atoi(matches[2])
	if err != nil {
		// Trường hợp này hiếm khi xảy ra nếu regex đúng, nhưng vẫn nên kiểm tra
		cc.logger.Error("Check-in: Could not parse Seat ID from QR '%s': %v", req.QRContent, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid Seat ID in QR content.",
			"data":    nil,
		})
		return
	}

	userID := c.GetHeader("X-User-ID")
	userRole := c.GetHeader("X-User-Role")
	note := fmt.Sprintf("Checked-in by User %s (Role: %s)", userID, userRole)

	cc.logger.Info("Attempting Check-in for QR: '%s', TripID: %s", req.QRContent, req.TripID)
	response, err := cc.checkinService.ProcessCheckin(c.Request.Context(), ticketID, int32(seatTicketID), req.TripID, note)
	if err != nil {
		cc.logger.Error("Check-in failed for QR '%s': %v", req.QRContent, err)
		var statusCode int
		errorMessage := err.Error()

		switch {
		case strings.Contains(errorMessage, "invalid or non-existent"):
			statusCode = http.StatusNotFound
		case strings.Contains(errorMessage, "mismatch") || strings.Contains(errorMessage, "not for this trip"):
			statusCode = http.StatusConflict
		case strings.Contains(errorMessage, "checkable state") || strings.Contains(errorMessage, "already checked in") || strings.Contains(errorMessage, "cancelled"):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
			errorMessage = "Failed to process check-in"
		}

		c.JSON(statusCode, gin.H{"code": statusCode, "message": errorMessage, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Check-in successful",
		"data":    response,
	})
}

// Handler mới để lấy tất cả check-in của một chuyến đi
func (cc *CheckinController) GetTripCheckinsHandler(c *gin.Context) {
	tripID := c.Param("tripID")
	if tripID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Trip ID is required in the URL path.",
			"data":    nil,
		})
		return
	}

	cc.logger.Info("Fetching check-ins for TripID: %s", tripID)
	response, err := cc.checkinService.GetCheckinsForTrip(c.Request.Context(), tripID)
	if err != nil {
		cc.logger.Error("Failed to fetch check-ins for TripID %s: %v", tripID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve check-in data.",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": fmt.Sprintf("Successfully retrieved %d check-in records for trip %s", len(response), tripID),
		"data":    response,
	})
}
