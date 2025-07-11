// file: pkg/routes/router.go
package routes

import (
	// Thay thế 'api_gateway' bằng tên module go thực tế của bạn
	"api_gateway/pkg/handlers"
	"api_gateway/pkg/middleware"
	"api_gateway/pkg/services"
	"api_gateway/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter khởi tạo và cấu hình tất cả các route API cho gateway.
func SetupRouter(
	router *gin.Engine,
	authService *utils.Auth,
	serviceRegistry *services.ServiceRegistry,
) {
	// --- Định nghĩa quyền RBAC ---
	// Map này định nghĩa các vai trò được phép cho các tiền tố đường dẫn (path prefixes) cụ thể.
	routeRolePermissions := map[string][]string{
		// Actions của User/Customer
		"/api/v1/auth/logout":          {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN", "ROLE_DRIVER"},
		"/api/v1/change-password":      {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN", "ROLE_DRIVER"},
		"/api/v1/customer/info":        {"ROLE_CUSTOMER", "ROLE_ADMIN", "ROLE_RECEPTION"},
		"/api/v1/customer/change-info": {"ROLE_CUSTOMER", "ROLE_ADMIN", "ROLE_RECEPTION"},
		"/api/v1/initiate-booking":     {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN", "ROLE_GUEST"}, // ĐÃ THÊM

		// Actions của Admin/Operator/Reception
		"/api/v1/create":        {"ROLE_ADMIN"},
		"/api/v1/users/by-role": {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/staff/tickets": {"ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN"},
		"/api/v1/create-seats":  {"ROLE_ADMIN", "ROLE_OPERATOR"},

		// Actions của Ticket:
		// - POST /api/v1/tickets là public cho guests.
		// - GET /api/v1/tickets được bảo vệ cho người dùng đã xác thực (để xem vé của chính họ).
		"/api/v1/tickets":           {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN"}, // ĐÃ BỎ COMMENT
		"/api/v1/tickets/all":       {"ROLE_ADMIN", "ROLE_OPERATOR", "ROLE_RECEPTION"},                  // Lấy tất cả vé (có phân trang)
		"/api/v1/public/ticket/:id": {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN"}, // Lấy thông tin vé công khai
		// Actions khác
		"/api/v1/checkin":    {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_DRIVER"},
		"/api/v1/token-test": {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN", "ROLE_DRIVER"},
		"/api/v1/employee":   {"ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN", "ROLE_DRIVER"},
		"/api/v1/invoices":   {"ROLE_CUSTOMER", "ROLE_ADMIN", "ROLE_RECEPTION"},
		"/api/v1/vnpay":      {"ROLE_CUSTOMER"},
		"/api/v1/email":      {"ROLE_CUSTOMER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN"},
		"/api/v1/ws/track":   {"ROLE_CUSTOMER", "ROLE_DRIVER", "ROLE_RECEPTION", "ROLE_OPERATOR", "ROLE_ADMIN", "ROLE_GUEST"}, // ĐÃ THÊM

		// Quản lý các entity của TripService (CUD được bảo vệ)
		"/api/v1/vehicles":     {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/routes":       {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/types":        {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/special-days": {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/stations":     {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/pickups":      {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/trips/driver": {"ROLE_DRIVER", "ROLE_OPERATOR"},
		"/api/v1/trips":        {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/provinces":    {"ROLE_ADMIN", "ROLE_OPERATOR"}, // ĐÃ THÊM

		// Accounts, News, Notifications, Shipments
		"/api/v1/accounts":       {"ROLE_CUSTOMER"},
		"/api/v1/notifications":  {"ROLE_ADMIN", "ROLE_OPERATOR", "ROLE_CUSTOMER"},
		"/api/v1/usersnoti":      {"ROLE_ADMIN", "ROLE_OPERATOR", "ROLE_CUSTOMER", "ROLE_GUEST"},
		"/api/v1/tripsshipments": {"ROLE_ADMIN", "ROLE_OPERATOR", "ROLE_RECEPTION"},
		"/api/v1/news":           {"ROLE_ADMIN", "ROLE_OPERATOR", "ROLE_CUSTOMER"},
		"/api/v1/shipments":      {"ROLE_ADMIN", "ROLE_OPERATOR", "ROLE_RECEPTION"},
		"/api/v1/shipinvoices":   {"ROLE_ADMIN", "ROLE_OPERATOR", "ROLE_RECEPTION"},

		//Dashboard
		"/api/v1/kpis":      {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/charts":    {"ROLE_ADMIN", "ROLE_OPERATOR"},
		"/api/v1/analytics": {"ROLE_ADMIN", "ROLE_OPERATOR"},
	}

	// Khởi tạo AuthMiddleware (kết hợp xác thực và phân quyền)
	authMw := middleware.AuthMiddleware(authService, routeRolePermissions)

	// --- Route của Gateway ---
	router.GET("/gateway/health", handlers.HealthCheckHandler)
	router.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "API Gateway is Active and Ready! 🌐"}) })

	apiWebhook := router.Group("webhooks/rest/webhook")

	{
		apiWebhook.POST("/", serviceRegistry.ProxyHandler) // Webhook cho Chat Bot
	}
	// --- API v1 Routes ---
	apiV1 := router.Group("/api/v1")

	// --- Public Routes (Không cần AuthMiddleware - Guest có thể truy cập) ---
	{
		// Xác thực
		publicAuthGroup := apiV1.Group("/auth")
		{
			publicAuthGroup.POST("/login", serviceRegistry.ProxyHandler)
			publicAuthGroup.POST("/refresh-token", serviceRegistry.ProxyHandler)
		}
		apiV1.POST("/signup", serviceRegistry.ProxyHandler)
		apiV1.POST("/forgot-password", serviceRegistry.ProxyHandler)
		apiV1.POST("/reset-password", serviceRegistry.ProxyHandler)
		apiV1.POST("/verify-otp", serviceRegistry.ProxyHandler)
		apiV1.POST("/resend-otp", serviceRegistry.ProxyHandler)

		// Public GET cho các entities
		apiV1.GET("/locations", serviceRegistry.ProxyHandler)
		apiV1.GET("/provinces", serviceRegistry.ProxyHandler)
		apiV1.GET("/provinces/:id", serviceRegistry.ProxyHandler)
		apiV1.GET("/routes", serviceRegistry.ProxyHandler)
		apiV1.GET("/routes/:id", serviceRegistry.ProxyHandler)
		apiV1.GET("/trips", serviceRegistry.ProxyHandler)
		apiV1.GET("/trips/search", serviceRegistry.ProxyHandler)
		apiV1.GET("/trips/:id", serviceRegistry.ProxyHandler)
		apiV1.GET("/trips/:id/seats", serviceRegistry.ProxyHandler)
		apiV1.GET("/stations", serviceRegistry.ProxyHandler)
		apiV1.GET("/stations/:id", serviceRegistry.ProxyHandler)
		apiV1.GET("/vehicles", serviceRegistry.ProxyHandler)
		apiV1.GET("/vehicles/:id", serviceRegistry.ProxyHandler)
		apiV1.GET("/types", serviceRegistry.ProxyHandler)
		apiV1.GET("/types/:id", serviceRegistry.ProxyHandler)
		apiV1.GET("/pickups", serviceRegistry.ProxyHandler)
		apiV1.GET("/pickups/:id", serviceRegistry.ProxyHandler)
		apiV1.GET("/special-days", serviceRegistry.ProxyHandler)
		apiV1.GET("/special-days/:id", serviceRegistry.ProxyHandler)

		apiV1.GET("/news", serviceRegistry.ProxyHandler)
		apiV1.GET("/news/:id", serviceRegistry.ProxyHandler)
		// News (Protected CUD)
		apiV1.POST("/news", authMw[0], authMw[1], serviceRegistry.ProxyHandler)
		apiV1.PUT("/news/:id", authMw[0], authMw[1], serviceRegistry.ProxyHandler)
		apiV1.DELETE("/news/:id", authMw[0], authMw[1], serviceRegistry.ProxyHandler)

		// Ticket & Trip (Public)
		apiV1.POST("/tickets", serviceRegistry.ProxyHandler) // SỬA ĐỔI: Guest có thể tạo vé
		apiV1.POST("/ticket-by-phone", serviceRegistry.ProxyHandler)
		apiV1.GET("/tickets-available/:id", serviceRegistry.ProxyHandler)
		apiV1.POST("/trips-available-seats", serviceRegistry.ProxyHandler)

		// Payment (Public)
		apiV1.POST("/payments", serviceRegistry.ProxyHandler)
		apiV1.GET("/payments", serviceRegistry.ProxyHandler)
		apiV1.GET("/invoices", serviceRegistry.ProxyHandler)
		apiV1.GET("/vnpay/return", serviceRegistry.ProxyHandler)
		apiV1.POST("/stripe/confirm-payment", serviceRegistry.ProxyHandler)
		apiV1.POST("/vnpay/create-payment", serviceRegistry.ProxyHandler)
		apiV1.POST("/stripe/create-payment-intent", serviceRegistry.ProxyHandler)
		apiV1.POST("/bank/create-payment-request", serviceRegistry.ProxyHandler)
		apiV1.POST("/bank/confirm-payment", serviceRegistry.ProxyHandler)
		apiV1.POST("/bank/payment-failed", serviceRegistry.ProxyHandler)
		apiV1.POST("/invoices/fail", serviceRegistry.ProxyHandler)

		// QR & Upload (Public)
		apiV1.GET("/qr/image", serviceRegistry.ProxyHandler)
		apiV1.GET("/qr/url", serviceRegistry.ProxyHandler)
		apiV1.POST("/upload/image", serviceRegistry.ProxyHandler)
	}

	// --- Protected Routes (Áp dụng AuthMiddleware - Yêu cầu đăng nhập) ---
	createUpdateDeleteHandler := serviceRegistry.ProxyHandler // Alias cho rõ ràng

	// General protected actions
	apiV1.POST("/auth/logout", authMw[0], authMw[1], serviceRegistry.ProxyHandler)
	apiV1.POST("/change-password", authMw[0], authMw[1], serviceRegistry.ProxyHandler)
	apiV1.POST("/email", authMw[0], authMw[1], serviceRegistry.ProxyHandler)
	apiV1.POST("/vnpay", authMw[0], authMw[1], serviceRegistry.ProxyHandler)

	// Websocket (Protected)
	websocketGroup := apiV1.Group("/ws")
	websocketGroup.Use(authMw...) // Áp dụng middleware cho cả group
	{
		websocketGroup.GET("/track/:bookingId", serviceRegistry.ProxyHandler)
	}

	// Initiate Booking (Protected)
	initticket := apiV1.Group("/initiate-booking")
	initticket.Use(authMw...)
	{
		initticket.POST("", serviceRegistry.ProxyHandler)
	}

	// Customer specific (Protected)
	customerProtected := apiV1.Group("/customer")
	customerProtected.Use(authMw...)
	{
		customerProtected.GET("/info", serviceRegistry.ProxyHandler)
		customerProtected.POST("/change-info", serviceRegistry.ProxyHandler)
	}

	// Protected CUD operations for entities
	tripProtected := apiV1.Group("/trips")
	tripProtected.Use(authMw...)
	{
		tripProtected.POST("", createUpdateDeleteHandler)
		tripProtected.PUT("/:id", createUpdateDeleteHandler)
		tripProtected.DELETE("/:id", createUpdateDeleteHandler)
		tripProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		tripProtected.GET("/status/:status", createUpdateDeleteHandler)
		tripProtected.GET("/driver", createUpdateDeleteHandler)
	}

	pickupProtected := apiV1.Group("/pickups")
	pickupProtected.Use(authMw...)
	{
		pickupProtected.POST("", createUpdateDeleteHandler)
		pickupProtected.PUT("/:id", createUpdateDeleteHandler)
		pickupProtected.DELETE("/:id", createUpdateDeleteHandler)
		pickupProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		pickupProtected.GET("/status/:status", createUpdateDeleteHandler)
		pickupProtected.GET("/byRoute/:id", createUpdateDeleteHandler)
	}

	provinceProtected := apiV1.Group("/provinces")
	provinceProtected.Use(authMw...)
	{
		provinceProtected.POST("", createUpdateDeleteHandler)
		provinceProtected.PUT("/:id", createUpdateDeleteHandler)
		provinceProtected.DELETE("/:id", createUpdateDeleteHandler)
		provinceProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		provinceProtected.GET("/status/:status", createUpdateDeleteHandler)
	}

	vehicleProtected := apiV1.Group("/vehicles")
	vehicleProtected.Use(authMw...)
	{
		vehicleProtected.POST("", createUpdateDeleteHandler)
		vehicleProtected.PUT("/:id", createUpdateDeleteHandler)
		vehicleProtected.DELETE("/:id", createUpdateDeleteHandler)
		vehicleProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		vehicleProtected.GET("/status/:status", createUpdateDeleteHandler)
	}

	routePathProtected := apiV1.Group("/routes")
	routePathProtected.Use(authMw...)
	{
		routePathProtected.POST("", createUpdateDeleteHandler)
		routePathProtected.PUT("/:id", createUpdateDeleteHandler)
		routePathProtected.DELETE("/:id", createUpdateDeleteHandler)
		routePathProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		routePathProtected.GET("/status/:status", createUpdateDeleteHandler)
	}

	specialDayProtected := apiV1.Group("/special-days")
	specialDayProtected.Use(authMw...)
	{
		specialDayProtected.POST("", createUpdateDeleteHandler)
		specialDayProtected.PUT("/:id", createUpdateDeleteHandler)
		specialDayProtected.DELETE("/:id", createUpdateDeleteHandler)
		specialDayProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		specialDayProtected.GET("/status/:status", createUpdateDeleteHandler)
	}

	stationProtected := apiV1.Group("/stations")
	stationProtected.Use(authMw...)
	{
		stationProtected.POST("", createUpdateDeleteHandler)
		stationProtected.PUT("/:id", createUpdateDeleteHandler)
		stationProtected.DELETE("/:id", createUpdateDeleteHandler)
		stationProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		// apiV1.GET("/stations/status/:status") đã là public, không cần định nghĩa lại ở đây
	}

	typeProtected := apiV1.Group("/types")
	typeProtected.Use(authMw...)
	{
		typeProtected.POST("", createUpdateDeleteHandler)
		typeProtected.PUT("/:id", createUpdateDeleteHandler)
		typeProtected.DELETE("/:id", createUpdateDeleteHandler)
		typeProtected.PUT("/:id/status/:status", createUpdateDeleteHandler)
		typeProtected.GET("/status/:status", createUpdateDeleteHandler)
	}

	//Dashboard protected routes
	dashboardProtected := apiV1.Group("/kpis")
	dashboardProtected.Use(authMw...)
	{
		dashboardProtected.GET("", serviceRegistry.ProxyHandler)
	}
	chartsProtected := apiV1.Group("/charts")
	chartsProtected.Use(authMw...)
	{
		chartsProtected.GET("/revenue-over-time", serviceRegistry.ProxyHandler)
		chartsProtected.GET("/ticket-distribution", serviceRegistry.ProxyHandler)
	}
	analyticsProtected := apiV1.Group("/analytics")
	analyticsProtected.Use(authMw...)
	{
		searches := analyticsProtected.Group("/searches")
		{
			searches.GET("/top-routes", serviceRegistry.ProxyHandler)
			searches.GET("/top-provinces", serviceRegistry.ProxyHandler)
			searches.GET("/by-hour", serviceRegistry.ProxyHandler)
			searches.GET("/over-time", serviceRegistry.ProxyHandler)
		}
	}

	// User management by admin/operator
	apiV1.POST("/create", authMw[0], authMw[1], serviceRegistry.ProxyHandler)
	apiV1.GET("/users/by-role", authMw[0], authMw[1], serviceRegistry.ProxyHandler)

	// Ticket service (Protected actions)
	ticketActionsProtected := apiV1.Group("/tickets")
	ticketActionsProtected.Use(authMw...)
	{
		// SỬA ĐỔI: POST đã được chuyển ra public.
		// Group này giờ chỉ bảo vệ các action cho user đã đăng nhập.
		ticketActionsProtected.GET("", serviceRegistry.ProxyHandler)     // Lấy vé của user
		ticketActionsProtected.GET("/:id", serviceRegistry.ProxyHandler) // Lấy chi tiết vé
		ticketActionsProtected.GET("/all", serviceRegistry.ProxyHandler) // Lấy tất cả vé (có phân trang)
	}

	//User get our ticket
	apiV1.GET("/public/ticket/:id", authMw[0], authMw[1], serviceRegistry.ProxyHandler) // Lấy thông tin vé công khai

	apiV1.POST("/staff/tickets", authMw[0], authMw[1], serviceRegistry.ProxyHandler)
	apiV1.POST("/create-seats", authMw[0], authMw[1], serviceRegistry.ProxyHandler)

	checkinAPI := apiV1.Group("/checkin")
	checkinAPI.Use(authMw...)
	{
		// SỬA ĐỔI: Không cần áp dụng authMw lần nữa vì group đã có
		checkinAPI.POST("/", serviceRegistry.ProxyHandler)
		checkinAPI.GET("/trip/:tripID", serviceRegistry.ProxyHandler)
	}

	apiV1.GET("/token-test", authMw[0], authMw[1], serviceRegistry.ProxyHandler)

	employee := apiV1.Group("/employee")
	employee.Use(authMw...)
	{
		employee.GET("/info", serviceRegistry.ProxyHandler)
	}

	// Shipments (Protected)
	tripShipments := apiV1.Group("/tripsshipments/:trip_id")
	tripShipments.Use(authMw...)
	{
		tripShipments.POST("/shipments", serviceRegistry.ProxyHandler) // Tạo shipment cho trip
		tripShipments.GET("/shipments", serviceRegistry.ProxyHandler)  // Lấy
		tripShipments.GET("/invoices", serviceRegistry.ProxyHandler)   // Lấy tất cả invoices cho trip
	}

	shipmentSpecific := apiV1.Group("/shipments")
	shipmentSpecific.Use(authMw...)
	{
		shipmentSpecific.GET("", serviceRegistry.ProxyHandler) // Tạo shipment
		shipmentSpecific.GET("/:id/invoice", serviceRegistry.ProxyHandler)
		shipmentSpecific.GET("/:id", serviceRegistry.ProxyHandler) // Lấy tất cả shipments với phân trang
	}

	// --- Invoice specific routes ---
	invoicesGroup := apiV1.Group("/shipinvoices")
	invoicesGroup.Use(authMw...)
	{
		invoicesGroup.GET("", serviceRegistry.ProxyHandler)     // Lấy tất cả invoices với phân trang
		invoicesGroup.GET("/:id", serviceRegistry.ProxyHandler) // Lấy một invoice theo ID
	}
	// Bank Account (Protected)
	accountRoutes := apiV1.Group("/accounts")
	accountRoutes.Use(authMw...)
	{
		accountRoutes.POST("", serviceRegistry.ProxyHandler)
		accountRoutes.GET("/me", serviceRegistry.ProxyHandler)
		accountRoutes.GET("", serviceRegistry.ProxyHandler)
		accountRoutes.POST("/deposit", serviceRegistry.ProxyHandler)
		accountRoutes.POST("/payment", serviceRegistry.ProxyHandler)
		accountRoutes.PATCH("/close", serviceRegistry.ProxyHandler)
		accountRoutes.GET("/history", serviceRegistry.ProxyHandler)
	}

	// Notifications (Protected)
	notificationsGroup := apiV1.Group("/notifications")
	notificationsGroup.Use(authMw...)
	{
		notificationsGroup.POST("/broadcast", serviceRegistry.ProxyHandler)
		notificationsGroup.POST("/user", serviceRegistry.ProxyHandler)
		notificationsGroup.GET("/broadcast", serviceRegistry.ProxyHandler)
		notificationsGroup.PUT("/:notification_id/read", serviceRegistry.ProxyHandler)
	}

	usersGroup := apiV1.Group("/usersnoti")
	usersGroup.Use(authMw...)
	{
		usersGroup.POST("/:user_id/fcm-token", serviceRegistry.ProxyHandler)
		usersGroup.GET("/:user_id/notifications", serviceRegistry.ProxyHandler)
		usersGroup.PUT("/:user_id/notifications/read-all", serviceRegistry.ProxyHandler)
	}
}
