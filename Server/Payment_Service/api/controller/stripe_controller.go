package controller

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"payment_service/config"
	"payment_service/domain/model"
	"payment_service/internal/service"
	"payment_service/pkg/utils"
)

// StripeController xử lý các API endpoint liên quan đến Stripe
type StripeController struct {
	stripeService  service.StripeServiceInterface
	invoiceService service.InvoiceServiceInterface // Keep if general invoice actions needed by controller
	cfg            *config.StripeConfig
}

// NewStripeController tạo một StripeController mới
func NewStripeController(stripeService service.StripeServiceInterface, invoiceService service.InvoiceServiceInterface, cfg *config.StripeConfig) *StripeController {
	return &StripeController{
		stripeService:  stripeService,
		invoiceService: invoiceService, // Retain for potential direct invoice operations
		cfg:            cfg,
	}
}

// CreatePaymentIntent ... (remains the same)
func (c *StripeController) CreatePaymentIntent(ctx *gin.Context) {
	var req model.InitialStripePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// Basic validation (Stripe has its own minimums, e.g., $0.50)
	// currencyLowerCase := strings.ToLower(req.Currency)
	// if (currencyLowerCase == "usd" && req.Amount < 50) || (currencyLowerCase == "vnd" && req.Amount < 10000) {
	//    utils.RespondWithError(ctx, http.StatusBadRequest, "Amount is below minimum", nil)
	//    return
	// }

	resp, err := c.stripeService.CreatePaymentIntent(ctx, req)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to create Payment Intent", err.Error())
		return
	}

	utils.RespondWithSuccess(ctx, http.StatusOK, "Payment Intent created successfully", resp)
}

// ConfirmStripePayment ... (remains the same)
func (c *StripeController) ConfirmStripePayment(ctx *gin.Context) {
	var req model.StripePaymentConfirmationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.PaymentIntentID == "" {
		utils.RespondWithError(ctx, http.StatusBadRequest, "PaymentIntentID is required", nil)
		return
	}

	invoice, err := c.stripeService.ConfirmPayment(ctx, req.PaymentIntentID)
	if err != nil {
		// Error could be API error from Stripe, or PI not succeeded, or DB update error
		// Example: Distinguish client errors vs server errors
		// if strings.Contains(err.Error(), "stripe: payment for PI") && strings.Contains(err.Error(), "is still processing") {
		// 	utils.RespondWithError(ctx, http.StatusAccepted, "Payment is still processing", err.Error()) // 202 Accepted
		// 	return
		// } else if strings.Contains(err.Error(), "stripe: payment for PI") && strings.Contains(err.Error(), "was not successful") {
		// 	utils.RespondWithError(ctx, http.StatusBadRequest, "Payment was not successful", err.Error())
		// 	return
		// }
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to confirm Stripe payment", err.Error())
		return
	}

	apiInvoice := c.invoiceService.MapDbInvoiceToAPIResponse(invoice) // Use injected invoiceService
	utils.RespondWithSuccess(ctx, http.StatusOK, "Stripe payment confirmed successfully", apiInvoice)
}

// HandleStripeWebhook ... (remains the same)
func (c *StripeController) HandleStripeWebhook(ctx *gin.Context) {
	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusServiceUnavailable, "Error reading webhook body", err.Error())
		return
	}

	signature := ctx.GetHeader("Stripe-Signature")
	if signature == "" && c.cfg.WebhookSecret != "" { // Only require if secret is configured
		utils.RespondWithError(ctx, http.StatusBadRequest, "Missing Stripe-Signature header", nil)
		return
	}

	if err := c.stripeService.HandleWebhook(ctx, payload, signature); err != nil {
		log.Printf("Error processing Stripe webhook: %v", err) // Log the actual error
		// Respond with BadRequest for signature errors, InternalServerError for others
		// This detail is usually handled within the service or based on error type.
		utils.RespondWithError(ctx, http.StatusBadRequest, "Error processing Stripe webhook", err.Error())
		return
	}

	utils.RespondWithSuccess(ctx, http.StatusOK, "Webhook received and processed", nil)
}

// RefundStripePayment handles initiating a Stripe refund.
// POST /api/v1/stripe/refund
func (c *StripeController) RefundStripePayment(ctx *gin.Context) {
	var req model.StripeInitiateRefundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload for Stripe refund", err.Error())
		return
	}

	// Call the stripe service to process the refund
	// The service will handle fetching invoice, calculating amount, calling Stripe API, and updating DB.
	updatedDbInvoice, err := c.stripeService.RefundPayment(ctx, req)
	if err != nil {
		// Service layer should return specific errors that can be mapped to HTTP status codes
		// e.g., if invoice not found, already refunded, Stripe API error, DB update error after successful refund.
		log.Printf("Error during Stripe refund for ticket %s: %v", req.TicketID, err)
		// Example of more granular error handling:
		// if errors.Is(err, service.ErrInvoiceNotFound) {
		// 	utils.RespondWithError(ctx, http.StatusNotFound, "Invoice not found for refund", err.Error())
		// } else if errors.Is(err, service.ErrInvoiceAlreadyRefunded) {
		// 	utils.RespondWithError(ctx, http.StatusBadRequest, "Invoice already refunded", err.Error())
		// } else if strings.Contains(err.Error(), "stripe API error") {
		// 	utils.RespondWithError(ctx, http.StatusFailedDependency, "Stripe API error during refund", err.Error()) // 424
		// } else if strings.Contains(err.Error(), "MANUAL INTERVENTION REQUIRED") {
		//  utils.RespondWithError(ctx, http.StatusInternalServerError, "Refund processed by Stripe but failed to update local data. Please check.", err.Error())
		// }
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to process Stripe refund", err.Error())
		return
	}

	// Map the updated DB invoice to API response
	apiInvoiceResponse := c.invoiceService.MapDbInvoiceToAPIResponse(updatedDbInvoice)
	utils.RespondWithSuccess(ctx, http.StatusOK, "Stripe refund processed successfully", apiInvoiceResponse)
}
