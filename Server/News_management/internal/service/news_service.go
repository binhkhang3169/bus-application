// internal/service/news_service.go
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	database "news-management/internal/db"
	"news-management/internal/model"
	"news-management/internal/repository"
	"news-management/pkg/kafkaclient"
	"time"

	"github.com/google/uuid"
	// Cần cho pgtype
)

// NewsService định nghĩa interface cho các nghiệp vụ liên quan đến tin tức
type NewsService interface {
	CreateNews(ctx context.Context, req model.CreateNewsRequest) (database.News, error)
	GetNewsByID(ctx context.Context, id uuid.UUID) (database.News, error)
	GetNewsList(ctx context.Context, limit, offset int) ([]database.News, error)
	UpdateNews(ctx context.Context, id uuid.UUID, req model.UpdateNewsRequest) (database.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error
}

type newsServiceImpl struct {
	repo      repository.NewsRepository
	publisher *kafkaclient.Publisher // << ADDED

}

// NewNewsService tạo một instance mới của newsServiceImpl
func NewNewsService(repo repository.NewsRepository, publisher *kafkaclient.Publisher) NewsService {
	return &newsServiceImpl{
		repo:      repo,
		publisher: publisher,
	}
}

// CreateNews xử lý nghiệp vụ tạo tin tức mới
func (s *newsServiceImpl) CreateNews(ctx context.Context, req model.CreateNewsRequest) (database.News, error) {
	// Validate input (có thể thêm các logic phức tạp hơn ở đây)
	if req.Title == "" || req.Content == "" || req.CreatedBy == "" {
		return database.News{}, errors.New("title, content, and created_by are required")
	}

	params := database.CreateNewsParams{
		Title:     req.Title,
		Content:   req.Content,
		CreatedBy: req.CreatedBy,
	}
	if req.ImageURL != "" {
		params.ImageUrl = req.ImageURL
	} else {
		params.ImageUrl = "" // Hoặc để mặc định là NULL trong DB
	}

	newNews, err := s.repo.CreateNews(ctx, params)
	if err != nil {
		return database.News{}, fmt.Errorf("failed to create news in repository: %w", err)
	}

	// 2. << ADDED: Gửi sự kiện thông báo đến Kafka sau khi tạo thành công >>
	go func() {
		// Tạo context riêng cho goroutine để không bị ảnh hưởng
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		notificationTopic := "notifications_topic" // Tên topic nên được lấy từ config

		event := kafkaclient.NotificationEvent{
			UserID:  nil, // UserID là nil để gửi thông báo cho tất cả mọi người (broadcast)
			Type:    "NEW_ARTICLE",
			Title:   "Tin tức mới!",
			Message: fmt.Sprintf("Có bài viết mới: %s", newNews.Title),
		}

		// Key có thể là nil hoặc một giá trị cố định cho broadcast
		if err := s.publisher.Publish(bgCtx, notificationTopic, nil, event); err != nil {
			// Ghi log lỗi nghiêm trọng nếu không thể gửi sự kiện
			log.Printf("CRITICAL: Failed to publish new article notification for news ID %s: %v", newNews.ID, err)
		} else {
			log.Printf("Successfully published notification for new article: %s", newNews.ID)
		}
	}()

	return newNews, nil
}

// GetNewsByID xử lý nghiệp vụ lấy chi tiết tin tức
func (s *newsServiceImpl) GetNewsByID(ctx context.Context, id uuid.UUID) (database.News, error) {
	news, err := s.repo.GetNewsByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // hoặc pgx.ErrNoRows nếu dùng pgx trực tiếp
			return database.News{}, errors.New("news not found")
		}
		return database.News{}, err
	}
	return news, nil
}

// GetNewsList xử lý nghiệp vụ lấy danh sách tin tức
func (s *newsServiceImpl) GetNewsList(ctx context.Context, limit, offset int) ([]database.News, error) {
	if limit <= 0 {
		limit = 10 // Giá trị mặc định
	}
	if offset < 0 {
		offset = 0 // Giá trị mặc định
	}
	params := database.ListNewsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	return s.repo.ListNews(ctx, params)
}

// UpdateNews xử lý nghiệp vụ cập nhật tin tức
func (s *newsServiceImpl) UpdateNews(ctx context.Context, id uuid.UUID, req model.UpdateNewsRequest) (database.News, error) {
	// Lấy tin tức hiện tại để kiểm tra tồn tại và lấy các giá trị không thay đổi
	existingNews, err := s.repo.GetNewsByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // hoặc pgx.ErrNoRows
			return database.News{}, errors.New("news not found for update")
		}
		return database.News{}, err
	}

	// Tạo params với giá trị hiện tại, sau đó cập nhật nếu có giá trị mới từ request
	updateParams := database.UpdateNewsParams{
		ID:       id,
		Title:    existingNews.Title,
		Content:  existingNews.Content,
		ImageUrl: existingNews.ImageUrl, // Giữ giá trị ImageUrl hiện tại
	}

	if req.Title != "" {
		updateParams.Title = req.Title
	}
	if req.Content != "" {
		updateParams.Content = req.Content
	}
	// Chỉ cập nhật ImageUrl nếu nó được cung cấp trong request
	// Nếu muốn xóa ImageUrl, client có thể gửi một giá trị đặc biệt hoặc để trống
	if req.ImageURL != "" { // Giả sử "" có nghĩa là không thay đổi, nếu muốn xoá thì client gửi "null" hoặc có field riêng
		updateParams.ImageUrl = req.ImageURL
	}
	// Nếu bạn muốn cho phép xóa ImageURL bằng cách gửi rỗng:
	// else if req.ImageURL == "" { // Cẩn thận: điều này có thể không phải ý muốn nếu "" là giá trị hợp lệ
	// 	updateParams.ImageUrl = pgtype.Text{Valid: false}
	// }

	return s.repo.UpdateNews(ctx, updateParams)
}

// DeleteNews xử lý nghiệp vụ xóa tin tức
func (s *newsServiceImpl) DeleteNews(ctx context.Context, id uuid.UUID) error {
	// Kiểm tra tin tức tồn tại trước khi xóa (tùy chọn, vì DB sẽ báo lỗi nếu không tìm thấy)
	_, err := s.repo.GetNewsByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // hoặc pgx.ErrNoRows
			return errors.New("news not found for deletion")
		}
		return err
	}
	return s.repo.DeleteNews(ctx, id)
}
