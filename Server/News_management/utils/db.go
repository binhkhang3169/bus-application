package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB thiết lập kết nối đến CSDL sử dụng pgxpool.
func ConnectDB(ctx context.Context, dataSourceName string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("không thể mở kết nối CSDL: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		db.Close() // Đóng kết nối nếu ping thất bại
		return nil, fmt.Errorf("không thể ping CSDL: %w", err)
	}

	log.Println("Kết nối CSDL thành công!")
	return db, nil
}
