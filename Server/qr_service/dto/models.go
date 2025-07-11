package dto

type TicketDetailForQR struct {
	SeatID    int32  `json:"seatId"`
	QRContent string `json:"qrContent"`
}

type OrderQRGenerationRequestEvent struct {
	OrderID       string              `json:"orderId"`
	CustomerEmail string              `json:"customerEmail"`
	CustomerName  string              `json:"customerName"`
	TotalPrice    float64             `json:"totalPrice"`
	Tickets       []TicketDetailForQR `json:"tickets"`
}

// DTO gửi tới Email Service
type TicketDetailForEmail struct {
	SeatInfo  string `json:"seatInfo"`
	QRCodeURL string `json:"qrCodeUrl"`
}

type OrderConfirmationData struct {
	CustomerName string                 `json:"customerName"`
	OrderID      string                 `json:"orderId"`
	TotalPrice   float64                `json:"totalPrice"`
	Tickets      []TicketDetailForEmail `json:"tickets"` // Danh sách vé đã có QR
}

type EmailRequestEvent struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Type  string `json:"type"`
}

// DTO cho các API response (giữ lại từ code gốc)
type QRInfoResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message,omitempty"`
	CloudinaryURL string `json:"cloudinary_url,omitempty"`
	PublicID      string `json:"public_id,omitempty"`
	Content       string `json:"content,omitempty"`
	Error         string `json:"error,omitempty"`
}
type UploadImageResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message,omitempty"`
	CloudinaryURL string `json:"cloudinary_url,omitempty"`
	PublicID      string `json:"public_id,omitempty"`
	Error         string `json:"error,omitempty"`
}
