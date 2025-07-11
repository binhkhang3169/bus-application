// internal/repository/news_repository.go
package repository

import (
	"context"
	"database/sql" // Hoặc "github.com/jackc/pgx/v5/pgxpool" nếu dùng pgx pool
	database "news-management/internal/db"

	"github.com/google/uuid"
	// Cần cho pgtype
)

// NewsRepository định nghĩa interface cho các thao tác với dữ liệu tin tức
type NewsRepository interface {
	CreateNews(ctx context.Context, arg database.CreateNewsParams) (database.News, error)
	GetNewsByID(ctx context.Context, id uuid.UUID) (database.News, error)
	ListNews(ctx context.Context, arg database.ListNewsParams) ([]database.News, error)
	UpdateNews(ctx context.Context, arg database.UpdateNewsParams) (database.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error
}

// newsRepositoryImpl triển khai NewsRepository
type newsRepositoryImpl struct {
	db      *sql.DB // Hoặc *pgxpool.Pool
	queries *database.Queries
}

// NewNewsRepository tạo một instance mới của newsRepositoryImpl
func NewNewsRepository(db *sql.DB) NewsRepository { // Hoặc *pgxpool.Pool
	return &newsRepositoryImpl{
		db:      db,
		queries: database.New(db), // database.New() nhận DBTX, có thể là *sql.DB hoặc *sql.Tx
	}
}

// CreateNews thêm một tin tức mới vào cơ sở dữ liệu
func (r *newsRepositoryImpl) CreateNews(ctx context.Context, arg database.CreateNewsParams) (database.News, error) {
	// Chuyển đổi từ string sang pgtype.Text cho các trường có thể NULL nếu cần
	// Trong database.yaml, nếu dùng sql_package: "pgx/v5", các kiểu NullXXX sẽ là pgtype.XXX
	// Ở đây, CreateNewsParams đã được database định nghĩa phù hợp.
	return r.queries.CreateNews(ctx, arg)
}

// GetNewsByID lấy một tin tức từ cơ sở dữ liệu bằng ID
func (r *newsRepositoryImpl) GetNewsByID(ctx context.Context, id uuid.UUID) (database.News, error) {
	return r.queries.GetNewsByID(ctx, id)
}

// ListNews lấy danh sách tin tức từ cơ sở dữ liệu với phân trang
func (r *newsRepositoryImpl) ListNews(ctx context.Context, arg database.ListNewsParams) ([]database.News, error) {
	return r.queries.ListNews(ctx, arg)
}

// UpdateNews cập nhật một tin tức trong cơ sở dữ liệu
func (r *newsRepositoryImpl) UpdateNews(ctx context.Context, arg database.UpdateNewsParams) (database.News, error) {
	// Xử lý pgtype.Text cho các trường có thể NULL
	// Ví dụ, nếu arg.Title là string và cột title có thể NULL, bạn cần chuyển đổi
	// Tuy nhiên, UpdateNewsParams của database thường đã xử lý việc này
	return r.queries.UpdateNews(ctx, arg)
}

// DeleteNews xóa một tin tức khỏi cơ sở dữ liệu bằng ID
func (r *newsRepositoryImpl) DeleteNews(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteNews(ctx, id)
}
