package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"payment_service/domain/model" // Your domain models
	"payment_service/internal/service"
	"payment_service/pkg/utils" // Your response utility
)

// BankController handles API endpoints related to bank payments.
type BankController struct {
	bankService    service.BankServiceInterface
	invoiceService service.InvoiceServiceInterface // For mapping responses, etc.
}

// NewBankController creates a new BankController.
func NewBankController(bankService service.BankServiceInterface, invoiceService service.InvoiceServiceInterface) *BankController {
	return &BankController{
		bankService:    bankService,
		invoiceService: invoiceService,
	}
}

// CreateBankPaymentRequestHandler godoc
// @Summary Request Bank Payment Details
// @Description Initiates a bank payment process by creating an invoice and returning bank details for the user to make a transfer.
// @Tags payments-bank
// @Accept json
// @Produce json
// @Param payment_request body model.InitialBankPaymentRequest true "Bank Payment Request"
// @Success 201 {object} utils.SuccessResponse{data=model.BankPaymentDetailsResponse} "Bank payment details provided"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload"
// @Failure 500 {object} utils.ErrorResponse "Failed to create bank payment request"
// @Router /bank/create-payment-request [post]
func (c *BankController) CreateBankPaymentRequestHandler(ctx *gin.Context) {
	var req model.InitialBankPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// Basic validation (can be more extensive)
	if req.Amount <= 0 {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Amount must be positive", nil)
		return
	}

	bankPaymentDetails, err := c.bankService.CreateBankPaymentRequest(ctx.Request.Context(), req)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to create bank payment request", err.Error())
		return
	}

	utils.RespondWithSuccess(ctx, http.StatusCreated, "Bank payment request initiated. Please use the provided details to make the transfer.", bankPaymentDetails)
}

// ConfirmBankPaymentHandler godoc
// @Summary Confirm Bank Payment
// @Description Allows an admin or an automated process to confirm receipt of a bank transfer.
// @Tags payments-bank
// @Accept json
// @Produce json
// @Param confirmation_request body model.BankPaymentConfirmationRequest true "Bank Payment Confirmation"
// @Success 200 {object} utils.SuccessResponse{data=model.InvoiceAPIResponse} "Bank payment confirmed successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload"
// @Failure 404 {object} utils.ErrorResponse "Invoice not found or not in a confirmable state"
// @Failure 500 {object} utils.ErrorResponse "Failed to confirm bank payment"
// @Router /bank/confirm-payment [post]
// @Security ApiKeyAuth // Example: This endpoint should be protected
func (c *BankController) ConfirmBankPaymentHandler(ctx *gin.Context) {
	var req model.BankPaymentConfirmationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload for confirmation", err.Error())
		return
	}

	// TODO: Add authentication and authorization for this endpoint.
	// Only authorized personnel or systems should be able to confirm payments.

	updatedDbInvoice, err := c.bankService.ConfirmBankPayment(ctx.Request.Context(), req)
	if err != nil {
		// Handle specific errors from service, e.g., invoice not found, wrong status
		// For now, a generic internal server error or bad request based on error content.
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to confirm bank payment", err.Error())
		return
	}

	// Map the db.Invoice to your API response model (e.g., model.InvoiceAPIResponse)
	// This mapping function would likely be in your invoiceService or a utility.
	apiInvoiceResponse := c.invoiceService.MapDbInvoiceToAPIResponse(updatedDbInvoice) // Assuming you have this method

	utils.RespondWithSuccess(ctx, http.StatusOK, "Bank payment confirmed successfully", apiInvoiceResponse)
}

// HandleBankPaymentFailedHandler godoc
// @Summary Mark Bank Payment as Failed
// @Description Allows an admin or an automated process to mark a bank payment as failed.
// @Tags payments-bank
// @Accept json
// @Produce json
// @Param failure_request body model.BankPaymentFailureRequest true "Bank Payment Failure Details" // Define this model
// @Success 200 {object} utils.SuccessResponse{data=model.InvoiceAPIResponse} "Bank payment marked as failed"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload"
// @Failure 404 {object} utils.ErrorResponse "Invoice not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to mark payment as failed"
// @Router /bank/payment-failed [post]
// @Security ApiKeyAuth // Example: This endpoint should be protected
func (c *BankController) HandleBankPaymentFailedHandler(ctx *gin.Context) {
	// Define model.BankPaymentFailureRequest:
	// type BankPaymentFailureRequest struct {
	//    InvoiceID uuid.UUID `json:"invoice_id" binding:"required"`
	//    Reason    string    `json:"reason" binding:"required"`
	// }
	var req struct { // Inline struct for brevity, define properly in models
		InvoiceID uuid.UUID `json:"invoice_id" binding:"required"`
		Reason    string    `json:"reason" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload for failure", err.Error())
		return
	}

	// TODO: Add authentication and authorization.

	failedDbInvoice, err := c.bankService.HandleBankPaymentFailed(ctx.Request.Context(), req.InvoiceID, req.Reason)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to mark bank payment as failed", err.Error())
		return
	}

	apiInvoiceResponse := c.invoiceService.MapDbInvoiceToAPIResponse(failedDbInvoice)
	utils.RespondWithSuccess(ctx, http.StatusOK, "Bank payment marked as failed", apiInvoiceResponse)
}
