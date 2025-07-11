package controller

import (
	"net/http"
	"payment_service/domain/model"
	"payment_service/internal/service"
	"payment_service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// StaffAssistedPaymentController handles HTTP requests for payments assisted by staff.
type StaffAssistedPaymentController struct {
	invoiceService service.InvoiceServiceInterface
	// Add other services if needed, e.g., authService for staff authentication
}

// NewStaffAssistedPaymentController creates a new StaffAssistedPaymentController.
func NewStaffAssistedPaymentController(invoiceService service.InvoiceServiceInterface) *StaffAssistedPaymentController {
	return &StaffAssistedPaymentController{
		invoiceService: invoiceService,
	}
}

// HandleDirectPayment godoc
// @Summary Process Direct Payment by Staff
// @Description Allows a staff member to record a payment made directly (e.g., cash, POS).
// @Tags payments-staff
// @Accept json
// @Produce json
// @Param payment_request body model.StaffDirectPaymentRequest true "Direct Payment Request"
// @Success 201 {object} utils.SuccessResponse{data=model.GetInvoiceResponse} "Direct payment processed successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload"
// @Failure 500 {object} utils.ErrorResponse "Failed to process direct payment"
// @Router /staff-payments/direct-payment [post]
// @Security ApiKeyAuth // TODO: Add appropriate staff authentication middleware
func (c *StaffAssistedPaymentController) HandleDirectPayment(ctx *gin.Context) {
	var req model.StaffDirectPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// TODO: Implement staff authentication and authorization here.
	// For example, check if the staff_id from token matches or has permission.

	dbInvoice, err := c.invoiceService.ProcessStaffDirectPayment(ctx.Request.Context(), req)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to process staff direct payment", err.Error())
		return
	}

	apiInvoiceResponse := c.invoiceService.MapDbInvoiceToAPIResponse(dbInvoice)
	utils.RespondWithSuccess(ctx, http.StatusCreated, "Direct payment processed successfully and invoice created", apiInvoiceResponse)
}
