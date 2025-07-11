package routes

import (
	"net/http"
	"notification-service/api/controller"
	"notification-service/internal/service"

	// "notification-service/internal/sse" // ĐÃ XÓA

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRouter không còn nhận sseManager
func SetupRouter(dbpool *pgxpool.Pool, notificationSvc service.NotificationService) *gin.Engine {
	r := gin.Default()

	// Khởi tạo controller không cần sseManager
	notificationCtrl := controller.NewNotificationController(notificationSvc)

	apiV1 := r.Group("/api/v1")
	{
		notificationsGroup := apiV1.Group("/notifications")
		{
			notificationsGroup.POST("/broadcast", notificationCtrl.CreateBroadcastNotification)
			notificationsGroup.POST("/user", notificationCtrl.CreateUserNotification)
			notificationsGroup.GET("/broadcast", notificationCtrl.GetBroadcastNotifications)
			// Sửa đổi route này một chút để an toàn hơn, nhận userID từ body
			notificationsGroup.PUT("/:notification_id/read", notificationCtrl.MarkNotificationAsRead)
		}

		// Đổi tên group để rõ ràng hơn
		usersGroup := apiV1.Group("/usersnoti")
		{
			// Route mới để đăng ký token
			usersGroup.POST("/:user_id/fcm-token", notificationCtrl.RegisterFCMToken)

			usersGroup.GET("/:user_id/notifications", notificationCtrl.GetUserNotifications)
			usersGroup.PUT("/:user_id/notifications/read-all", notificationCtrl.MarkAllUserNotificationsAsRead)
			// usersGroup.GET("/:user_id/notifications/stream", notificationCtrl.StreamUserNotifications) // ĐÃ XÓA
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	return r
}
