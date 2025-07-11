package main

import (
	"go-bigquery-dashboard/internal/handler"
	"go-bigquery-dashboard/internal/repository"
	"log"

	// <-- 1. Import thư viện CORS
	"github.com/gin-gonic/gin"
)

const credentialsFile = "credentials.json"

func main() {
	// Khởi tạo Repositories
	dashboardRepo, err := repository.NewDashboardRepository(credentialsFile)
	if err != nil {
		log.Fatalf("Failed to create dashboard repository: %v", err)
	}
	defer dashboardRepo.Close()

	tripSearchRepo, err := repository.NewTripSearchRepository(credentialsFile)
	if err != nil {
		log.Fatalf("Failed to create trip search repository: %v", err)
	}
	defer tripSearchRepo.Close()

	// Khởi tạo Handlers
	dashboardHandler := handler.NewDashboardHandler(dashboardRepo)
	tripSearchHandler := handler.NewTripSearchHandler(tripSearchRepo)

	// Thiết lập Gin Router
	r := gin.Default()

	// // --- 2. Thêm Middleware CORS ---
	// // Cấu hình để chỉ cho phép frontend từ localhost:3000 truy cập
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// --- 3. Định tuyến API với tiền tố /api/v1 ---
	api := r.Group("/api/v1")
	{
		// API phân tích tài chính
		api.GET("/kpis", dashboardHandler.GetKPIs)
		charts := api.Group("/charts")
		{
			charts.GET("/revenue-over-time", dashboardHandler.GetRevenueOverTime)
			charts.GET("/ticket-distribution", dashboardHandler.GetTicketDistribution)
		}

		// API phân tích tìm kiếm
		analytics := api.Group("/analytics")
		{
			searches := analytics.Group("/searches")
			{
				searches.GET("/top-routes", tripSearchHandler.GetTopRoutes)
				searches.GET("/top-provinces", tripSearchHandler.GetTopProvinces)
				searches.GET("/by-hour", tripSearchHandler.GetSearchesByHour)
				searches.GET("/over-time", tripSearchHandler.GetSearchesOverTime)
			}
		}
	}

	// Chạy server
	log.Println("Starting server on :8091")
	if err := r.Run(":8091"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
