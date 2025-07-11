package models

type VNPayPaymentResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PaymentURL string `json:"payment_url"`
		TxnRef     int    `json:"txn_ref"`
	} `json:"data"`
	InvoiceID string `json:"invoice_id"`
}
