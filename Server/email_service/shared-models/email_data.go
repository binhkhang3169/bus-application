// File: shared-models/email_data.go
package models

// Dữ liệu chung cho các loại email đơn giản
type GenericEmailData struct {
	Data string `json:"data"` // Có thể là OTP, URL, hoặc một thông báo đơn giản
}

// Thông tin về điểm đi/đến
type LocationInfo struct {
	LocationName string `json:"locationName"` // Ví dụ: "Bến xe Miền Đông"
	Address      string `json:"address"`      // Ví dụ: "Quận Bình Thạnh, TP.HCM"
	Time         string `json:"time"`         // Ví dụ: "9:00 PM"
	Date         string `json:"date"`         // Ví dụ: "Thứ Ba, 20/02/2025"
}

// Dữ liệu cho email xác nhận vé (phù hợp với template bạn cung cấp)
type TicketConfirmationData struct {
	CustomerName    string       `json:"customerName"`    // "John"
	TicketID        string       `json:"ticketId"`        // Mã vé để hiển thị
	TripDate        string       `json:"tripDate"`        // Ngày chính của chuyến đi
	TripTime        string       `json:"tripTime"`        // Giờ chính của chuyến đi
	DepartureInfo   LocationInfo `json:"departureInfo"`   // Thông tin điểm đi
	ArrivalInfo     LocationInfo `json:"arrivalInfo"`     // Thông tin điểm đến
	ConfirmationURL string       `json:"confirmationUrl"` // URL cho nút "Confirm Now"
	QRCodeImageURL  string       `json:"qrCodeImageUrl"`  // URL của hình ảnh QR code
}

// Dữ liệu cho email thông báo hoàn tiền
type RefundNotificationData struct {
	CustomerName string `json:"customerName"`
	TicketID     string `json:"ticketId"`
}
