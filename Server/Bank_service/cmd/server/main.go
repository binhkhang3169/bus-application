package main

import (
	"log"
	"net/http"
	"time"

	"bank/api/route"
	"bank/config" // Sẽ được tạo bởi sqlc
	"bank/internal/repository"
	"bank/internal/service" // Thư mục util cho các hàm tiện ích (ví dụ: random)
	"bank/pkg/kafkaclient"
	"bank/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Tải cấu hình
	cfg, err := config.LoadConfig(".") // Đọc từ file .env trong thư mục gốc
	if err != nil {
		log.Fatalf("Không thể tải cấu hình: %v", err)
	}

	// Kết nối CSDL
	conn, err := utils.ConnectDB(cfg.DBDriver, cfg.DBSource())
	if err != nil {
		log.Fatalf("Không thể kết nối tới CSDL: %v", err)
	}
	defer conn.Close()

	kafkaClient, err := kafkaclient.NewPublisher(cfg)
	if err != nil {
		log.Fatalf("Không thể khởi tạo Kafka client: %v", err)
	}
	defer kafkaClient.Close()

	// Khởi tạo store (chứa các query được sinh bởi sqlc)
	store := repository.NewStore(conn) // db.NewStore sẽ được sqlc tạo ra

	// Khởi tạo các tầng
	accountRepo := repository.NewAccountRepository(store)
	accountSvc := service.NewAccountService(accountRepo, kafkaClient)
	// Khởi tạo các service khác nếu có...

	// Thiết lập Gin router
	router := gin.Default()

	// Setup Gin server
	server := &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Setup routes
	// Truyền các service cần thiết vào route setup
	route.SetupRoutes(router, accountSvc, cfg) // Truyền cfg nếu cần thiết cho middleware hoặc controller

	log.Printf("Server đang chạy tại địa chỉ %s", cfg.ServerAddress())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Lỗi khi khởi chạy server: %s\n", err)
	}
}
