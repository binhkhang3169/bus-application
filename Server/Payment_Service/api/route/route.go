package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"payment_service/api/controller" // Đảm bảo import đúng controller
)

// SetupRoutes cấu hình tất cả các API route cho ứng dụng
func SetupRoutes(
	r *gin.Engine,
	vnpayCtrl *controller.VNPayController,
	stripeCtrl *controller.StripeController,
	bankCtrl *controller.BankController,
	staffCtrl *controller.StaffAssistedPaymentController,
) {
	apiV1 := r.Group("/api/v1")

	// VNPay routes
	vnpayRoutes := apiV1.Group("/vnpay")
	{
		vnpayRoutes.POST("/create-payment", vnpayCtrl.CreatePayment)
		vnpayRoutes.GET("/return", vnpayCtrl.HandleReturn) // VNPay return URL
		vnpayRoutes.GET("/ipn", vnpayCtrl.HandleIPN)       // VNPay IPN URL
		vnpayRoutes.POST("/query-transaction", vnpayCtrl.VNPayQueryTransaction)
		vnpayRoutes.POST("/refund-transaction", vnpayCtrl.VNPayRefundTransaction) // Expects VNPayInitiateRefundRequest
	}

	// Stripe routes
	stripeRoutes := apiV1.Group("/stripe")
	{
		stripeRoutes.POST("/create-payment-intent", stripeCtrl.CreatePaymentIntent)
		stripeRoutes.POST("/confirm-payment", stripeCtrl.ConfirmStripePayment)
		stripeRoutes.POST("/webhook", stripeCtrl.HandleStripeWebhook)
		stripeRoutes.POST("/refund", stripeCtrl.RefundStripePayment) // New Stripe refund route
	}
	// Bank Transfer routes
	bankRoutes := apiV1.Group("/bank")
	{
		bankRoutes.POST("/create-payment-request", bankCtrl.CreateBankPaymentRequestHandler)
		// The confirm-payment and payment-failed routes should be protected (e.g., admin only)
		// bankRoutes.POST("/confirm-payment", authMiddleware, bankCtrl.ConfirmBankPaymentHandler)
		// bankRoutes.POST("/payment-failed", authMiddleware, bankCtrl.HandleBankPaymentFailedHandler)
		bankRoutes.POST("/confirm-payment", bankCtrl.ConfirmBankPaymentHandler)     // Add appropriate middleware
		bankRoutes.POST("/payment-failed", bankCtrl.HandleBankPaymentFailedHandler) // Add appropriate middleware

		// Add refund route if needed
		// bankRoutes.POST("/refund", authMiddleware, bankCtrl.RefundBankPaymentHandler)
	}
	// Invoice routes (can be shared or have a dedicated InvoiceController)
	// Assuming VNPayController handles these for now, or you can refactor to a new InvoiceController.
	invoiceRoutes := apiV1.Group("/invoices")
	{
		invoiceRoutes.GET("/:id", vnpayCtrl.GetInvoice) // GET /api/v1/invoices/{invoice_id}
		// This route needs auth to determine the customer.
		invoiceRoutes.GET("/customer", vnpayCtrl.GetInvoicesByCustomer) // GET /api/v1/invoices/customer (gets customer_id from token or query for demo)
		invoiceRoutes.POST("/fail", vnpayCtrl.FailInvoice)
	}
	staffPaymentRoutes := apiV1.Group("/staff-payments")
	// TODO: Apply staff-only authentication middleware to this group or specific routes
	// staffPaymentRoutes.Use(YourStaffAuthMiddleware())
	{
		staffPaymentRoutes.POST("/direct-payment", staffCtrl.HandleDirectPayment)
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
}
