package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import driver PostgreSQL

	"payment_service/api/controller"
	"payment_service/api/route"
	"payment_service/config"
	"payment_service/internal/repository"
	"payment_service/internal/service"
	"payment_service/internal/worker"
	"payment_service/pkg/kafkaclient"
	"payment_service/pkg/redisclient"
	"payment_service/pkg/utils"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection using database/sql
	dbConn, err := ConnectToDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Ping database to ensure connection is live
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	// Initialize JWT Auth utility (nếu có)
	var authUtil *utils.Auth
	if cfg.JWT.SecretKey != "" {
		authUtil = utils.NewAuth(cfg.JWT.SecretKey)
	} else {
		log.Println("Warning: JWT_SECRET_KEY is not set. Authentication features might be limited.")
	}
	kafkaClient, err := kafkaclient.NewPublisher(cfg.KafkaConfig)
	if err != nil {
		log.Fatalf("Không thể khởi tạo Kafka client: %v", err)
	}
	defer kafkaClient.Close()
	redisClient := redisclient.NewClient(cfg.RedisConfig.URL)

	// Initialize repositories
	// Sử dụng *sql.DB cho NewInvoiceRepository
	invoiceRepo := repository.NewInvoiceRepository(dbConn)

	// Initialize services
	// Truyền interface repository cho service
	invoiceService := service.NewInvoiceService(invoiceRepo, kafkaClient, redisClient)
	vnpayService := service.NewVNPayService(&cfg.VNPay, invoiceService) // Thêm ServerConfig nếu cần cho ReturnURL
	stripeService := service.NewStripeService(&cfg.Stripe, invoiceService)
	bankService := service.NewBankService(invoiceRepo, invoiceService, "http://bank-service:8086", &http.Client{})

	// Initialize controllers
	vnpayController := controller.NewVNPayController(*vnpayService, invoiceService, &cfg.VNPay, authUtil)
	stripeController := controller.NewStripeController(stripeService, invoiceService, &cfg.Stripe)
	bankController := controller.NewBankController(bankService, invoiceService)
	staffCtrl := controller.NewStaffAssistedPaymentController(invoiceService)

	expirySubscriber := worker.NewExpirySubscriber(redisClient, invoiceService)
	go expirySubscriber.Start(context.Background())

	// Initialize Gin router
	// gin.SetMode(gin.ReleaseMode) // Chuyển sang ReleaseMode cho production
	router := gin.Default()

	// Setup routes
	route.SetupRoutes(router, vnpayController, stripeController, bankController, staffCtrl)

	// Configure server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
		// ReadTimeout:  10 * time.Second, // Thêm timeout nếu cần
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give the server 5 seconds to finish in-flight requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// ConnectToDatabase establishes a connection to the database using database/sql
func ConnectToDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	// Ví dụ DSN cho PostgreSQL: "postgres://user:password@host:port/dbname?sslmode=disable"
	// Hoặc "host=myhost port=myport user=gorm dbname=gorm password=mypassword sslmode=disable TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
	if strings.ToLower(cfg.Driver) != "postgres" {
		// Nếu bạn muốn hỗ trợ các driver khác, cần điều chỉnh DSN
		log.Printf("Warning: Database driver is '%s'. DSN format might need adjustment.", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Thiết lập connection pool (tùy chọn nhưng khuyến khích)
	db.SetMaxOpenConns(25)                 // Số lượng connection tối đa được mở
	db.SetMaxIdleConns(5)                  // Số lượng connection tối đa ở trạng thái idle
	db.SetConnMaxLifetime(5 * time.Minute) // Thời gian sống tối đa của một connection

	// Kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		db.Close() // Đóng connection nếu ping thất bại
		return nil, fmt.Errorf("failed to ping database after opening connection: %w", err)
	}

	return db, nil
}
