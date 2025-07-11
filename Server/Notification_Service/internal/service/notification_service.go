package service

import (
	"context"
	"fmt"
	"log"
	"notification-service/internal/db"
	"notification-service/internal/model"
	"notification-service/internal/repository"
	"strings"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Helper function to check if error indicates invalid token
func isTokenInvalid(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "registration-token-not-registered") ||
		strings.Contains(errStr, "invalid-registration-token") ||
		strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "unregistered")
}

// Interface được cập nhật với phương thức mới
type NotificationService interface {
	CreateNotification(ctx context.Context, req model.CreateNotificationRequest) (db.Notification, error)
	RegisterFCMToken(ctx context.Context, userID string, token string) error
	GetNotificationsForUser(ctx context.Context, userID string, limit, offset int32) ([]db.Notification, error)
	GetBroadcastNotifications(ctx context.Context, limit, offset int32) ([]db.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationIDStr string, userID string) (db.Notification, error)
	MarkAllUserNotificationsAsRead(ctx context.Context, userID string) ([]db.Notification, error)
}

// struct không còn sseManager
type notificationService struct {
	repo      repository.Store
	fcmClient *messaging.Client
	fsClient  *firestore.Client
}

// NewNotificationService không còn nhận sseManager
func NewNotificationService(repo repository.Store, fcm *messaging.Client, firestore *firestore.Client) NotificationService {
	return &notificationService{
		repo:      repo,
		fcmClient: fcm,
		fsClient:  firestore,
	}
}

// Phương thức mới để đăng ký token
func (s *notificationService) RegisterFCMToken(ctx context.Context, userID string, token string) error {
	params := db.RegisterFCMTokenParams{
		UserID: userID,
		Token:  token,
	}
	_, err := s.repo.RegisterFCMToken(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to register FCM token in db: %w", err)
	}
	log.Printf("Successfully registered/updated FCM token for user %s", userID)
	return nil
}

// CreateNotification được tái cấu trúc hoàn toàn
func (s *notificationService) CreateNotification(ctx context.Context, req model.CreateNotificationRequest) (db.Notification, error) {
	var createdNotification db.Notification

	// 1. Lưu thông báo vào DB (PostgreSQL)
	err := s.repo.ExecTx(ctx, func(qtx *db.Queries) error {
		var pgUserID pgtype.Text
		if req.UserID != nil && *req.UserID != "" {
			pgUserID = pgtype.Text{String: *req.UserID, Valid: true}
		} else {
			pgUserID = pgtype.Text{Valid: false} // For broadcast
		}

		params := db.CreateNotificationParams{
			UserID:  pgUserID,
			Type:    req.Type,
			Title:   req.Title,
			Message: req.Message,
		}

		var errTx error
		createdNotification, errTx = qtx.CreateNotification(ctx, params)
		return errTx
	})

	if err != nil {
		return db.Notification{}, fmt.Errorf("failed to create notification in transaction: %w", err)
	}
	log.Printf("Successfully created notification %s in DB", createdNotification.ID.String())

	// 2. Lưu thông báo vào Firestore (để client có thể lắng nghe real-time nếu cần)
	_, _, err = s.fsClient.Collection("notifications").Add(ctx, map[string]interface{}{
		"id":         createdNotification.ID.String(),
		"user_id":    createdNotification.UserID.String,
		"type":       createdNotification.Type,
		"title":      createdNotification.Title,
		"message":    createdNotification.Message,
		"is_read":    false,
		"created_at": createdNotification.CreatedAt,
	})
	if err != nil {
		log.Printf("Failed to add notification to Firestore (non-critical): %v", err)
	}

	// 3. Lấy tokens và gửi Push Notification qua FCM
	var tokens []string
	if req.UserID != nil && *req.UserID != "" { // Gửi cho user cụ thể
		tokens, err = s.repo.GetFCMTokensByUserID(ctx, *req.UserID)
		if err != nil {
			log.Printf("Could not get FCM tokens for user %s: %v", *req.UserID, err)
		}
	} else { // Gửi broadcast
		tokens, err = s.repo.GetAllFCMTokens(ctx)
		if err != nil {
			log.Printf("Could not get all FCM tokens for broadcast: %v", err)
		}
	}

	if len(tokens) > 0 {
		log.Printf("Sending notification to %d tokens", len(tokens))
		s.sendFCMToTokens(ctx, tokens, createdNotification.Title, createdNotification.Message, createdNotification.ID.String())
	} else {
		log.Printf("No FCM tokens found for this request.")
	}

	return createdNotification, nil
}

// sendFCMToTokens: Send FCM to tokens individually (replacing deprecated batch method)
func (s *notificationService) sendFCMToTokens(ctx context.Context, tokens []string, title, body, notificationID string) {
	successCount := 0
	failureCount := 0
	var invalidTokens []string

	for _, token := range tokens {
		// Create individual message
		message := &messaging.Message{
			Data: map[string]string{
				"title":          title,
				"body":           body,
				"notificationId": notificationID,
				"click_action":   "FLUTTER_NOTIFICATION_CLICK",
			},
			Token: token,
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Alert: &messaging.ApsAlert{
							Title: title,
							Body:  body,
						},
						ContentAvailable: true,
					},
				},
			},
		}

		// Send individual message
		response, err := s.fcmClient.Send(ctx, message)
		if err != nil {
			failureCount++
			log.Printf("Failed to send FCM to token %s: %v", token, err)

			// Check if token is invalid based on error message
			if messaging.IsInvalidArgument(err) || isTokenInvalid(err) {
				invalidTokens = append(invalidTokens, token)
				log.Printf("Token %s is invalid/unregistered", token)
			}
		} else {
			successCount++
			log.Printf("Successfully sent FCM to token %s, message ID: %s", token, response)
		}
	}

	log.Printf("FCM sending completed. Success: %d, Failure: %d", successCount, failureCount)

	// Clean up invalid tokens
	if len(invalidTokens) > 0 {
		log.Printf("Found %d invalid/unregistered tokens to be cleaned up", len(invalidTokens))
		go func() {
			for _, token := range invalidTokens {
				// Note: You'll need to add DeleteFCMToken method to your repository
				// For now, we'll just log the tokens that should be deleted
				log.Printf("Should delete invalid token: %s", token)
				// Uncomment when you have the DeleteFCMToken method:
				// err := s.repo.DeleteFCMToken(context.Background(), token)
				// if err != nil {
				//     log.Printf("Failed to delete invalid token %s: %v", token, err)
				// } else {
				//     log.Printf("Successfully deleted invalid token %s", token)
				// }
			}
		}()
	}
}

// Các phương thức Get/MarkRead/MarkAllRead không thay đổi
func (s *notificationService) GetNotificationsForUser(ctx context.Context, userID string, limit, offset int32) ([]db.Notification, error) {
	pgUserID := pgtype.Text{String: userID, Valid: true}
	params := db.GetNotificationsByUserIDParams{
		UserID: pgUserID,
		Limit:  limit,
		Offset: offset,
	}
	return s.repo.GetNotificationsByUserID(ctx, params)
}

func (s *notificationService) GetBroadcastNotifications(ctx context.Context, limit, offset int32) ([]db.Notification, error) {
	params := db.GetBroadcastNotificationsParams{
		Limit:  limit,
		Offset: offset,
	}
	return s.repo.GetBroadcastNotifications(ctx, params)
}

func (s *notificationService) MarkNotificationAsRead(ctx context.Context, notificationIDStr string, userID string) (db.Notification, error) {
	notificationIDUUID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		log.Printf("Invalid notification ID format: %v", err)
		return db.Notification{}, err
	}

	pgNotificationID := pgtype.UUID{Bytes: notificationIDUUID, Valid: true}
	pgUserID := pgtype.Text{String: userID, Valid: true}

	params := db.MarkNotificationAsReadParams{
		ID:     pgNotificationID,
		UserID: pgUserID,
	}
	return s.repo.MarkNotificationAsRead(ctx, params)
}

func (s *notificationService) MarkAllUserNotificationsAsRead(ctx context.Context, userID string) ([]db.Notification, error) {
	pgUserID := pgtype.Text{String: userID, Valid: true}
	return s.repo.MarkAllUserNotificationsAsRead(ctx, pgUserID)
}
