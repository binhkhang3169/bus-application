// cmd/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"news-management/api/controller"
	"news-management/api/routes" // Alias để tránh trùng tên với package routes của Gin
	"news-management/config"
	"news-management/internal/repository"
	"news-management/internal/service"
	"news-management/pkg/kafkaclient"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import driver PostgreSQL
)

func main() {
	// 1. Tải cấu hình
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Kết nối Database
	// Lưu ý: sqlc.yaml đang dùng "sql_package: pgx/v5",
	// nhưng NewNewsRepository đang nhận *sql.DB.
	// Để nhất quán, nếu dùng pgx/v5 trong sqlc, bạn nên dùng pgxpool cho kết nối.
	// Ở đây, để đơn giản, ta vẫn dùng *sql.DB với lib/pq.
	// Nếu bạn muốn dùng pgxpool:
	// connPool, err := pgxpool.New(context.Background(), cfg.DBSource())
	// if err != nil {
	//  log.Fatalf("Unable to connect to database: %v\n", err)
	// }
	// defer connPool.Close()
	// dbInstance := connPool // Gán connPool cho biến sẽ truyền vào repository

	// Sử dụng database/sql với lib/pq
	db, err := sql.Open("postgres", cfg.DBSource())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Kiểm tra kết nối
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	// Cài đặt connection pool (quan trọng cho performance)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	publisher, err := kafkaclient.NewPublisher(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka publisher: %v", err)
	}
	defer publisher.Close()

	// 3. Khởi tạo các tầng
	newsRepo := repository.NewNewsRepository(db) // Truyền *sql.DB
	newsSvc := service.NewNewsService(newsRepo, publisher)
	newsCtrl := controller.NewNewsController(newsSvc)

	// 4. Khởi tạo Gin router
	router := gin.Default()

	// Middleware (ví dụ: logging, CORS - nếu cần)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 5. Đăng ký routes
	apiGroup := router.Group("/api/v1") // Base path cho API
	routes.SetupNewsRoutes(apiGroup, newsCtrl)

	// Endpoint kiểm tra sức khỏe
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// 6. Chạy server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s\n", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
