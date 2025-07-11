// file: ticket-service/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"ticket-service/api/controllers"
	"ticket-service/api/routes"
	"ticket-service/config"
	"ticket-service/internal/consumers"
	"ticket-service/internal/db"
	"ticket-service/internal/repositories"
	"ticket-service/internal/services"
	"ticket-service/internal/workers"
	"ticket-service/pkg/emailclient"
	"ticket-service/pkg/kafkaclient"
	"ticket-service/pkg/utils"
	"ticket-service/pkg/websocket"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 1. Tải cấu hình và khởi tạo logger
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	logger := utils.NewDefaultLogger()
	logger.Info("Starting ticket service...")

	// 2. Kết nối đến các dịch vụ bên ngoài (Database, Redis)
	sqlDB, err := ConnectToDatabase(cfg)
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer sqlDB.Close()
	logger.Info("Connected to PostgreSQL")

	redisClient, err := ConnectToRedis(cfg.Redis.URL)
	if err != nil {
		logger.Error("Failed to connect to redis: %v", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	logger.Info("Connected to Redis")

	// 3. Khởi tạo các thành phần cốt lõi (Kafka, Repositories, Services)
	kafkaPublisher, err := kafkaclient.NewPublisher(cfg.Kafka, logger)
	if err != nil {
		logger.Error("Failed to create kafka publisher: %v", err)
		os.Exit(1)
	}
	defer kafkaPublisher.Close()

	auth := utils.NewAuth(cfg.JWT.SecretKey)
	util := utils.NewUtils()
	query := db.New(sqlDB)

	emailClient := emailclient.NewEmailClient(kafkaPublisher, logger, cfg)
	ticketRepo := repositories.NewTicketRepository(sqlDB, redisClient, util, logger)
	manaRepo := repositories.NewManagerTicket(sqlDB, redisClient, ticketRepo, logger)
	checkRepo := repositories.NewCheckinRepository(sqlDB, logger)

	var ticketService services.ITicketService // Khai báo trước để giải quyết phụ thuộc vòng
	manaService := services.NewManagerTicketService(manaRepo, ticketRepo, logger, cfg, kafkaPublisher, emailClient)
	ticketService = services.NewTicketService(ticketRepo, util, logger, cfg, kafkaPublisher, redisClient)
	checkService := services.NewCheckinService(checkRepo, logger)

	// 4. Khởi tạo và chạy các Worker/Consumer trong Goroutine
	consumerCtx, consumerCancel := context.WithCancel(context.Background())

	// Khởi tạo WebSocket Manager với callback
	onAckTimeout := func(bookingID string, ticketID string) {
		logger.Info("[Main] ACK Timeout triggered for booking %s, ticket %s. Initiating cancellation.", bookingID, ticketID)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := ticketService.CancelTicket(ctx, ticketID); err != nil {
				logger.Error("[Main-Error] Failed to cancel ticket %s on timeout: %v", ticketID, err)
			}
		}()
	}
	wsManager := websocket.NewManager(redisClient, onAckTimeout)

	// Chạy các worker
	outboxWorker := workers.NewOutboxPoller(query, kafkaPublisher, logger)
	go outboxWorker.Start(consumerCtx)

	timeoutWorker := workers.NewTimeoutWorker(redisClient, ticketService, logger)
	go timeoutWorker.Start(consumerCtx)

	// Chạy các consumer với context đã tạo
	tripConsumer := consumers.NewTripConsumer(cfg, manaService, logger)
	go tripConsumer.Start(consumerCtx)

	ticketStatusConsumer := consumers.NewTicketStatusConsumer(cfg, manaService, logger)
	go ticketStatusConsumer.Start(consumerCtx)

	bookingRequestConsumer := consumers.NewBookingRequestConsumer(cfg, ticketService, redisClient, logger)
	go bookingRequestConsumer.Start(consumerCtx)

	if err := ticketRepo.SubscribeToSeatStatusChanges(consumerCtx); err != nil {
		logger.Error("Failed to subscribe to seat status changes: %v", err)
	}

	// 5. Khởi tạo và chạy Web Server
	gin.SetMode(cfg.Server.GinMode)
	router := gin.Default()

	// Khởi tạo Controllers
	manaController := controllers.NewManagerTicketController(manaService)
	ticketController := controllers.NewTicketController(ticketService)
	checkController := controllers.NewCheckinController(checkService, logger)
	testController := controllers.NewTokenTestController(auth)
	routes.SetupRoutes(router, manaController, ticketController, testController, checkController, wsManager)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// 6. Chặn luồng main và xử lý shutdown an toàn (Graceful Shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Chương trình sẽ dừng ở đây cho đến khi nhận được tín hiệu (ví dụ: Ctrl+C)
	<-quit

	logger.Info("Shutting down server and background workers...")

	// Gửi tín hiệu dừng đến tất cả các goroutine đang lắng nghe consumerCtx
	consumerCancel()

	// Cho server 5 giây để xử lý các request còn lại
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	// (Tùy chọn) Chờ một chút để các goroutine dọn dẹp
	time.Sleep(1 * time.Second)
	logger.Info("Server exited properly")
}

// Hàm kết nối đến Database
func ConnectToDatabase(cfg config.Config) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)
	dbConn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = dbConn.PingContext(ctx); err != nil {
		dbConn.Close()
		return nil, err
	}
	return dbConn, nil
}

// Hàm kết nối đến Redis (Upstash)
func ConnectToRedis(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("không thể phân giải URL của Redis: %w", err)
	}
	redisClient := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("không thể kết nối đến Redis (Upstash): %w", err)
	}
	return redisClient, nil
}
