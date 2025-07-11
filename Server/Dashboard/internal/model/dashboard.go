package model

// StatItem là một cấu trúc chung cho các dữ liệu thống kê theo danh mục.
type StatItem struct {
	Category string  `json:"category"`
	Value    float64 `json:"value"`
}

// TimeSeriesDataPoint là điểm dữ liệu cho biểu đồ theo thời gian.
type TimeSeriesDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// KpiData chứa các chỉ số hiệu suất chính.
type KpiData struct {
	TotalRevenue  float64 `json:"total_revenue"`
	TotalTickets  int64   `json:"total_tickets"`
	TotalInvoices int64   `json:"total_invoices"`
	UnpaidTickets int64   `json:"unpaid_tickets"`
}

// TicketDistributionData chứa dữ liệu phân bổ vé.
type TicketDistributionData struct {
	ByChannel []StatItem `json:"by_channel"`
	ByType    []StatItem `json:"by_type"`
}
