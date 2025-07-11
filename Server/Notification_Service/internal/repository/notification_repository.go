// notification-service/internal/repository/notification_repository.go
package repository

import (
	"context"
	"fmt"
	db "notification-service/internal/db" // Giả sử module name là "notification-service"

	"github.com/jackc/pgx/v5" // Cần cho pgx.Tx
	// "github.com/jackc/pgx/v5/pgtype" // Vẫn cần nếu có các hàm wrapper còn lại
)

// Store interface bao gồm tất cả các hàm của db.Querier (do sqlc tạo ra)
// và thêm phương thức ExecTx để quản lý transaction.
type Store interface {
	db.Querier
	ExecTx(ctx context.Context, fn func(*db.Queries) error) error
}

// SQLStore cung cấp tất cả các chức năng để thực thi các truy vấn DB.
type SQLStore struct {
	*db.Queries         // Queries được tạo bởi sqlc (khởi tạo với connPool)
	connPool    db.DBTX // Giữ connection pool (ví dụ *pgxpool.Pool)
}

// NewStore tạo một Store mới.
// connPool nên là một kiểu tương thích với pgx/v5 như *pgxpool.Pool.
func NewStore(connPool db.DBTX) Store {
	return &SQLStore{
		Queries:  db.New(connPool), // db.New từ sqlc mong đợi DBTX
		connPool: connPool,
	}
}

// ExecTx thực thi một hàm fn bên trong một transaction của cơ sở dữ liệu.
// Nếu hàm fn trả về lỗi, transaction sẽ được rollback.
// Nếu hàm fn không có lỗi, transaction sẽ được commit.
func (s *SQLStore) ExecTx(ctx context.Context, fn func(*db.Queries) error) error {
	// Cần một đối tượng có thể bắt đầu transaction (ví dụ: *pgxpool.Pool)
	// db.DBTX không đảm bảo có phương thức Begin().
	// *pgxpool.Pool và pgx.Tx (cho nested transactions, nếu DB hỗ trợ) có Begin().
	txStarter, ok := s.connPool.(interface {
		Begin(context.Context) (pgx.Tx, error)
	})
	if !ok {
		return fmt.Errorf("DBTX connection does not support starting transactions")
	}

	tx, err := txStarter.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Đảm bảo rollback nếu có lỗi hoặc hàm kết thúc mà không commit
	// Lưu ý: tx.Rollback() là no-op nếu tx đã được commit.
	defer tx.Rollback(ctx)

	// s.Queries là *db.Queries được khởi tạo với connPool.
	// Sử dụng phương thức WithTx của nó để có được một *db.Queries mới hoạt động trên transaction tx.
	qtx := s.Queries.WithTx(tx)
	err = fn(qtx) // Thực thi hàm callback với Queries của transaction
	if err != nil {
		// Nếu fn trả về lỗi, defer tx.Rollback(ctx) sẽ xử lý.
		return fmt.Errorf("transaction callback failed: %w", err)
	}

	// Nếu không có lỗi, commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Các phương thức wrapper như CreateNotification, GetNotificationsByUserID trong SQLStore
// không còn thực sự cần thiết nếu service layer gọi trực tiếp các phương thức từ db.Querier
// (thông qua store.SomeQuery(ctx, params)) hoặc sử dụng ExecTx.
// Nếu giữ lại, chúng không nên tự quản lý transaction.
// Ví dụ:
/*
func (s *SQLStore) CreateNotification(ctx context.Context, arg db.CreateNotificationParams) (db.Notification, error) {
	return s.Queries.CreateNotification(ctx, arg)
}

func (s *SQLStore) GetNotificationsByUserID(ctx context.Context, arg db.GetNotificationsByUserIDParams) ([]db.Notification, error) {
	return s.Queries.GetNotificationsByUserID(ctx, arg)
}
*/
