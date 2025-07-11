// file: main.go
package main

import (
	// Thay thế 'api_gateway' bằng tên module go thực tế của bạn
	"api_gateway/config"
	"api_gateway/pkg/routes"
	"api_gateway/pkg/services"
	"api_gateway/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Tải cấu hình
	config.LoadEnv()
	jwtCfg := config.LoadJWTConfig()
	serviceURLs := config.LoadServiceURLs()
	serverCfg := config.LoadServerConfig()

	// 2. Khởi tạo Gin Router
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 3. Khởi tạo Service Xác thực
	authService := utils.NewAuth(jwtCfg.SecretKey)

	// 4. Khởi tạo Service Registry và đăng ký các microservice
	registry := services.NewServiceRegistry()
	registerAllServices(registry, serviceURLs) // Tách ra hàm riêng cho gọn gàng

	// 5. Setup Routes
	// Truyền các thành phần đã khởi tạo vào hàm cấu hình router.
	routes.SetupRouter(router, authService, registry)

	// 6. Khởi động Server
	log.Printf("✅ API Gateway is starting on port %s...", serverCfg.Port)
	if err := router.Run(":" + serverCfg.Port); err != nil {
		log.Fatalf("❌ API Gateway failed to start: %v", err)
	}
}

// registerAllServices gom việc đăng ký các service vào một nơi
func registerAllServices(registry *services.ServiceRegistry, serviceURLs config.ServiceURLs) {
	// Email Service
	registry.RegisterService("email-service", serviceURLs.EmailServiceURL, "/api/v1/email", 1)

	// Payment Services
	registry.RegisterService("payment-service-vnpay", serviceURLs.PaymentServiceURL, "/api/v1/vnpay", 2)
	registry.RegisterService("payment-service-stripe", serviceURLs.PaymentServiceURL, "/api/v1/stripe", 2)
	registry.RegisterService("payment-service-invoices", serviceURLs.PaymentServiceURL, "/api/v1/invoices", 2)
	registry.RegisterService("payment-service-generic", serviceURLs.PaymentServiceURL, "/api/v1/payments", 1)
	registry.RegisterService("payment-service-bank", serviceURLs.PaymentServiceURL, "/api/v1/bank", 1)

	// Trip Services
	registry.RegisterService("trip-service-locations", serviceURLs.TripServiceURL, "/api/v1/locations", 1)
	registry.RegisterService("trip-service-pickups", serviceURLs.TripServiceURL, "/api/v1/pickups", 1)
	registry.RegisterService("trip-service-provinces", serviceURLs.TripServiceURL, "/api/v1/provinces", 1)
	registry.RegisterService("trip-service-vehicles", serviceURLs.TripServiceURL, "/api/v1/vehicles", 1)
	registry.RegisterService("trip-service-routes", serviceURLs.TripServiceURL, "/api/v1/routes", 1)
	registry.RegisterService("trip-service-special-days", serviceURLs.TripServiceURL, "/api/v1/special-days", 1)
	registry.RegisterService("trip-service-stations", serviceURLs.TripServiceURL, "/api/v1/stations", 1)
	registry.RegisterService("trip-service-trips", serviceURLs.TripServiceURL, "/api/v1/trips", 1)
	registry.RegisterService("trip-service-types", serviceURLs.TripServiceURL, "/api/v1/types", 1)

	// Ticket Services
	registry.RegisterService("ticket-service-initiate-booking", serviceURLs.TicketServiceURL, "/api/v1/initiate-booking", 2) // NEW
	registry.RegisterService("ticket-service-track", serviceURLs.TicketServiceURL, "/api/v1/ws/track", 2)                    // NEW
	registry.RegisterService("ticket-service-main", serviceURLs.TicketServiceURL, "/api/v1/tickets", 1)
	registry.RegisterService("ticket-service-phone", serviceURLs.TicketServiceURL, "/api/v1/ticket-by-phone", 2)
	registry.RegisterService("ticket-service-staff", serviceURLs.TicketServiceURL, "/api/v1/staff/tickets", 2)
	registry.RegisterService("ticket-service-available", serviceURLs.TicketServiceURL, "/api/v1/tickets-available", 2)
	registry.RegisterService("ticket-service-trips-seats", serviceURLs.TicketServiceURL, "/api/v1/trips-available-seats", 2)
	registry.RegisterService("ticket-service-create-seats", serviceURLs.TicketServiceURL, "/api/v1/create-seats", 2)
	registry.RegisterService("ticket-service-checkin", serviceURLs.TicketServiceURL, "/api/v1/checkin", 2)
	registry.RegisterService("ticket-service-token-test", serviceURLs.TicketServiceURL, "/api/v1/token-test", 2)
	registry.RegisterService("ticket-service-checkin-trip", serviceURLs.TicketServiceURL, "/api/v1/public", 2)

	// User Services
	registry.RegisterService("user-service-auth", serviceURLs.UserServiceURL, "/api/v1/auth", 2)
	registry.RegisterService("user-service-password", serviceURLs.UserServiceURL, "/api/v1/change-password", 2)
	registry.RegisterService("user-service-signup", serviceURLs.UserServiceURL, "/api/v1/signup", 1)
	registry.RegisterService("user-service-public", serviceURLs.UserServiceURL, "/api/v1/", 1) // Rule chung cho các public user routes
	registry.RegisterService("user-service-users", serviceURLs.UserServiceURL, "/api/v1/create", 1)
	registry.RegisterService("user-service-verify", serviceURLs.UserServiceURL, "/api/v1/employee", 1)

	//Bank Services
	registry.RegisterService("bank-service-accounts", serviceURLs.BankServiceURL, "/api/v1/accounts", 1)
	//News Services
	registry.RegisterService("news-service-news", serviceURLs.NewsServiceURL, "/api/v1/news", 1)
	//Notification services
	registry.RegisterService("notification-service-notifications", serviceURLs.NotificationServiceURL, "/api/v1/notifications", 1)
	registry.RegisterService("notification-service-users", serviceURLs.NotificationServiceURL, "/api/v1/usersnoti", 1)

	//Shipment services
	registry.RegisterService("shipment-service-shipments", serviceURLs.ShipServiceURL, "/api/v1/shipments", 1)
	registry.RegisterService("shipment-service-invoices", serviceURLs.ShipServiceURL, "/api/v1/tripsshipments", 1)
	registry.RegisterService("shipment-service-trips", serviceURLs.ShipServiceURL, "/api/v1/shipinvoices", 1)

	//Qr services
	registry.RegisterService("qr-service-qr", serviceURLs.QrServiceURL, "/api/v1/qr", 1)
	registry.RegisterService("qr-service-upload", serviceURLs.QrServiceURL, "/api/v1/upload", 1)
	// Chat services
	registry.RegisterService("chat-service-messages", serviceURLs.ChatServiceURL, "/webhooks/rest/webhook", 1)

	registry.RegisterService("dashboard-service-kpis", serviceURLs.DashboardURL, "/api/v1/kpis", 1)
	registry.RegisterService("dashboard-service-revenue", serviceURLs.DashboardURL, "/api/v1/charts", 1)
	registry.RegisterService("dashboard-service-ticket", serviceURLs.DashboardURL, "/api/v1/analytics", 1)
}
