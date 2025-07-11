package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"payment_service/config"
	"payment_service/domain/model"
	// Assuming this is updated after sqlc generate
)

// VNPayService handles the VNPay payment integration
type VNPayService struct {
	config     *config.VNPayConfig
	invoiceSvc InvoiceServiceInterface // Changed to InvoiceServiceInterface
}

// NewVNPayService creates a new VNPay service
func NewVNPayService(cfg *config.VNPayConfig, invoiceSvc InvoiceServiceInterface) *VNPayService { // Changed to InvoiceServiceInterface
	return &VNPayService{
		config:     cfg,
		invoiceSvc: invoiceSvc,
	}
}

// CreatePayment creates a new payment URL for VNPay
func (s *VNPayService) CreatePayment(ctx context.Context, req model.VNPayPaymentRequest) (*model.VNPayPaymentResponse, error) {
	// Initialize random seed (deprecated, use rand.New(rand.NewSource(time.Now().UnixNano())) for new code)
	// For simplicity in this context, keeping rand.Seed
	rand.Seed(time.Now().UnixNano())

	txnRef := strconv.Itoa(rand.Intn(9999999) + 1000000) // Generate a unique enough TxnRef

	invoice, err := s.invoiceSvc.CreateInvoiceForVNPay(ctx, req, txnRef)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	now := time.Now()
	createDate := now.Format("20060102150405")
	expireTime := now.Add(15 * time.Minute).Format("20060102150405")

	// VNPay expects amount in VND (integer, multiplied by 100 if it were cents, but it's base unit for VND)
	// req.Amount is float64, VNPay expects integer string for vnp_Amount
	amountVND := int(req.Amount * 100) // This is standard for VNPay, treating amount as base unit * 100

	inputData := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_TmnCode":    s.config.TmnCode,
		"vnp_Amount":     strconv.Itoa(amountVND),
		"vnp_Command":    "pay",
		"vnp_CreateDate": createDate,
		"vnp_CurrCode":   "VND",
		"vnp_IpAddr":     "127.0.0.1", // This should ideally be passed from the controller or obtained from request context
		"vnp_Locale":     req.Language,
		"vnp_OrderInfo":  fmt.Sprintf("Thanh toan cho ve %s, hoa don %s", req.TicketID, invoice.InvoiceNumber),
		"vnp_OrderType":  "other", // Or map req.InvoiceType
		"vnp_ReturnUrl":  s.config.ReturnURL,
		"vnp_TxnRef":     txnRef,
		"vnp_ExpireDate": expireTime,
	}

	if req.BankCode != "" {
		inputData["vnp_BankCode"] = req.BankCode
	}

	keys := make([]string, 0, len(inputData))
	for k := range inputData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryBuilder strings.Builder
	var hashDataBuilder strings.Builder

	for i, k := range keys {
		encodedKey := url.QueryEscape(k)
		encodedValue := url.QueryEscape(inputData[k])

		if i > 0 {
			queryBuilder.WriteString("&")
			hashDataBuilder.WriteString("&")
		}
		queryBuilder.WriteString(encodedKey)
		queryBuilder.WriteString("=")
		queryBuilder.WriteString(encodedValue)

		hashDataBuilder.WriteString(encodedKey) // hashData uses raw key as per VNPay docs for v2.1.0 example
		hashDataBuilder.WriteString("=")
		hashDataBuilder.WriteString(encodedValue) // raw value as per VNPay docs for v2.1.0 example
	}

	vnpURL := s.config.VNPayURL + "?" + queryBuilder.String()
	hashData := hashDataBuilder.String()

	hmacObj := hmac.New(sha512.New, []byte(s.config.HashSecret))
	hmacObj.Write([]byte(hashData))
	vnpSecureHash := hex.EncodeToString(hmacObj.Sum(nil))

	vnpURL += "&vnp_SecureHash=" + vnpSecureHash

	response := &model.VNPayPaymentResponse{
		PaymentURL: vnpURL,
		TxnRef:     txnRef, // Return the generated TxnRef
		InvoiceID:  invoice.InvoiceID,
	}

	return response, nil
}

// ProcessReturn processes the return from VNPay payment gateway
func (s *VNPayService) ProcessReturn(ctx context.Context, queryParams url.Values) (*model.VNPayReturnResponse, error) {
	// Get the secure hash from the query
	vnpSecureHash := queryParams.Get("vnp_SecureHash")

	// Create a map to store all "vnp_" parameters
	inputData := make(map[string]string)
	for key, values := range queryParams {
		if strings.HasPrefix(key, "vnp_") && key != "vnp_SecureHash" {
			inputData[key] = values[0]
		}
	}

	// Sort the keys
	keys := make([]string, 0, len(inputData))
	for k := range inputData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the hash data
	var hashDataBuilder strings.Builder
	for i, k := range keys {
		if i > 0 {
			hashDataBuilder.WriteString("&")
		}
		hashDataBuilder.WriteString(url.QueryEscape(k))
		hashDataBuilder.WriteString("=")
		hashDataBuilder.WriteString(url.QueryEscape(inputData[k]))
	}

	// Calculate secure hash
	hashData := hashDataBuilder.String()
	hmacObj := hmac.New(sha512.New, []byte(s.config.HashSecret))
	hmacObj.Write([]byte(hashData))
	secureHash := hex.EncodeToString(hmacObj.Sum(nil))

	// Verify the secure hash
	isValidSignature := secureHash == vnpSecureHash

	amountStr := queryParams.Get("vnp_Amount")
	amountInt, _ := strconv.Atoi(amountStr)
	amount := float64(amountInt) / 100.0 // Amount from VNPay is base unit * 100

	txnRef := queryParams.Get("vnp_TxnRef")
	responseCode := queryParams.Get("vnp_ResponseCode")
	orderInfo := queryParams.Get("vnp_OrderInfo")
	bankCode := queryParams.Get("vnp_BankCode")
	payDate := queryParams.Get("vnp_PayDate")
	transactionNo := queryParams.Get("vnp_TransactionNo")

	invoice, err := s.invoiceSvc.GetInvoiceByVNPayTxnRef(ctx, txnRef)
	if err != nil {
		// If invoice not found, it's a critical issue or invalid txnRef
		log.Printf("Error: ProcessReturn: VNPay TxnRef %s not found in DB: %v", txnRef, err)
		// Potentially return an error or a response indicating failure without invoice linkage
		return &model.VNPayReturnResponse{
			IsValid:        isValidSignature, // Signature check might still be valid
			TransactionNo:  transactionNo,
			Amount:         amount,
			OrderInfo:      orderInfo,
			ResponseCode:   responseCode,
			BankCode:       bankCode,
			PaymentTime:    payDate,
			TransactionRef: txnRef,
			Message:        "Transaction reference not found or error fetching invoice.",
		}, fmt.Errorf("invoice for TxnRef %s not found: %w", txnRef, err)
	}

	var resultMessage string
	if isValidSignature {
		if responseCode == "00" { // Payment successful
			// Check if invoice is already completed to prevent reprocessing
			if invoice.PaymentStatus.String == string(model.PaymentStatusPending) {
				_, updateErr := s.invoiceSvc.UpdateInvoiceStatusForVNPaySuccess(ctx, txnRef, bankCode, transactionNo, payDate)
				if updateErr != nil {
					log.Printf("Error: ProcessReturn: Failed to update invoice for VNPay success (TxnRef: %s): %v", txnRef, updateErr)
					resultMessage = "Payment successful, but internal update failed. Please contact support."
					// Return error as this is an internal processing failure
					return nil, fmt.Errorf("failed to update invoice for VNPay success (TxnRef: %s): %w", txnRef, updateErr)
				}
				resultMessage = "Payment successful"
			} else {
				resultMessage = fmt.Sprintf("Payment already processed (Status: %s)", invoice.PaymentStatus.String)
			}
		} else { // Payment failed or canceled by user
			failureReason := fmt.Sprintf("VNPay Response Code: %s. Message: %s", responseCode, vnPayMessage(responseCode))
			if invoice.PaymentStatus.String == string(model.PaymentStatusPending) {
				_, updateErr := s.invoiceSvc.UpdateInvoiceStatusForPaymentFailure(ctx, txnRef, model.PaymentMethodVNPay, failureReason)
				if updateErr != nil {
					log.Printf("Error: ProcessReturn: Failed to update invoice for VNPay failure (TxnRef: %s): %v", txnRef, updateErr)
				}
			}
			resultMessage = fmt.Sprintf("Payment failed. Reason: %s", failureReason)
		}
	} else {
		resultMessage = "Invalid signature. Payment integrity compromised."
		log.Printf("Error: ProcessReturn: Invalid signature for VNPay TxnRef %s", txnRef)
		// Do not update DB if signature is invalid
	}

	return &model.VNPayReturnResponse{
		IsValid:        isValidSignature,
		TransactionNo:  transactionNo,
		Amount:         amount,
		OrderInfo:      orderInfo,
		ResponseCode:   responseCode,
		BankCode:       bankCode,
		PaymentTime:    payDate,
		TransactionRef: txnRef,
		Message:        resultMessage,
		InvoiceID:      invoice.InvoiceID,
	}, nil
}

// ProcessIPN processes the Instant Payment Notification from VNPay
func (s *VNPayService) ProcessIPN(ctx context.Context, queryParams url.Values) (*model.VNPayIPNResponse, error) {
	vnpSecureHash := queryParams.Get("vnp_SecureHash")
	inputData := make(map[string]string)
	for key, values := range queryParams {
		if strings.HasPrefix(key, "vnp_") && key != "vnp_SecureHash" {
			inputData[key] = values[0]
		}
	}

	keys := make([]string, 0, len(inputData))
	for k := range inputData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var hashDataBuilder strings.Builder
	for i, k := range keys {
		if i > 0 {
			hashDataBuilder.WriteString("&")
		}
		hashDataBuilder.WriteString(k) // Raw key
		hashDataBuilder.WriteString("=")
		hashDataBuilder.WriteString(inputData[k]) // Raw value
	}
	hashData := hashDataBuilder.String()

	hmacObj := hmac.New(sha512.New, []byte(s.config.HashSecret))
	hmacObj.Write([]byte(hashData))
	calculatedSecureHash := hex.EncodeToString(hmacObj.Sum(nil))

	ipnResponse := &model.VNPayIPNResponse{}

	if calculatedSecureHash != vnpSecureHash {
		ipnResponse.RspCode = "97" // Invalid Signature
		ipnResponse.Message = "Invalid Signature"
		log.Printf("IPN Error: Invalid signature. TxnRef: %s", queryParams.Get("vnp_TxnRef"))
		return ipnResponse, nil // No further processing, return 200 OK with error code to VNPay
	}

	txnRef := queryParams.Get("vnp_TxnRef")
	vnpAmountStr := queryParams.Get("vnp_Amount")
	responseCode := queryParams.Get("vnp_ResponseCode") // VNPay uses vnp_ResponseCode in IPN as well
	// transactionStatus := queryParams.Get("vnp_TransactionStatus") // Usually 00 for successful transaction in IPN

	invoice, err := s.invoiceSvc.GetInvoiceByVNPayTxnRef(ctx, txnRef)
	if err != nil {
		ipnResponse.RspCode = "01" // Order not found
		ipnResponse.Message = "Order not found"
		log.Printf("IPN Error: Order not found for TxnRef %s: %v", txnRef, err)
		return ipnResponse, nil
	}

	// Verify amount (VNPay amount is base unit * 100)
	vnpAmountInt, convErr := strconv.Atoi(vnpAmountStr)
	if convErr != nil {
		ipnResponse.RspCode = "04" // Invalid amount (format error)
		ipnResponse.Message = "Invalid amount format"
		log.Printf("IPN Error: Invalid amount format for TxnRef %s: %s", txnRef, vnpAmountStr)
		return ipnResponse, nil
	}

	// Compare with invoice.FinalAmount (which is float64 in base currency unit)
	// VNPay sends amount as integer (base unit * 100), so convert invoice.FinalAmount to integer in same format
	expectedAmountVND := int(invoice.FinalAmount * 100)
	if vnpAmountInt != expectedAmountVND {
		ipnResponse.RspCode = "04" // Invalid amount
		ipnResponse.Message = "Invalid amount"
		log.Printf("IPN Error: Amount mismatch for TxnRef %s. Expected: %d, Got: %d", txnRef, expectedAmountVND, vnpAmountInt)
		return ipnResponse, nil
	}

	// Check if invoice is already processed (COMPLETED or FAILED)
	if invoice.PaymentStatus.String == string(model.PaymentStatusCompleted) || invoice.PaymentStatus.String == string(model.PaymentStatusFailed) {
		// If already COMPLETED, VNPay expects 00. If FAILED and IPN says success, this is an issue.
		// For simplicity, if already completed, assume it's a duplicate IPN.
		if invoice.PaymentStatus.String == string(model.PaymentStatusCompleted) && responseCode == "00" {
			ipnResponse.RspCode = "00" // Already confirmed, acknowledge success
			ipnResponse.Message = "Order already confirmed"
		} else {
			ipnResponse.RspCode = "02" // Order already confirmed (generic)
			ipnResponse.Message = "Order already confirmed with different status or IPN indicates failure for already completed order"
		}
		log.Printf("IPN Info: Order TxnRef %s already processed. Current status: %s. IPN ResponseCode: %s", txnRef, invoice.PaymentStatus.String, responseCode)
		return ipnResponse, nil
	}

	// At this point, signature is valid, order exists, amount matches, and order is PENDING.
	if responseCode == "00" { // Payment successful
		bankCode := queryParams.Get("vnp_BankCode")
		payDate := queryParams.Get("vnp_PayDate")
		transactionNo := queryParams.Get("vnp_TransactionNo")

		_, updateErr := s.invoiceSvc.UpdateInvoiceStatusForVNPaySuccess(ctx, txnRef, bankCode, transactionNo, payDate)
		if updateErr != nil {
			ipnResponse.RspCode = "99" // Internal error, VNPay will retry
			ipnResponse.Message = "System error while updating order"
			log.Printf("IPN Error: Failed to update invoice for VNPay success (TxnRef: %s): %v", txnRef, updateErr)
			// Return actual error to controller to decide HTTP status for VNPay (usually 200 OK with specific RspCode)
			return ipnResponse, fmt.Errorf("IPN: failed to update invoice for TxnRef %s: %w", txnRef, updateErr)
		}
		ipnResponse.RspCode = "00"
		ipnResponse.Message = "Confirm Success"
		log.Printf("IPN Success: Payment confirmed for TxnRef %s. Invoice updated.", txnRef)
	} else { // Payment failed
		failureReason := fmt.Sprintf("VNPay IPN Response Code: %s. Message: %s", responseCode, vnPayMessage(responseCode))
		_, updateErr := s.invoiceSvc.UpdateInvoiceStatusForPaymentFailure(ctx, txnRef, model.PaymentMethodVNPay, failureReason)
		if updateErr != nil {
			ipnResponse.RspCode = "99"
			ipnResponse.Message = "System error while updating failed order"
			log.Printf("IPN Error: Failed to update invoice for VNPay failure (TxnRef: %s): %v", txnRef, updateErr)
			return ipnResponse, fmt.Errorf("IPN: failed to update failed invoice for TxnRef %s: %w", txnRef, updateErr)
		}
		// For failed payments, VNPay might not expect "00".
		// Responding with the failure code received or a generic failure acknowledgement.
		// However, VNPay IPN documentation implies they expect 00/Message for successful processing of the IPN itself,
		// not necessarily for the transaction outcome. Let's send 00 to acknowledge IPN receipt and processing.
		ipnResponse.RspCode = "00" // Acknowledge IPN processing was successful
		ipnResponse.Message = "Confirm Success (transaction failed, status updated)"
		log.Printf("IPN Info: Payment failed for TxnRef %s. Reason: %s. Invoice updated to FAILED.", txnRef, failureReason)
	}

	return ipnResponse, nil
}

// QueryTransaction prepares data for querying a transaction status with VNPay
func (s *VNPayService) QueryTransaction(ctx context.Context, req model.VNPayQueryRequest, ipAddr string) (map[string]string, error) {
	rand.Seed(time.Now().UnixNano())                          // Consider a better random source for production
	requestID := strconv.FormatInt(time.Now().UnixNano(), 10) // Unique request ID
	createDate := time.Now().Format("20060102150405")

	// Ensure TxnRef and TransactionDate are from the request
	if req.TxnRef == "" || req.TransactionDate == "" {
		return nil, fmt.Errorf("TxnRef and TransactionDate are required for VNPay query")
	}

	dataRequest := map[string]string{
		"vnp_RequestId":       requestID,
		"vnp_Version":         "2.1.0",
		"vnp_Command":         "querydr",
		"vnp_TmnCode":         s.config.TmnCode,
		"vnp_TxnRef":          req.TxnRef, // Merchant's transaction reference
		"vnp_OrderInfo":       fmt.Sprintf("Query transaction status for %s", req.TxnRef),
		"vnp_TransactionDate": req.TransactionDate, // Original transaction date YYYYMMDD
		"vnp_CreateDate":      createDate,          // Current request creation date
		"vnp_IpAddr":          ipAddr,
	}

	// VNPay's documented hash string for querydr:
	// vnp_RequestId|vnp_Version|vnp_Command|vnp_TmnCode|vnp_TxnRef|vnp_TransactionDate|vnp_CreateDate|vnp_IpAddr|vnp_OrderInfo
	// Note: The order matters and must match VNPay's specification.
	// Some VNPay docs show a slightly different field order or set for hashing for different commands. Always verify.
	// For querydr, often it's: RequestId|Version|Command|TmnCode|TxnRef|TransactionDate|CreateDate|IpAddr|OrderInfo
	// Re-checking typical VNPay pattern:
	// It's usually | separated values of specific fields IN ORDER.
	// Let's assume the fields must be sorted by key name first, then concatenated for the hash.
	// However, VNPay docs sometimes specify an explicit order for hash string.
	// Sticking to sorted keys as a general approach if explicit order is not strictly clear or varies by API version.
	// *Correction*: VNPay often specifies the exact string format for hashing, not sorted keys.
	// For queryDR, the typical hash string format from their PHP SDK example is:
	// $vnp_RequestId."|".$vnp_Version."|".$vnp_Command."|".$vnp_TmnCode."|".$vnp_TxnRef."|".$vnp_TransactionDate."|".$vnp_CreateDate."|".$vnp_IpAddr."|".$vnp_OrderInfo;
	// This order is important.

	hashDataString := strings.Join([]string{
		dataRequest["vnp_RequestId"],
		dataRequest["vnp_Version"],
		dataRequest["vnp_Command"],
		dataRequest["vnp_TmnCode"],
		dataRequest["vnp_TxnRef"],
		dataRequest["vnp_TransactionDate"],
		dataRequest["vnp_CreateDate"],
		dataRequest["vnp_IpAddr"],
		dataRequest["vnp_OrderInfo"],
	}, "|")

	hmacObj := hmac.New(sha512.New, []byte(s.config.HashSecret))
	hmacObj.Write([]byte(hashDataString))
	secureHash := hex.EncodeToString(hmacObj.Sum(nil))
	dataRequest["vnp_SecureHash"] = secureHash

	return dataRequest, nil
}

// RefundTransaction prepares data for refunding a transaction.
// The actual API call to VNPay is expected to be made by the caller with this data.
// The 'reason' parameter is for internal logging/updating invoice notes.
func (s *VNPayService) RefundTransaction(ctx context.Context, req model.VNPayRefundRequest, ipAddr string, reason string) (map[string]string, error) {
	rand.Seed(time.Now().UnixNano())
	requestID := strconv.FormatInt(time.Now().UnixNano(), 10)
	createDate := time.Now().Format("20060102150405")

	// VNPay expects amount in integer format (base unit * 100)
	amountToRefundVND := int(req.Amount * 100)

	refundData := map[string]string{
		"vnp_RequestId":       requestID,
		"vnp_Version":         "2.1.0",
		"vnp_Command":         "refund",
		"vnp_TmnCode":         s.config.TmnCode,
		"vnp_TransactionType": req.TransactionType, // "02": full, "03": partial
		"vnp_TxnRef":          req.TxnRef,          // Original merchant transaction ref
		"vnp_Amount":          strconv.Itoa(amountToRefundVND),
		"vnp_OrderInfo":       fmt.Sprintf("Hoan tien cho giao dich %s. Ly do: %s", req.TxnRef, reason),
		"vnp_TransactionNo":   "0",                 // Use "0" if original vnp_TransactionNo is unknown or not applicable
		"vnp_TransactionDate": req.TransactionDate, // Original payment date YYYYMMDD
		"vnp_CreateDate":      createDate,          // Refund request creation date
		"vnp_CreateBy":        req.CreateBy,        // User/system initiating refund
		"vnp_IpAddr":          ipAddr,
	}

	// Hash string for refund, order of fields is critical and per VNPay spec.
	// Example: RequestId|Version|Command|TmnCode|TransactionType|TxnRef|Amount|TransactionNo|TransactionDate|CreateBy|CreateDate|IpAddr|OrderInfo
	hashDataString := strings.Join([]string{
		refundData["vnp_RequestId"],
		refundData["vnp_Version"],
		refundData["vnp_Command"],
		refundData["vnp_TmnCode"],
		refundData["vnp_TransactionType"],
		refundData["vnp_TxnRef"],
		refundData["vnp_Amount"],
		refundData["vnp_TransactionNo"],
		refundData["vnp_TransactionDate"],
		refundData["vnp_CreateBy"],
		refundData["vnp_CreateDate"],
		refundData["vnp_IpAddr"],
		refundData["vnp_OrderInfo"],
	}, "|")

	hmacObj := hmac.New(sha512.New, []byte(s.config.HashSecret))
	hmacObj.Write([]byte(hashDataString))
	secureHash := hex.EncodeToString(hmacObj.Sum(nil))
	refundData["vnp_SecureHash"] = secureHash

	// After preparing data (and presumably after VNPay confirms refund via IPN or query):
	// Update local invoice status.
	// This direct update here might be premature if VNPay's refund is not immediate.
	// Ideally, update based on a subsequent IPN or query result for the refund.
	// For now, if the request is to *initiate* and then assume it will be processed:
	invoice, err := s.invoiceSvc.GetInvoiceByVNPayTxnRef(ctx, req.TxnRef)
	if err != nil {
		log.Printf("Error: RefundTransaction: Failed to get invoice by TxnRef %s for status update: %v", req.TxnRef, err)
		// Return prepared data, but log the failure to find invoice for update
		return refundData, fmt.Errorf("failed to retrieve invoice for status update post-refund prep: %w", err)
	}

	// Using an empty string for refundSpecificIdentifier as VNPay refund API response is not processed here directly
	_, updateErr := s.invoiceSvc.UpdateInvoiceStatusForRefund(ctx, invoice.InvoiceID, reason, "")
	if updateErr != nil {
		log.Printf("Error: RefundTransaction: Failed to update invoice %s to refunded: %v", invoice.InvoiceID, updateErr)
		// Return prepared data, but log the failure
		return refundData, fmt.Errorf("failed to update invoice status to refunded: %w", updateErr)
	}
	log.Printf("Info: RefundTransaction: Invoice %s status updated to REFUNDED (or attempted). Reason: %s", invoice.InvoiceID, reason)

	// The ticket status update is handled by InvoiceService.UpdateInvoiceStatusForRefund

	return refundData, nil // Returns data to be POSTed to VNPay
}

// updateTicketStatus is now part of InvoiceService, keep this for reference if direct call needed previously.
// For this refactoring, it's removed from vnpay_service to avoid duplication.
// func (s *VNPayService) updateTicketStatus(ctx context.Context, ticketID string, statusCode string) error { ... }

// vnPayMessage maps VNPay response codes to human-readable messages (simplified)
func vnPayMessage(responseCode string) string {
	// This is a simplified map. Refer to official VNPay documentation for all codes.
	switch responseCode {
	case "00":
		return "Giao dịch thành công"
	case "07":
		return "Trừ tiền thành công. Giao dịch bị nghi ngờ (liên hệ VNPAY)"
	case "09":
		return "Thẻ/Tài khoản chưa đăng ký Internet Banking tại Ngân hàng"
	case "10":
		return "Thẻ/Tài khoản xác thực không thành công/ nhập sai mật khẩu quá số lần quy định"
	case "11":
		return "Đã hết hạn chờ thanh toán. Xin quý khách vui lòng thực hiện lại giao dịch."
	case "12":
		return "Thẻ/Tài khoản bị khóa"
	case "13":
		return "Quý khách nhập sai mật khẩu xác thực giao dịch (OTP)"
	case "24":
		return "Giao dịch không thành công do: Khách hàng hủy giao dịch"
	case "51":
		return "Tài khoản không đủ số dư để thực hiện giao dịch"
	case "65":
		return "Giao dịch không thành công do: Tài khoản của quý khách đã vượt quá hạn mức giao dịch trong ngày"
	case "75":
		return "Ngân hàng thanh toán đang bảo trì"
	case "79":
		return "Giao dịch không thành công do: Quý khách nhập sai mật khẩu thanh toán quá số lần quy định. Xin quý khách vui lòng thực hiện lại giao dịch"
	case "99":
		return "Lỗi không xác định. Vui lòng liên hệ tổng đài VNPay."
	default:
		return fmt.Sprintf("Lỗi không xác định (Mã: %s)", responseCode)
	}
}
