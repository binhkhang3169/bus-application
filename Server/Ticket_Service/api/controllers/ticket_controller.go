package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"ticket-service/domain/models"
	"ticket-service/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RoleCustomer  = "ROLE_CUSTOMER"
	RoleReception = "ROLE_RECEPTION"
	RoleAdmin     = "ROLE_ADMIN"
	RoleOperator  = "ROLE_OPERATOR"
)

var (
	ErrTripNotFoundFromService = errors.New("trip not found in service")
	ErrConflictFromService     = errors.New("conflict in service")
)

type TicketController struct {
	ticketService services.ITicketService
}

func NewTicketController(ticketService services.ITicketService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
	}
}

// InitiateBookingHandler được cập nhật để gọi service
func (t *TicketController) InitiateBookingHandler(c *gin.Context) {
	// 1. Lấy thông tin từ request như cũ
	var customerID *int
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr != "" {
		id, err := strconv.Atoi(userIDStr)
		if err == nil {
			customerID = &id
		}
	}

	var input models.TicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid request body: " + err.Error()})
		return
	}

	var customerNullInt sql.NullInt32
	if customerID != nil {
		customerNullInt = sql.NullInt32{Int32: int32(*customerID), Valid: true}
	}

	// 2. Tạo bookingId
	bookingID := uuid.New().String()

	// 3. GỌI SERVICE ĐỂ ĐƯA YÊU CẦU VÀO HÀNG ĐỢI
	// Sử dụng context của request cho tác vụ nhanh này
	err := t.ticketService.QueueNewBooking(c.Request.Context(), bookingID, &input, customerNullInt)
	if err != nil {
		// Nếu service không đưa vào hàng đợi được, trả về lỗi server
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to accept booking request. Please try again later.",
		})
		return
	}

	// 4. Phản hồi ngay lập tức cho client
	c.JSON(http.StatusAccepted, gin.H{
		"code":    http.StatusAccepted,
		"message": "Booking request accepted. Please track the progress via WebSocket.",
		"data": gin.H{
			"bookingId": bookingID,
		},
	})
}
func (t *TicketController) GetTicketHandler(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userRole := c.GetHeader("X-User-Role")

	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Missing X-User-ID header. Authentication required.",
			"data":    nil,
		})
		return
	}

	customerID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid X-User-ID format.",
			"data":    nil,
		})
		return
	}

	ticketID := c.Param("id")
	ticket, err := t.ticketService.GetTicketByID(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Ticket not found.",
			"data":    nil,
		})
		return
	}

	if userRole == RoleCustomer && int(ticket.CustomerID.Int32) != customerID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "Access denied. You are not authorized to view this ticket.",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Ticket retrieved successfully",
		"data": gin.H{
			"ticket":      ticket,
			"retrievedBy": customerID,
		},
	})
}

func (t *TicketController) GetInfoTicketHandler(c *gin.Context) {
	var input models.TicketInfoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	ticket, err := t.ticketService.GetInfoTicketByPhone(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Ticket not found.",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Ticket information retrieved successfully",
		"data": gin.H{
			"ticket": ticket,
		},
	})
}

func (t *TicketController) GetAllTicketHandler(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userRole := c.GetHeader("X-User-Role")

	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Missing X-User-ID header. Authentication required.",
			"data":    nil,
		})
		return
	}

	if userRole != RoleCustomer {
		// Future: Add logic for admins/staff to fetch all tickets or based on other params
	}

	customerID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid X-User-ID format.",
			"data":    nil,
		})
		return
	}

	customerIDPtr := sql.NullInt32{
		Int32: int32(customerID),
		Valid: true,
	}

	tickets, err := t.ticketService.GetTicketByCustomer(c.Request.Context(), customerIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve tickets.",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Tickets retrieved successfully",
		"data":    tickets,
	})
}

// CreateTicketHandler creates a ticket. Payment initiation is no longer part of this step.
func (t *TicketController) CreateTicketHandler(c *gin.Context) {
	var customerID *int
	userIDStr := c.GetHeader("X-User-ID")

	if userIDStr != "" {
		id, err := strconv.Atoi(userIDStr)
		if err == nil {
			customerID = &id
		} else {
			fmt.Printf("Warning: Could not parse X-User-ID '%s' for CreateTicketHandler: %v\n", userIDStr, err)
		}
	} else {
		fmt.Printf("Info: No X-User-ID header provided for CreateTicketHandler. Assuming anonymous or staff booking.\n")
	}

	var input models.TicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	if customerID == nil && input.Phone == "" {
		// This validation might change depending on your rules for anonymous vs identified bookings
		// if customerID is implicitly derived from a logged-in user and not passed in the body.
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Missing phone number for unauthenticated booking, or X-User-ID header for authenticated booking.",
			"data":    nil,
		})
		return
	}

	var customerNullInt sql.NullInt32
	if customerID != nil {
		customerNullInt = sql.NullInt32{
			Int32: int32(*customerID),
			Valid: true,
		}
	} else {
		// For unauthenticated users, customerNullInt remains Valid=false.
		// If input.CustomerID was in JSON, it's in input.CustomerID.
		// If your service logic needs to distinguish between "no customer ID header" and "customer ID 0 from body",
		// that logic would be in the service layer.
		customerNullInt = sql.NullInt32{
			Int32: 0,     // Use CustomerID from input if provided, otherwise it's 0
			Valid: false, // Valid if CustomerID from input is not 0
		}
	}

	// Call the modified CreateTicket service method
	ticket, err := t.ticketService.CreateTicket(c.Request.Context(), &input, customerNullInt)
	if err != nil {
		// Determine appropriate status code based on error type
		// For example, if err is due to seat unavailability (conflict)
		// For now, using a generic StatusConflict but can be refined
		statusCode := http.StatusConflict
		// if errors.Is(err, services.ErrSeatsNotAvailable) { // Example specific error
		// 	statusCode = http.StatusConflict
		// } else if errors.Is(err, someOtherValidationError) {
		// 	statusCode = http.StatusBadRequest
		// }
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": "Failed to create ticket: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "Ticket created successfully. Payment processing is a separate step.", // Updated message
		"data": gin.H{
			"ticket": ticket,
			// "payment_url" is removed as it's no longer returned by the service
		},
	})
}

func (t *TicketController) CreateTicketByStaffHandler(c *gin.Context) {
	staffIDStr := c.GetHeader("X-User-ID")
	staffRole := c.GetHeader("X-User-Role")

	if staffIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Missing X-User-ID header. Staff authentication required.",
			"data":    nil,
		})
		return
	}

	isAuthorizedStaff := false
	switch staffRole {
	case RoleReception, RoleAdmin, RoleOperator:
		isAuthorizedStaff = true
	}

	if !isAuthorizedStaff {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": fmt.Sprintf("Access denied. Role '%s' is not authorized for this action.", staffRole),
			"data":    nil,
		})
		return
	}

	var input models.TicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	ticket, err := t.ticketService.CreateTicketByStaff(c.Request.Context(), &input, staffIDStr)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Failed to create ticket by staff: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "Ticket created successfully by staff",
		"data":    gin.H{"ticket": ticket},
	})
}

func (t *TicketController) GetAvailableHandler(c *gin.Context) {
	id := c.Param("id")
	seats, err := t.ticketService.GetAvailableSeatsByTripID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTripNotFoundFromService) { // Example specific error from service
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "Trip not found.", "data": nil})
		} else if errors.Is(err, ErrConflictFromService) { // Example for another specific error
			c.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "message": "Could not retrieve seats due to a conflict: " + err.Error(), "data": nil})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to retrieve available seats: " + err.Error(), "data": nil})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Available seats retrieved successfully",
		"data":    gin.H{"seats": seats},
	})
}

func (t *TicketController) GetAvailableMultiTripsHandler(c *gin.Context) {
	var input struct {
		TripIDs []string `json:"trip_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	type tripResult struct {
		TripID string
		Seats  []models.SeatReturn
		Error  error
	}

	resultChan := make(chan tripResult, len(input.TripIDs))

	for _, tripID := range input.TripIDs {
		go func(id string) {
			seats, errService := t.ticketService.GetAvailableSeatsByTripID(ctx, id)
			resultChan <- tripResult{
				TripID: id,
				Seats:  seats,
				Error:  errService,
			}
		}(tripID)
	}

	resultData := make(map[string][]models.SeatReturn)
	var errorMessages []string

	for i := 0; i < len(input.TripIDs); i++ {
		select {
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":    http.StatusRequestTimeout,
				"message": "Request timed out while fetching seats for multiple trips.",
				"data":    gin.H{"retrieved_seats": resultData, "errors": append(errorMessages, "overall timeout")},
			})
			return
		case res := <-resultChan:
			if res.Error != nil {
				errMsg := fmt.Sprintf("Error fetching seats for trip %s: %v", res.TripID, res.Error)
				if errors.Is(res.Error, ErrTripNotFoundFromService) {
					errMsg = fmt.Sprintf("Trip %s not found.", res.TripID)
				}
				errorMessages = append(errorMessages, errMsg)
			} else {
				resultData[res.TripID] = res.Seats
			}
		}
	}

	if len(errorMessages) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Retrieved seats for multiple trips with some errors.",
			"data": gin.H{
				"seats":   resultData,
				"errors":  errorMessages,
				"partial": true,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Available seats for multiple trips retrieved successfully.",
		"data":    gin.H{"seats": resultData},
	})
}

func (t *TicketController) GetAllTicketsPaginatedHandler(c *gin.Context) {
	userRole := c.GetHeader("X-User-Role")

	// Authorization: Only Admins or Operators can access this
	isAuthorized := false
	switch userRole {
	case RoleAdmin, RoleOperator:
		isAuthorized = true
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": fmt.Sprintf("Access denied. Role '%s' is not authorized for this action.", userRole),
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginatedResult, err := t.ticketService.GetAllTickets(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve tickets: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All tickets retrieved successfully",
		"data":    paginatedResult,
	})
}

func (t *TicketController) GetPublicTicketInfoByIDHandler(c *gin.Context) {

	userRole := c.GetHeader("X-User-Role")

	// Authorization: Only Admins or Operators can access this
	isAuthorized := false
	switch userRole {
	case RoleAdmin, RoleOperator:
		isAuthorized = true
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": fmt.Sprintf("Access denied. Role '%s' is not authorized for this action.", userRole),
		})
		return
	}

	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Ticket ID is required.",
		})
		return
	}

	ticket, err := t.ticketService.GetTicketByID(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve ticket information.",
		})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Ticket not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Ticket information retrieved successfully",
		"data":    gin.H{"ticket": ticket},
	})
}
