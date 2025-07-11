package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import a blank Pq driver
)

// ConnectDB thiết lập kết nối đến CSDL.
func ConnectDB(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("không thể mở kết nối CSDL: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close() // Đóng kết nối nếu ping thất bại
		return nil, fmt.Errorf("không thể ping CSDL: %w", err)
	}

	log.Println("Kết nối CSDL thành công!")
	return db, nil
}
