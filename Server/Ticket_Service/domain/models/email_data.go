package models

// Dữ liệu cho email xác nhận vé
type TicketConfirmationData struct {
	CustomerName string   `json:"customerName"`
	TicketID     string   `json:"ticketId"`
	TripInfo     string   `json:"tripInfo"` // Có thể mở rộng thành một struct chi tiết hơn
	Price        float64  `json:"price"`
	QRCodeURLs   []string `json:"qrCodeUrls"`
}

// Dữ liệu cho email thông báo hoàn tiền
type RefundNotificationData struct {
	CustomerName string `json:"customerName"`
	TicketID     string `json:"ticketId"`
}
