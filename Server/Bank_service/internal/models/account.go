package models

import (
	"time"
)

// AccountResponse định nghĩa cấu trúc trả về cho client khi lấy thông tin tài khoản.
// Có thể tùy chỉnh để chỉ trả về các trường cần thiết.
type AccountResponse struct {
	ID        int64     `json:"id"`
	OwnerName string    `json:"owner_name"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAccountRequest định nghĩa cấu trúc request để tạo tài khoản mới.
type CreateAccountRequest struct {
	OwnerName string `json:"owner_name" binding:"required"`
	Currency  string `json:"currency" binding:"required,currency"` // currency là custom validator tag
	Balance   int64  `json:"balance"`                              // Cho phép nhập balance ban đầu, mặc định là 0 nếu không có
}

// DepositRequest định nghĩa cấu trúc request để nạp tiền.
type DepositRequest struct {
	Amount   int64  `json:"amount" binding:"required,gt=0"`       // Số tiền phải lớn hơn 0
	Currency string `json:"currency" binding:"required,currency"` // Kiểm tra currency nếu cần
}

// PaymentRequest định nghĩa cấu trúc request để thanh toán.
type PaymentRequest struct {
	Amount   int64  `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

// UpdateStatusRequest định nghĩa cấu trúc request để cập nhật trạng thái tài khoản.
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active closed"` // Chỉ chấp nhận 'active' hoặc 'closed'
}

// GetAccountRequest định nghĩa params cho việc lấy thông tin tài khoản (qua URI).
type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// ListAccountsRequest định nghĩa query params cho việc liệt kê tài khoản.
type ListAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

// Helper function để chuyển đổi db.Account (sqlc generated) sang AccountResponse
// Bạn sẽ cần import package db của sqlc vào đây
// Ví dụ: import "bank/db/sqlc"
// func ToAccountResponse(dbAccount db.Account) AccountResponse {
// 	return AccountResponse{
// 		ID:        dbAccount.ID,
// 		OwnerName: dbAccount.OwnerName,
// 		Balance:   dbAccount.Balance,
// 		Currency:  dbAccount.Currency,
// 		Status:    dbAccount.Status,
// 		CreatedAt: dbAccount.CreatedAt,
// 		UpdatedAt: dbAccount.UpdatedAt,
// 	}
// }

type TransactionType string

const (
	TransactionTypeCreateAccount TransactionType = "CREATE_ACCOUNT"
	TransactionTypeDeposit       TransactionType = "DEPOSIT"
	TransactionTypePayment       TransactionType = "PAYMENT"
	TransactionTypeCloseAccount  TransactionType = "CLOSE_ACCOUNT"
)

// ListTransactionHistoryRequest defines parameters for listing transaction history.
type ListTransactionHistoryRequest struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

// TransactionHistoryResponse defines the API response for a single transaction history entry.
type TransactionHistoryResponse struct {
	ID                   int64           `json:"id"`
	AccountID            int64           `json:"account_id"`
	TransactionType      TransactionType `json:"transaction_type"`
	Amount               *int64          `json:"amount,omitempty"`   // Pointer to allow null
	Currency             *string         `json:"currency,omitempty"` // Pointer to allow null
	TransactionTimestamp string          `json:"transaction_timestamp"`
	Description          string          `json:"description,omitempty"`
}
