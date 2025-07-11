package controller

import (
	"fmt"
	"log"
	"net/http"
	"payment_service/domain/model"
	"payment_service/pkg/utils"
	"strings"

	"payment_service/config"
	"payment_service/internal/service" // Ensure this points to your service package

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VNPayController xử lý các API endpoint liên quan đến VNPay
type VNPayController struct {
	vnpaySvc   service.VNPayService            // Use interface
	invoiceSvc service.InvoiceServiceInterface // Use interface
	cfg        *config.VNPayConfig
	auth       *utils.Auth // Placeholder for auth utility
}

// NewVNPayController tạo một VNPayController mới
func NewVNPayController(
	vnpaySvc service.VNPayService,
	invoiceSvc service.InvoiceServiceInterface,
	cfg *config.VNPayConfig,
	auth *utils.Auth,
) *VNPayController {
	return &VNPayController{
		vnpaySvc:   vnpaySvc,
		invoiceSvc: invoiceSvc,
		cfg:        cfg,
		auth:       auth,
	}
}

// CreatePayment xử lý yêu cầu tạo thanh toán VNPay
// POST /api/v1/vnpay/create-payment
func (c *VNPayController) CreatePayment(ctx *gin.Context) {
	var req model.VNPayPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// clientIP := ctx.ClientIP() // Gin's ClientIP() might need trusted proxies configuration
	// For simplicity, this part is often handled in service or passed directly if available
	// For now, vnpaySvc CreatePayment might use a default or expect it if critical

	resp, err := c.vnpaySvc.CreatePayment(ctx, req) // Assuming service handles IP
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to create VNPay payment", err.Error())
		return
	}

	utils.RespondWithSuccess(ctx, http.StatusOK, "VNPay payment URL created successfully", resp)
}

// HandleReturn xử lý khi VNPay redirect về Return URL
// GET /api/v1/vnpay/return
func (c *VNPayController) HandleReturn(ctx *gin.Context) {
	queryParams := ctx.Request.URL.Query()

	// The service method will validate signature and update DB
	resp, err := c.vnpaySvc.ProcessReturn(ctx, queryParams)
	if err != nil {
		// Error might be due to invalid signature, DB update failure, etc.
		// The error message from service should be descriptive.
		log.Printf("Error processing VNPay return: %v. Query: %s", err, queryParams.Encode())
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to process VNPay return", err.Error())
		return
	}

	// Typically redirect to a frontend page:
	// frontendURL := fmt.Sprintf("%s?invoice_id=%s&status=%s&message=%s",
	//    c.cfg.ClientReturnURL, resp.InvoiceID.String(), resp.ResponseCode, url.QueryEscape(resp.Message))
	// ctx.Redirect(http.StatusFound, frontendURL)
	// For API-only service, returning JSON is fine:
	utils.RespondWithSuccess(ctx, http.StatusOK, "VNPay return processed", resp)
}

// HandleIPN xử lý Instant Payment Notification từ VNPay
// GET /api/v1/vnpay/ipn (VNPay typically uses GET for IPN)
func (c *VNPayController) HandleIPN(ctx *gin.Context) {
	queryParams := ctx.Request.URL.Query()
	log.Printf("Received VNPay IPN: %s", queryParams.Encode())

	ipnRespModel, err := c.vnpaySvc.ProcessIPN(ctx, queryParams)
	if err != nil {
		// This error is from the service layer, indicating a system problem during IPN processing (e.g., DB error)
		// ProcessIPN itself should log specifics. The controller just relays a generic server error to VNPay.
		log.Printf("CRITICAL: System error processing VNPay IPN: %v. Query: %s", err, queryParams.Encode())
		// VNPay expects a specific JSON response format for IPN.
		// Even with internal error, respond in VNPay's required format.
		// RspCode "99" is often for system errors from merchant side.
		ctx.JSON(http.StatusOK, model.VNPayIPNResponse{RspCode: "99", Message: "Internal Server Error"})
		return
	}

	// Send the structured response back to VNPay
	log.Printf("Responding to VNPay IPN with RspCode: %s, Message: %s", ipnRespModel.RspCode, ipnRespModel.Message)
	ctx.JSON(http.StatusOK, ipnRespModel)
}

// GetInvoice lấy thông tin hóa đơn theo ID
// GET /api/v1/invoices/:id
func (c *VNPayController) GetInvoice(ctx *gin.Context) {
	idStr := ctx.Param("id")
	invoiceID, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid invoice ID format", err.Error())
		return
	}

	dbInvoice, err := c.invoiceSvc.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondWithError(ctx, http.StatusNotFound, "Invoice not found", nil)
		} else {
			utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve invoice", err.Error())
		}
		return
	}

	apiInvoice := c.invoiceSvc.MapDbInvoiceToAPIResponse(dbInvoice)
	utils.RespondWithSuccess(ctx, http.StatusOK, "Invoice retrieved successfully", apiInvoice)
}

// GetInvoicesByCustomer ... (remains the same)
func (c *VNPayController) GetInvoicesByCustomer(ctx *gin.Context) {
	customerIDStr := ctx.Query("customer_id")
	if customerIDStr == "" {
		// Basic auth / JWT logic placeholder
		// customerIDFromToken, err := c.auth.GetCustomerIDFromToken(ctx.GetHeader("Authorization"))
		// if err != nil {
		//    utils.RespondWithError(ctx, http.StatusUnauthorized, "Unauthorized", nil)
		//    return
		// }
		// customerIDStr = customerIDFromToken
		utils.RespondWithError(ctx, http.StatusBadRequest, "Customer ID is required via query param for this example", nil)
		return
	}

	dbInvoices, err := c.invoiceSvc.GetInvoicesByCustomerID(ctx, customerIDStr)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve invoices for customer", err.Error())
		return
	}

	if len(dbInvoices) == 0 {
		utils.RespondWithSuccess(ctx, http.StatusOK, "No invoices found for this customer", []model.GetInvoiceResponse{})
		return
	}

	apiInvoices := c.invoiceSvc.MapDbInvoicesToAPIResponses(dbInvoices)
	utils.RespondWithSuccess(ctx, http.StatusOK, "Invoices retrieved successfully", apiInvoices)
}

// VNPayQueryTransaction ... (remains similar, service prepares data)
func (c *VNPayController) VNPayQueryTransaction(ctx *gin.Context) {
	var req model.VNPayQueryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload for VNPay query", err.Error())
		return
	}
	clientIP := ctx.ClientIP()
	if clientIP == "" {
		clientIP = "127.0.0.1" // Fallback
	}

	// The service prepares the data to be sent to VNPay
	queryData, err := c.vnpaySvc.QueryTransaction(ctx, req, clientIP)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to prepare VNPay query transaction data", err.Error())
		return
	}

	// In a real scenario, the backend might make the call to VNPay:
	// vnpayResponse, err := c.vnpaySvc.ExecuteQuery(ctx, queryData) ...
	// For now, returning prepared data and VNPay API URL
	utils.RespondWithSuccess(ctx, http.StatusOK, "VNPay query data prepared. POST this to VNPay.", gin.H{
		"prepared_data": queryData,
		"vnpay_api_url": c.cfg.TransactionAPI, // Assuming cfg.TransactionAPI is for query/refund
	})
}

// VNPayRefundTransaction handles initiating a VNPay refund.
// POST /api/v1/vnpay/refund-transaction
func (c *VNPayController) VNPayRefundTransaction(ctx *gin.Context) {
	var req model.VNPayInitiateRefundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload for VNPay refund initiation", err.Error())
		return
	}

	clientIP := ctx.ClientIP()
	if clientIP == "" {
		clientIP = "127.0.0.1" // Fallback
	}

	// 1. Get the completed invoice by TicketID
	invoice, err := c.invoiceSvc.GetLatestCompletedInvoiceByTicketID(ctx, req.TicketID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondWithError(ctx, http.StatusNotFound, fmt.Sprintf("No completed invoice found for ticket ID %s to refund", req.TicketID), err.Error())
		} else {
			utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve invoice for refund", err.Error())
		}
		return
	}

	// Ensure it's a VNPay invoice
	if !invoice.PaymentMethod.Valid || invoice.PaymentMethod.String != string(model.PaymentMethodVNPay) || !invoice.VnpayTxnRef.Valid {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invoice is not a refundable VNPay transaction", nil)
		return
	}
	if invoice.PaymentStatus.String == string(model.PaymentStatusRefunded) {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invoice is already refunded", nil)
		return
	}
	if invoice.PaymentStatus.String != string(model.PaymentStatusCompleted) {
		utils.RespondWithError(ctx, http.StatusBadRequest, fmt.Sprintf("Invoice status is %s, cannot refund", invoice.PaymentStatus.String), nil)
		return
	}

	// 2. Calculate refund amount
	refundAmount := invoice.FinalAmount * (1.0 - req.PercentageDeduction)
	if refundAmount < 0 { // Should not happen with percentage_deduction <= 1
		refundAmount = 0
	}
	// VNPay amounts are typically positive. Round to 2 decimal places for float comparison.
	refundAmount = utils.RoundFloat(refundAmount, 2)

	// 3. Determine TransactionType
	var transactionType string
	// Compare float64 carefully or convert to integer cents for comparison if that's how VNPay determines this.
	// For simplicity: if percentage deduction is 0, assume full refund of the final_amount.
	// VNPay "02" for full refund of original transaction, "03" for partial.
	// If the calculated refundAmount equals the invoice.FinalAmount, it's a full refund of that final_amount.
	if req.PercentageDeduction == 0.0 && refundAmount == invoice.FinalAmount {
		transactionType = "02" // Full refund of the final amount
	} else {
		transactionType = "03" // Partial refund
	}
	if refundAmount == 0 && invoice.FinalAmount > 0 { // If refunding zero for a non-zero payment
		transactionType = "03" // Still a type of partial refund (refunding zero amount)
	}

	// 4. Get TransactionDate (original payment date in YYYYMMDD)
	var transactionDateStr string
	if invoice.VnpayPayDate.Valid && len(invoice.VnpayPayDate.String) >= 8 {
		transactionDateStr = invoice.VnpayPayDate.String[:8] // YYYYMMDD from YYYYMMDDHHMMSS
	} else if invoice.IssueDate.Valid {
		transactionDateStr = invoice.IssueDate.Time.Format("20060102")
	} else {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Cannot determine original transaction date for refund", nil)
		return
	}

	// 5. Prepare VNPayRefundRequest for the service
	vnpayRefundServiceReq := model.VNPayRefundRequest{
		TxnRef:          invoice.VnpayTxnRef.String,
		TransactionDate: transactionDateStr,
		Amount:          refundAmount,
		TransactionType: transactionType,
		CreateBy:        req.RefundInitiator,
	}

	// 6. Call VNPayService to get refund request data (service will also update DB optimistically or based on IPN)
	// The service's RefundTransaction now also takes the 'reason'.
	refundDataToVNPay, err := c.vnpaySvc.RefundTransaction(ctx, vnpayRefundServiceReq, clientIP, req.Reason)
	if err != nil {
		// Check if it's a "already refunded" or "not completed" type error from the service if it pre-checks
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to prepare/process VNPay refund data", err.Error())
		return
	}

	// The refundDataToVNPay is what should be POSTed to VNPay.
	// The controller or a dedicated client would make this HTTP POST.
	// For this example, we return the prepared data.
	// The invoice status update to REFUNDED is now handled within vnpaySvc.RefundTransaction (or should be triggered by IPN).
	utils.RespondWithSuccess(ctx, http.StatusOK, "VNPay refund data prepared. POST this to VNPay. Invoice status updated/pending.", gin.H{
		"prepared_refund_data": refundDataToVNPay,
		"vnpay_refund_api_url": c.cfg.TransactionAPI, // Ensure this is the correct refund API URL
		"invoice_id":           invoice.InvoiceID,
		"refunded_amount":      refundAmount,
	})
}

func (c *VNPayController) FailInvoice(ctx *gin.Context) {
	var req model.TicketStatusFailureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload, 'reason' is required", err.Error())
		return
	}

	// Lấy thông tin hóa đơn để biết payment method
	invoice, err := c.invoiceSvc.GetInvoiceByID(ctx, req.Invoice)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "Invoice not found", err.Error())
		return
	}

	// Gọi hàm dùng chung để hủy hóa đơn
	updatedInvoice, err := c.invoiceSvc.UpdateInvoiceStatusForPaymentFailureForUUID(
		ctx,
		invoice.InvoiceID, // Dùng ID thay vì identifier
		model.PaymentMethod(invoice.PaymentMethod.String),
		req.Reason,
	)

	if err != nil {
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Failed to fail invoice", err.Error())
		return
	}

	apiInvoice := c.invoiceSvc.MapDbInvoiceToAPIResponse(updatedInvoice)
	utils.RespondWithSuccess(ctx, http.StatusOK, "Invoice marked as failed successfully", apiInvoice)
}
