package utils

import (
	db "bank/internal/db" // Đường dẫn tới package db của sqlc
	"bank/internal/models"
)

// ToAccountResponse chuyển đổi từ db.Account (do sqlc tạo ra) sang model.AccountResponse.
// Điều này giúp tách biệt cấu trúc dữ liệu CSDL và cấu trúc response API.
func ToAccountResponse(dbAccount db.Account) models.AccountResponse {
	return models.AccountResponse{
		ID:        dbAccount.ID,
		OwnerName: dbAccount.OwnerName,
		Balance:   dbAccount.Balance,
		Currency:  dbAccount.Currency,
		Status:    dbAccount.Status,
		CreatedAt: dbAccount.CreatedAt,
		UpdatedAt: dbAccount.UpdatedAt,
	}
}

// ToAccountResponses chuyển đổi một slice db.Account sang slice model.AccountResponse.
func ToAccountResponses(dbAccounts []db.Account) []models.AccountResponse {
	responses := make([]models.AccountResponse, len(dbAccounts))
	for i, acc := range dbAccounts {
		responses[i] = ToAccountResponse(acc)
	}
	return responses
}
