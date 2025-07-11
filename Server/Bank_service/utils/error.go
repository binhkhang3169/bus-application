package utils

import (
	"bank/internal/models"
	"errors"
	"fmt"
	"net/http"
)

// AppError là cấu trúc lỗi tùy chỉnh cho ứng dụng.
type AppError struct {
	Code    int    `json:"-"` // HTTP status code, không export ra JSON message
	Message string `json:"message"`
	Err     error  `json:"-"` // Lỗi gốc, không export ra JSON message
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Message: %s, OriginalError: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError tạo một AppError mới.
func NewAppError(message string, code int, originalError ...error) *AppError {
	var err error
	if len(originalError) > 0 {
		err = originalError[0]
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Các hàm tạo lỗi cụ thể
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(message, http.StatusBadRequest, err)
}

func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(message, http.StatusNotFound, err)
}

func NewInternalServerError(message string, err error) *AppError {
	// Log lỗi ở đây nếu cần
	// log.Printf("Internal Server Error: %s, Original: %v\n", message, err)
	return NewAppError(message, http.StatusInternalServerError, err)
}

func NewUnauthorizedError(message string, err error) *AppError {
	return NewAppError(message, http.StatusUnauthorized, err)
}

func NewForbiddenError(message string, err error) *AppError {
	return NewAppError(message, http.StatusForbidden, err)
}

// TransactionError dùng để gói lỗi từ transaction và lỗi rollback.
type TransactionError struct {
	OriginalErr error
	RollbackErr error
}

func (te *TransactionError) Error() string {
	return fmt.Sprintf("lỗi giao dịch: %v (lỗi rollback: %v)", te.OriginalErr, te.RollbackErr)
}

func NewTransactionError(originalErr, rollbackErr error) *TransactionError {
	return &TransactionError{OriginalErr: originalErr, RollbackErr: rollbackErr}
}

// DetermineStatusCode ánh xạ lỗi nghiệp vụ sang HTTP status code
func DetermineStatusCode(err error) int {
	switch {
	case errors.Is(err, models.ErrAccountNotFound):
		return http.StatusNotFound
	case errors.Is(err, models.ErrInsufficientFunds):
		return http.StatusPaymentRequired // 402 Payment Required
	case errors.Is(err, models.ErrInvalidAccountStatus):
		return http.StatusUnprocessableEntity // 422 Unprocessable Entity
	case errors.Is(err, models.ErrCurrencyMismatch):
		return http.StatusUnprocessableEntity
	// Thêm các case khác nếu cần
	default:
		return http.StatusInternalServerError
	}
}

// HandleServiceError chuyển đổi lỗi từ service layer sang AppError cho controller
func HandleServiceError(err error, defaultMessage string) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) { // Nếu lỗi đã là AppError, trả về nó
		return appErr
	}

	// Nếu là các lỗi nghiệp vụ đã biết từ service
	statusCode := DetermineStatusCode(err)
	if statusCode != http.StatusInternalServerError {
		return NewAppError(err.Error(), statusCode, err)
	}

	// Nếu không, coi là lỗi server nội bộ
	return NewInternalServerError(defaultMessage, err)
}
