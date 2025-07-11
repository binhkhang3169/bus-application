// internal/model/news.go
package model

import (
	"news-management/internal/db"
	"time"

	"github.com/google/uuid"
)

// CreateNewsRequest định nghĩa cấu trúc dữ liệu cho yêu cầu tạo tin tức mới
type CreateNewsRequest struct {
	Title     string `json:"title" binding:"required"`
	ImageURL  string `json:"image_url"`
	Content   string `json:"content" binding:"required"`
	CreatedBy string `json:"created_by" binding:"required"`
}

// UpdateNewsRequest định nghĩa cấu trúc dữ liệu cho yêu cầu cập nhật tin tức
type UpdateNewsRequest struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Content  string `json:"content"`
}

// NewsResponse định nghĩa cấu trúc dữ liệu trả về cho client
type NewsResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToNewsResponse chuyển đổi từ sqlc.News sang NewsResponse
func ToNewsResponse(news db.News) NewsResponse {
	return NewsResponse{
		ID:        news.ID,
		Title:     news.Title,
		ImageURL:  news.ImageUrl, // sqlc tạo ImageUrl là sql.NullString
		Content:   news.Content,
		CreatedBy: news.CreatedBy,
		CreatedAt: news.CreatedAt, // sqlc tạo CreatedAt là pgtype.Timestamptz
		UpdatedAt: news.UpdatedAt, // sqlc tạo UpdatedAt là pgtype.Timestamptz
	}
}

// ToListNewsResponse chuyển đổi danh sách sqlc.News sang danh sách NewsResponse
func ToListNewsResponse(newsList []db.News) []NewsResponse {
	responses := make([]NewsResponse, len(newsList))
	for i, news := range newsList {
		responses[i] = ToNewsResponse(news)
	}
	return responses
}
