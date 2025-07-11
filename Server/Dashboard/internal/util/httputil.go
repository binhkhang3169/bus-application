package util

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ParseTimeRange lấy và xác thực start_date, end_date từ query params.
// Nó trả về ngày bắt đầu và ngày kết thúc (đã được +1 để truy vấn BigQuery)
func ParseTimeRange(c *gin.Context) (string, string, bool) {
	const layout = "2006-01-02"
	// Mặc định: 30 ngày gần nhất
	defaultEnd := time.Now()
	defaultStart := defaultEnd.AddDate(0, 0, -30)

	startDateStr := c.DefaultQuery("start_date", defaultStart.Format(layout))
	endDateStr := c.DefaultQuery("end_date", defaultEnd.Format(layout))

	// Xác thực định dạng ngày
	_, errStart := time.Parse(layout, startDateStr)
	endDate, errEnd := time.Parse(layout, endDateStr)
	if errStart != nil || errEnd != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return "", "", false
	}

	// BigQuery WHERE clause (field < endDate) sẽ không bao gồm ngày cuối cùng.
	// Vì vậy, ta cần lấy ngày kế tiếp của endDate để đảm bảo bao trọn ngày cuối.
	nextDayEndStr := endDate.AddDate(0, 0, 1).Format(layout)

	return startDateStr, nextDayEndStr, true
}
