package route

import (
	"bank/api/controller" // Sẽ tạo sau nếu cần
	"bank/config"
	"bank/internal/service"
	"bank/utils"

	// Cho custom validator

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"       // Thêm import này
	"github.com/go-playground/validator/v10" // Thêm import này
)

// SetupRoutes thiết lập tất cả các routes cho ứng dụng.
func SetupRoutes(router *gin.Engine, accountSvc service.AccountService, cfg config.Config) {
	// Đăng ký custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", utils.ValidCurrency) // Đăng ký validator 'currency'
	}

	// Khởi tạo controllers
	accountController := controller.NewAccountController(accountSvc)

	// Nhóm routes cho API v1
	apiV1 := router.Group("/api/v1")
	{
		// Middleware (ví dụ: logging, auth) có thể được thêm ở đây
		// apiV1.Use(middleware.LoggingMiddleware())
		// apiV1.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // Nếu có JWT

		accountRoutes := apiV1.Group("/accounts")
		{
			// === CÁC THAY ĐỔI CHÍNH NẰM Ở ĐÂY ===

			// POST /api/v1/accounts - Vẫn giữ nguyên để tạo tài khoản mới
			accountRoutes.POST("", accountController.CreateAccount)

			// GET /api/v1/accounts/me - Lấy thông tin tài khoản của user đang đăng nhập (từ header)
			// Sử dụng "/me" để rõ ràng hơn thay vì dùng chung GET /accounts
			accountRoutes.GET("/me", accountController.GetMyAccount)

			// GET /api/v1/accounts - Dành cho admin để liệt kê tất cả các tài khoản (có phân trang)
			accountRoutes.GET("", accountController.ListAccounts)

			// Các endpoint thao tác trên tài khoản của "tôi" (dựa vào header)
			accountRoutes.POST("/deposit", accountController.DepositToMyAccount)
			accountRoutes.POST("/payment", accountController.MakePaymentOnMyAccount)
			accountRoutes.PATCH("/close", accountController.CloseMyAccount)
			accountRoutes.GET("/history", accountController.GetMyTransactionHistory)

			/*
				Lưu ý: Các route cũ sử dụng /:id đã được thay thế.
				- accountRoutes.GET("/:id", ...) -> accountRoutes.GET("/me", ...)
				- accountRoutes.POST("/:id/deposit", ...) -> accountRoutes.POST("/deposit", ...)
				- ... và tương tự cho các route khác
			*/
		}

		// Thêm các nhóm route khác ở đây (ví dụ: /users, /transactions)
	}

	// Route cho health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}
