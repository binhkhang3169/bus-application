// payment_service/pkg/redisclient/client.go
package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewClient tạo một Redis client mới từ một URL kết nối (tương thích với Upstash)
func NewClient(redisURL string) *redis.Client {
	// THAY ĐỔI: Sử dụng redis.ParseURL để tự động cấu hình TLS và Auth từ URL
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Sprintf("Không thể phân giải URL của Redis: %v", err))
	}

	rdb := redis.NewClient(opts)

	// Kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(fmt.Sprintf("Không thể kết nối đến Redis (Upstash): %v", err))
	}

	fmt.Println("Đã kết nối thành công đến Redis (Upstash).")
	return rdb
}
