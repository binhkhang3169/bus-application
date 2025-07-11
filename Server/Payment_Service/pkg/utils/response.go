package utils

import (
	"github.com/gin-gonic/gin"
)

// StandardAPIResponse định nghĩa cấu trúc response chuẩn
type StandardAPIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError định nghĩa cấu trúc lỗi
type APIError struct {
	Code    int         `json:"code,omitempty"` // Mã lỗi HTTP hoặc mã lỗi nội bộ
	Details interface{} `json:"details,omitempty"`
}

// RespondWithSuccess gửi response thành công
func RespondWithSuccess(c *gin.Context, httpStatusCode int, message string, data interface{}) {
	c.JSON(httpStatusCode, StandardAPIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondWithError gửi response lỗi
func RespondWithError(c *gin.Context, httpStatusCode int, message string, errorDetails interface{}) {
	errResponse := StandardAPIResponse{
		Success: false,
		Message: message,
	}
	if errorDetails != nil {
		errResponse.Error = &APIError{
			Code:    httpStatusCode, // Có thể dùng mã lỗi nội bộ khác nếu cần
			Details: errorDetails,
		}
	}
	c.JSON(httpStatusCode, errResponse)
}
