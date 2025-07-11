package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notification-service/api/routes"
	"notification-service/config"
	"notification-service/internal/kafka"
	"notification-service/internal/repository"
	"notification-service/internal/service"

	// "notification-service/internal/sse" // ĐÃ XÓA
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	firebase "firebase.google.com/go"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"
)

func main() {
	// Khởi tạo Firebase
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error initializing FCM client: %v\n", err)
	}

	fsClient, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v\n", err)
	}
	defer fsClient.Close()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Kết nối DB
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslmode)
	dbpool, err := pgxpool.New(context.Background(), dbConnectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()
	log.Println("Successfully connected to the database!")

	// Initialize SSE Manager // ĐÃ XÓA
	// sseManager := sse.NewSSEManager()

	// Khởi tạo Repository, Service (không còn sseManager)
	store := repository.NewStore(dbpool)
	notificationService := service.NewNotificationService(store, fcmClient, fsClient)

	// Khởi tạo Gin router (không còn sseManager)
	router := routes.SetupRouter(dbpool, notificationService)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	// Khởi chạy các goroutine (Kafka consumer, HTTP server)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go kafka.StartConsumer(ctx, &wg, cfg, notificationService)

	go func() {
		log.Printf("HTTP server starting on port %s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", cfg.HTTPPort, err)
		}
	}()

	// Xử lý shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server and consumer...")

	cancel()

	shutdownCtx, serverShutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer serverShutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Waiting for Kafka consumer to shut down...")
	wg.Wait()
	log.Println("Kafka consumer shut down.")
	log.Println("Server exiting")
}
