package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shipment-service/api/controller" // sqlc generated
	"shipment-service/api/routes"
	"shipment-service/config"
	"shipment-service/internal/repository"
	"shipment-service/internal/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	// _ "github.com/joho/godotenv/autoload" // Autoloads .env, or load manually in config
)

func main() {
	// Load configuration
	// The .env file path can be relative to the execution directory.
	// If running `go run cmd/main.go` from project root, ".env" is correct.
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBSource)
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

	// Initialize Repositories
	// Repositories are now stateless and receive the querier (db.DBTX) per call from the service.
	shipmentRepo := repository.NewShipmentRepository()
	invoiceRepo := repository.NewInvoiceRepository()

	// Initialize Services
	// The service now takes the dbpool to manage transactions.
	shipmentSvc := service.NewShipmentService(shipmentRepo, invoiceRepo, db)

	// Initialize Controllers
	shipmentCtrl := controller.NewShipmentController(shipmentSvc)

	// Setup Gin router
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, shipmentCtrl)

	// Start server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
		// ReadTimeout:  5 * time.Second, // Example timeouts
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout:  120 * time.Second,
	}

	// Goroutine for starting the server
	go func() {
		log.Printf("Server starting on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	// signal.Notify registers the given channel to receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Block until a signal is received.
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the requests it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
