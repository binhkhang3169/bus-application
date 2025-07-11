package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending              PaymentStatus = "PENDING"
	PaymentStatusCompleted            PaymentStatus = "COMPLETED"
	PaymentStatusFailed               PaymentStatus = "FAILED"
	PaymentStatusRefunded             PaymentStatus = "REFUNDED"
	PaymentStatusAwaitingConfirmation PaymentStatus = "AWAITING_CONFIRMATION" // For bank transfers
	PaymentStatusRequiresAction       PaymentStatus = "REQUIRES_ACTION"       // From Stripe
)

const (
	PaymentMethodStaffCash     PaymentMethod = "STAFF_CASH"
	PaymentMethodStaffCard     PaymentMethod = "STAFF_POS_CARD"
	PaymentMethodStaffTransfer PaymentMethod = "STAFF_TRANSFER" // If staff confirms an immediate transfer
	PaymentMethodStaffOther    PaymentMethod = "STAFF_OTHER"
)

// PaymentMethod định nghĩa các phương thức thanh toán
type PaymentMethod string

const (
	PaymentMethodVNPay  PaymentMethod = "VNPAY"
	PaymentMethodStripe PaymentMethod = "STRIPE"
	PaymentMethodBank   PaymentMethod = "BANK"
)

type StaffDirectPaymentRequest struct {
	CustomerID           string  `json:"customer_id" binding:"required"`
	TicketID             string  `json:"ticket_id" binding:"required"`
	Amount               float64 `json:"amount" binding:"required,gt=0"`
	Currency             string  `json:"currency" binding:"required"` // e.g., "VND", "USD"
	PaymentMethod        string  `json:"payment_method" binding:"required,oneof=STAFF_CASH STAFF_POS_CARD STAFF_TRANSFER STAFF_OTHER"`
	TransactionReference string  `json:"transaction_reference,omitempty"` // e.g., POS transaction ID, Cheque number, staff note
	StaffID              string  `json:"staff_id" binding:"required"`     // ID of the staff member processing the payment
	Notes                string  `json:"notes,omitempty"`                 // Additional notes
	InvoiceType          string  `json:"invoice_type,omitempty"`
	DiscountAmount       float64 `json:"discount_amount,omitempty"`
	TaxAmount            float64 `json:"tax_amount,omitempty"`
}

// VNPayPaymentRequest là request body cho việc tạo thanh toán VNPay
type VNPayPaymentRequest struct {
	Amount         float64 `json:"amount" binding:"required"`
	BankCode       string  `json:"bank_code"` // Tùy chọn: Cho chuyển hướng ngân hàng cụ thể
	Language       string  `json:"language" binding:"required,oneof=vn en"`
	InvoiceType    string  `json:"invoice_type"` // Ví dụ: "flight_ticket", "service_fee"
	CustomerID     string  `json:"customer_id" binding:"required"`
	TicketID       string  `json:"ticket_id" binding:"required"` // Liên kết đến vé đang được thanh toán
	DiscountAmount float64 `json:"discount_amount"`
	TaxAmount      float64 `json:"tax_amount"`
	Notes          string  `json:"notes"` // Ghi chú thêm cho hóa đơn
}

// VNPayPaymentResponse là response trả về sau khi tạo yêu cầu thanh toán VNPay
type VNPayPaymentResponse struct {
	PaymentURL string    `json:"payment_url"`
	TxnRef     string    `json:"txn_ref"`    // VNPay TxnRef (thường là chuỗi số)
	InvoiceID  uuid.UUID `json:"invoice_id"` // ID hóa đơn nội bộ của bạn
}

// VNPayReturnResponse là response khi VNPay redirect về Return URL
type VNPayReturnResponse struct {
	IsValid        bool      `json:"is_valid_signature"`
	TransactionNo  string    `json:"vnp_transaction_no"` // Mã giao dịch của VNPay
	Amount         float64   `json:"amount"`             // Số tiền đã thanh toán
	OrderInfo      string    `json:"order_info"`
	ResponseCode   string    `json:"vnp_response_code"` // Mã phản hồi của VNPay
	BankCode       string    `json:"vnp_bank_code"`     // Mã ngân hàng
	PaymentTime    string    `json:"vnp_pay_date"`      // Thời gian thanh toán (YYYYMMDDHHMMSS)
	TransactionRef string    `json:"vnp_txn_ref"`       // Mã tham chiếu giao dịch của bạn
	Message        string    `json:"message"`           // Thông báo kết quả (ví dụ: "Thanh toán thành công")
	InvoiceID      uuid.UUID `json:"invoice_id"`        // ID hóa đơn nội bộ của bạn
}

// VNPayIPNResponse là struct cho VNPay IPN response mà service của bạn sẽ trả về cho VNPay
type VNPayIPNResponse struct {
	RspCode string `json:"RspCode"` // Mã phản hồi của VNPay (ví dụ: "00")
	Message string `json:"Message"` // Thông báo của VNPay (ví dụ: "Confirm Success")
}

// InitialStripePaymentRequest là request từ frontend để tạo PaymentIntent với Stripe
type InitialStripePaymentRequest struct {
	Amount         int64  `json:"amount" binding:"required"` // Số tiền ở đơn vị nhỏ nhất (ví dụ: cents cho USD, hoặc đơn vị cơ sở cho VND)
	Currency       string `json:"currency" binding:"required,oneof=usd vnd"`
	InvoiceType    string `json:"invoice_type,omitempty"`
	CustomerID     string `json:"customer_id" binding:"required"`
	TicketID       string `json:"ticket_id" binding:"required"`
	DiscountAmount int64  `json:"discount_amount,omitempty"` // Để ghi nhận, số tiền thực tế gửi cho Stripe là final_amount
	TaxAmount      int64  `json:"tax_amount,omitempty"`      // Để ghi nhận
	Notes          string `json:"notes,omitempty"`           // Ghi chú thêm cho hóa đơn
}

// StripePaymentIntentResponse trả về cho frontend sau khi tạo PaymentIntent
type StripePaymentIntentResponse struct {
	ClientSecret    string    `json:"client_secret"`
	PaymentIntentID string    `json:"payment_intent_id"`
	InvoiceID       uuid.UUID `json:"invoice_id"` // ID hóa đơn nội bộ của bạn
	PublishableKey  string    `json:"publishable_key"`
}

// StripePaymentConfirmationRequest gửi từ frontend sau khi Stripe.js xác nhận thanh toán phía client
type StripePaymentConfirmationRequest struct {
	PaymentIntentID string `json:"payment_intent_id" binding:"required"`
}

// StripeWebhookEvent là struct chung cho webhook đến (bạn sẽ parse event.Data.Object cụ thể)
type StripeWebhookEvent struct {
	ID              string                 `json:"id"`
	APIVersion      string                 `json:"api_version"`
	Data            StripeWebhookEventData `json:"data"`
	Type            string                 `json:"type"` // Ví dụ: "payment_intent.succeeded"
	Created         int64                  `json:"created"`
	Livemode        bool                   `json:"livemode"`
	PendingWebhooks int                    `json:"pending_webhooks"`
	Request         struct {
		ID             string `json:"id"`
		IdempotencyKey string `json:"idempotency_key"`
	} `json:"request"`
	Object string `json:"object"` // "event"
}

// StripeWebhookEventData chứa object Stripe từ webhook
type StripeWebhookEventData struct {
	Object map[string]interface{} `json:"object"` // Đây sẽ là object Stripe, ví dụ PaymentIntent
}

// StandardAPIResponse là một wrapper response chung cho API
type StandardAPIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError định nghĩa cấu trúc lỗi cho API response
type APIError struct {
	Code    int         `json:"code,omitempty"` // Mã lỗi HTTP hoặc mã lỗi nội bộ
	Details interface{} `json:"details,omitempty"`
}

// GetInvoiceResponse là response khi lấy thông tin một hóa đơn
type GetInvoiceResponse struct {
	InvoiceID                  uuid.UUID `json:"invoice_id"`
	InvoiceNumber              string    `json:"invoice_number"`
	InvoiceType                string    `json:"invoice_type,omitempty"`
	CustomerID                 string    `json:"customer_id"`
	TicketID                   string    `json:"ticket_id"`
	TotalAmount                float64   `json:"total_amount"`       // Hiển thị dạng float cho dễ đọc
	DiscountAmount             float64   `json:"discount_amount"`    // Hiển thị dạng float
	TaxAmount                  float64   `json:"tax_amount"`         // Hiển thị dạng float
	FinalAmount                float64   `json:"final_amount"`       // Hiển thị dạng float
	Currency                   string    `json:"currency,omitempty"` // Added currency
	PaymentStatus              string    `json:"payment_status"`
	PaymentMethod              string    `json:"payment_method,omitempty"`
	IssueDate                  string    `json:"issue_date,omitempty"` // Format YYYY-MM-DD HH:MM:SS
	Notes                      string    `json:"notes,omitempty"`
	CreatedAt                  string    `json:"created_at"`
	UpdatedAt                  string    `json:"updated_at"`
	VNPayTxnRef                string    `json:"vnpay_txn_ref,omitempty"`
	VNPayBankCode              string    `json:"vnpay_bank_code,omitempty"`
	VNPayTxnNo                 string    `json:"vnpay_txn_no,omitempty"`
	VNPayPayDate               string    `json:"vnpay_pay_date,omitempty"`
	StripePaymentIntentID      string    `json:"stripe_payment_intent_id,omitempty"`
	StripeChargeID             string    `json:"stripe_charge_id,omitempty"`
	StripeCustomerID           string    `json:"stripe_customer_id,omitempty"`
	StripePaymentMethodDetails string    `json:"stripe_payment_method_details,omitempty"`
}

// VNPayQueryRequest cho VNPay QueryDR
type VNPayQueryRequest struct {
	TxnRef          string `json:"txn_ref" binding:"required"`          // Mã giao dịch của merchant (vnp_TxnRef)
	TransactionDate string `json:"transaction_date" binding:"required"` // Ngày giao dịch (YYYYMMDD)
}

// VNPayRefundRequest cho VNPay Refund (sử dụng bởi service)
type VNPayRefundRequest struct {
	TxnRef          string  `json:"txn_ref" binding:"required"`          // Mã giao_dịch cần hoàn tiền (vnp_TxnRef)
	TransactionDate string  `json:"transaction_date" binding:"required"` // Ngày giao dịch (YYYYMMDD)
	Amount          float64 `json:"amount" binding:"required"`           // Số tiền cần hoàn
	TransactionType string  `json:"transaction_type" binding:"required"` // Loại hoàn tiền "02": Toàn phần, "03": Một phần
	CreateBy        string  `json:"create_by" binding:"required"`        // Người yêu cầu hoàn tiền
}

// VNPayInitiateRefundRequest là request body cho việc khởi tạo hoàn tiền VNPay từ controller
type VNPayInitiateRefundRequest struct {
	TicketID            string  `json:"ticket_id" binding:"required"`
	PercentageDeduction float64 `json:"percentage_deduction" binding:"min=0,max=1"` // Ví dụ: 0.0 cho không trừ, 0.1 cho trừ 10%
	Reason              string  `json:"reason" binding:"required"`
	RefundInitiator     string  `json:"refund_initiator" binding:"required"` // Người/hệ thống yêu cầu hoàn tiền
}

// StripeInitiateRefundRequest là request body cho việc khởi tạo hoàn tiền Stripe từ controller
type StripeInitiateRefundRequest struct {
	TicketID            string  `json:"ticket_id" binding:"required"`
	PercentageDeduction float64 `json:"percentage_deduction" binding:"min=0,max=1"` // Ví dụ: 0.0 cho không trừ, 0.1 cho trừ 10%
	Reason              string  `json:"reason" binding:"required"`
	// RefundInitiator có thể thêm nếu cần cho audit log hoặc metadata Stripe
}

// PaymentMethod represents the type of payment.

// InitialBankPaymentRequest is used to initiate a bank payment.
type InitialBankPaymentRequest struct {
	CustomerID  string  `json:"customer_id" binding:"required"`
	TicketID    string  `json:"ticket_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required"` // e.g., "VND", "USD"
	InvoiceType string  `json:"invoice_type,omitempty"`
	// Add any other fields needed to create the initial invoice
}

// BankPaymentDetailsResponse contains the bank account details your customers will use to pay.
type BankPaymentDetailsResponse struct {
	InvoiceID                    uuid.UUID `json:"invoice_id"`
	BankTransferCode             string    `json:"bank_transfer_code"` // Unique code for user to put in reference
	PayableAmount                float64   `json:"payable_amount"`
	Currency                     string    `json:"currency"`
	OurBankAccountName           string    `json:"our_bank_account_name"`
	OurBankAccountNumber         string    `json:"our_bank_account_number"`
	OurBankName                  string    `json:"our_bank_name"`
	PaymentReferenceInstructions string    `json:"payment_reference_instructions"`
	DueDate                      time.Time `json:"due_date,omitempty"`
	// You might include a QR code string for payment here (e.g., VietQR)
}

// BankPaymentConfirmationRequest is used by an admin or an automated system to confirm a bank payment.
type BankPaymentConfirmationRequest struct {
	InvoiceID             uuid.UUID `json:"invoice_id,omitempty"`         // Confirm by InvoiceID
	BankTransferCode      string    `json:"bank_transfer_code,omitempty"` // Or confirm by BankTransferCode
	PayerAccountName      string    `json:"payer_account_name,omitempty"`
	PayerAccountNumber    string    `json:"payer_account_number,omitempty"`
	PayerBankName         string    `json:"payer_bank_name,omitempty"`
	BankTransactionID     string    `json:"bank_transaction_id,omitempty"` // The bank's reference for the received payment
	AmountReceived        float64   `json:"amount_received" binding:"required"`
	CurrencyReceived      string    `json:"currency_received" binding:"required"`
	ConfirmationTimestamp time.Time `json:"confirmation_timestamp" binding:"required"`
	ConfirmationDetails   string    `json:"confirmation_details,omitempty"` // e.g., Screenshot URL, notes
	ConfirmedBy           string    `json:"confirmed_by,omitempty"`         // User ID of admin or system
}

// BankRefundRequest might be needed later
// type BankRefundRequest struct {
// 	InvoiceID string `json:"invoice_id" binding:"required"`
// 	Amount    float64 `json:"amount,omitempty"` // Optional: if partial refund, otherwise full
// 	Reason    string `json:"reason" binding:"required"`
// }

const (
	TicketStatusPaid     = "1" // Trạng thái vé: Đã thanh toán
	TicketStatusRefunded = "2" // Trạng thái vé: Đã hoàn tiền
	TicketStatusFailed   = "3" // Trạng thái vé: Thanh toán thất bại/Đã hủy
)

type TicketStatusFailureRequest struct {
	Invoice uuid.UUID `json:"invoice_id" binding:"required"`
	Reason  string    `json:"reason" binding:"required"`
}
