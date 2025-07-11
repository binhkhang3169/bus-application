package controller

import (
	"log"
	"net/http"
	"notification-service/internal/model"
	"notification-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Đã loại bỏ sseManager
type NotificationController struct {
	service service.NotificationService
}

// NewNotificationController không còn nhận sseManager
func NewNotificationController(svc service.NotificationService) *NotificationController {
	return &NotificationController{service: svc}
}

// ========= HANDLER MỚI =========

// RegisterFCMToken godoc
// @Summary Register a new FCM token for a user
// @Description Associates a Firebase Cloud Messaging (FCM) token with a user ID.
// @Tags users
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Param token_request body model.RegisterFCMTokenRequest true "FCM Token"
// @Success 200 {object} gin.H{"message": "Token registered successfully"}
// @Failure 400 {object} gin.H{"error": "string"}
// @Failure 500 {object} gin.H{"error": "string"}
// @Router /users/{user_id}/fcm-token [post]
func (c *NotificationController) RegisterFCMToken(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var req model.RegisterFCMTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := c.service.RegisterFCMToken(ctx.Request.Context(), userID, req.Token); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register FCM token: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Token registered successfully"})
}

// CreateBroadcastNotification không thay đổi
func (c *NotificationController) CreateBroadcastNotification(ctx *gin.Context) {
	var req model.CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	req.UserID = nil // Ensure it's a broadcast

	notification, err := c.service.CreateNotification(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, notification)
}

// CreateUserNotification không thay đổi
func (c *NotificationController) CreateUserNotification(ctx *gin.Context) {
	var req model.CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if req.UserID == nil || *req.UserID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required for user-specific notification"})
		return
	}

	notification, err := c.service.CreateNotification(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, notification)
}

// Các handler khác (GetUserNotifications, GetBroadcastNotifications, ...) giữ nguyên
// ... (code các handler khác không đổi)
func (c *NotificationController) GetUserNotifications(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	notifications, err := c.service.GetNotificationsForUser(ctx.Request.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notifications)
}

func (c *NotificationController) GetBroadcastNotifications(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	notifications, err := c.service.GetBroadcastNotifications(ctx.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve broadcast notifications: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notifications)
}

func (c *NotificationController) MarkNotificationAsRead(ctx *gin.Context) {
	notificationID := ctx.Param("notification_id")
	// Lấy userID từ body hoặc authenticated user context thay vì query param sẽ an toàn hơn
	var body struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required in body"})
		return
	}

	notification, err := c.service.MarkNotificationAsRead(ctx.Request.Context(), notificationID, body.UserID)
	if err != nil {
		log.Printf("Error in MarkNotificationAsRead controller: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notification)
}

func (c *NotificationController) MarkAllUserNotificationsAsRead(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	updatedNotifications, err := c.service.MarkAllUserNotificationsAsRead(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all notifications as read: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedNotifications)
}
