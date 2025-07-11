package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// ========= CÁC STRUCT DÙNG CHUNG CHO EMAIL =========

// EmailRequest is the structure of the data sent via Kafka.
type EmailRequest struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"` // Contains the JSON string of the data for the template
	Type  string `json:"type"`
}

// Struct for OTP email
type OTPData struct {
	Data         string `json:"data"`
	CustomerName string `json:"customerName"`
	CurrentDate  string `json:"currentDate"`
}

// Struct for staff account email
type StaffAccountData struct {
	CustomerName string `json:"customerName"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	CurrentDate  string `json:"currentDate"`
}

// Struct for a single ticket in the order confirmation email
type TicketData struct {
	SeatInfo  string `json:"seatInfo"`
	QRCodeURL string `json:"qrCodeURL"`
}

// Struct for the order confirmation email
type OrderConfirmationData struct {
	CustomerName string       `json:"customerName"`
	OrderID      string       `json:"orderId"`
	TotalPrice   float64      `json:"totalPrice"`
	Tickets      []TicketData `json:"tickets"`
}

// ========= CÁC BIẾN CẤU HÌNH TOÀN CỤC =========

var (
	smtpServer, smtpPort, smtpUsername, smtpPassword, smtpFrom string
	kafkaSeeds                                                 []string
	kafkaEmailTopic, kafkaGroupID                              string
	kafkaEnableTLS                                             bool
	kafkaSASLUser, kafkaSASLPass                               string
)

// ========= HÀM MAIN VÀ KHỞI TẠO =========

func main() {
	// Load environment variables from .env file.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables from OS.")
	}

	// Read and initialize configuration.
	initConfig()

	// Run Kafka consumer in a goroutine to not block the main thread.
	go setupKafkaConsumer()

	// Set up Gin web server.
	r := gin.Default()
	api := r.Group("/api/v1")
	{
		// This endpoint is still useful for testing or sending mail directly (synchronously).
		api.POST("/email", sendEmailHandler)
	}
	r.GET("/health", healthCheckHandler)

	// Start the server.
	port := getEnv("PORT", "8085")
	log.Printf("Email microservice starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initConfig reads configuration from environment variables.
func initConfig() {
	smtpServer = getEnv("SMTP_SERVER", "smtp.gmail.com")
	smtpPort = getEnv("SMTP_PORT", "587")
	smtpUsername = getEnv("SMTP_USERNAME", "")
	smtpPassword = getEnv("SMTP_PASSWORD", "")
	smtpFrom = getEnv("SMTP_FROM", "")

	// Read new Kafka configuration for franz-go
	kafkaSeeds = strings.Split(getEnv("KAFKA_SEEDS", "localhost:9092"), ",")
	kafkaEmailTopic = getEnv("KAFKA_EMAIL_TOPIC", "email_requests")
	kafkaGroupID = getEnv("KAFKA_GROUP_ID_EMAIL", "email_service_group")
	kafkaEnableTLS, _ = strconv.ParseBool(getEnv("KAFKA_ENABLE_TLS", "false"))
	kafkaSASLUser = getEnv("KAFKA_SASL_USER", "")
	kafkaSASLPass = getEnv("KAFKA_SASL_PASS", "")

	if smtpUsername == "" || smtpPassword == "" || smtpFrom == "" {
		log.Fatal("FATAL: SMTP credentials (SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM) are not configured.")
	}
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		log.Fatal("FATAL: 'templates' directory not found. HTML email sending will fail.")
	}
}

// getEnv is a helper function to safely read environment variables.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ========= HTTP HANDLERS =========

func sendEmailHandler(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sendEmail(req); err != nil {
		log.Printf("Error in sendEmailHandler: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "up",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ========= LOGIC GỬI EMAIL =========

// sendEmail is the core logic function to process and dispatch an email.
func sendEmail(req EmailRequest) error {
	var templateFile string
	var templateData interface{}

	// Choose template and unmarshal JSON data based on req.Type.
	switch req.Type {
	case "otp":
		templateFile = "otp.html"
		var data OTPData
		if err := json.Unmarshal([]byte(req.Body), &data); err != nil {
			return fmt.Errorf("failed to unmarshal OTPData for type %s: %w", req.Type, err)
		}
		templateData = data

	case "staff_account_info":
		templateFile = "staff_account_info.html"
		var data StaffAccountData
		if err := json.Unmarshal([]byte(req.Body), &data); err != nil {
			return fmt.Errorf("failed to unmarshal StaffAccountData for type %s: %w", req.Type, err)
		}
		templateData = data

	case "order_confirmation_payment":
		templateFile = "order_confirmation.html"
		var data OrderConfirmationData
		if err := json.Unmarshal([]byte(req.Body), &data); err != nil {
			return fmt.Errorf("failed to unmarshal OrderConfirmationData for type %s: %w", req.Type, err)
		}
		templateData = data

	default:
		return fmt.Errorf("unsupported email type: %s", req.Type)
	}

	// Parse and execute the HTML template.
	templatePath := filepath.Join("templates", templateFile)

	// The template uses the 'printf' function which is not a default.
	// We must create a new template, add the function to its map, and then parse the file.
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"printf": fmt.Sprintf,
	}).ParseFiles(templatePath)

	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var renderedBody bytes.Buffer
	if err := tmpl.Execute(&renderedBody, templateData); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	// Send email via SMTP.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	addr := smtpServer + ":" + smtpPort
	subject := "Subject: " + req.Title + "\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	message := []byte(subject + mime + "\r\n" + renderedBody.String())

	if err := smtp.SendMail(addr, auth, smtpFrom, []string{req.To}, message); err != nil {
		return fmt.Errorf("smtp.SendMail failed: %w", err)
	}

	log.Printf("Successfully sent email to %s (Type: %s)", req.To, req.Type)
	return nil
}

// ========= KAFKA CONSUMER VỚI FRANZ-GO =========

// setupKafkaConsumer listens for email requests from Kafka.
func setupKafkaConsumer() {
	opts := []kgo.Opt{
		kgo.SeedBrokers(kafkaSeeds...),
		kgo.ConsumerGroup(kafkaGroupID),
		kgo.ConsumeTopics(kafkaEmailTopic),
		// Disable auto-commit to manually control offset committing
		kgo.DisableAutoCommit(),
	}

	if kafkaEnableTLS {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}
	if kafkaSASLUser != "" && kafkaSASLPass != "" {
		opts = append(opts, kgo.SASL(scram.Auth{
			User: kafkaSASLUser,
			Pass: kafkaSASLPass,
		}.AsSha256Mechanism()))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Fatalf("Consumer: Failed to create Kafka client: %v", err)
	}
	defer client.Close()

	log.Println("Kafka consumer started. Waiting for messages on topic:", kafkaEmailTopic)
	ctx := context.Background()

	for {
		fetches := client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			// Non-retriable errors are logged here
			log.Printf("Consumer: Error polling fetches: %v", errs)
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var req EmailRequest
			if err := json.Unmarshal(record.Value, &req); err != nil {
				log.Printf("Consumer: Error unmarshalling message: %v. Discarding (poison pill).", err)
				// Commit the faulty message to avoid reprocessing
				if err := client.CommitRecords(ctx, record); err != nil {
					log.Printf("Consumer: Failed to commit poison pill message: %v", err)
				}
				continue
			}

			log.Printf("Consumer: Received email request for %s (Type: %s)", req.To, req.Type)

			if err := sendEmail(req); err != nil {
				log.Printf("Consumer: Error sending email for %s: %v. Message will NOT be committed and will be retried.", req.To, err)
				// Do not commit on error, so Kafka can redeliver the message
			} else {
				log.Printf("Consumer: Successfully processed email request for %s. Committing offset.", req.To)
				// Manually commit the offset after successful processing
				if err := client.CommitRecords(ctx, record); err != nil {
					log.Printf("Consumer: Failed to commit message after successful processing: %v", err)
				}
			}
		}
	}
}
