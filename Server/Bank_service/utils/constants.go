package utils

// AccountStatus định nghĩa các trạng thái của tài khoản.
type AccountStatus string

const (
	AccountStatusActive AccountStatus = "active"
	AccountStatusClosed AccountStatus = "closed"
	// Thêm các trạng thái khác nếu cần, ví dụ: "frozen", "pending_verification"
)

// SupportedCurrencies định nghĩa các loại tiền tệ được hỗ trợ.
// Có thể tải từ config hoặc DB trong ứng dụng thực tế.
var SupportedCurrencies = map[string]bool{
	"VND": true,
	"USD": true,
	"EUR": true,
}

// IsSupportedCurrency kiểm tra xem một loại tiền tệ có được hỗ trợ không.
func IsSupportedCurrency(currency string) bool {
	_, ok := SupportedCurrencies[currency]
	return ok
}
