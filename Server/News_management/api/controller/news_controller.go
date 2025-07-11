package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"news-management/internal/model"
	"news-management/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5" // Để check lỗi pgx.ErrNoRows
)

// NewsController chứa các handler cho API tin tức
type NewsController struct {
	newsService service.NewsService
}

// NewNewsController tạo một instance mới của NewsController
func NewNewsController(newsService service.NewsService) *NewsController {
	return &NewsController{newsService: newsService}
}

// CreateNews godoc
// @Summary Tạo tin tức mới
// @Description Tạo một bản ghi tin tức mới với các thông tin được cung cấp.
// @Tags news
// @Accept json
// @Produce json
// @Param news body model.CreateNewsRequest true "Thông tin tin tức"
// @Success 201 {object} model.NewsResponse "Tin tức vừa tạo"
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /news [post]
func (ctrl *NewsController) CreateNews(c *gin.Context) {
	var req model.CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdNews, err := ctrl.newsService.CreateNews(c.Request.Context(), req)
	if err != nil {
		// Phân biệt lỗi do nghiệp vụ (ví dụ: validation) và lỗi hệ thống
		if err.Error() == "title, content, and created_by are required" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, model.ToNewsResponse(createdNews))
}

// GetNewsByID godoc
// @Summary Lấy chi tiết một tin tức
// @Description Lấy thông tin chi tiết của một bản ghi tin tức dựa trên ID.
// @Tags news
// @Produce json
// @Param id path string true "ID của tin tức (UUID)"
// @Success 200 {object} model.NewsResponse "Chi tiết tin tức"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy tin tức"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /news/{id} [get]
func (ctrl *NewsController) GetNewsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	news, err := ctrl.newsService.GetNewsByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "news not found" || errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, model.ToNewsResponse(news))
}

// GetNewsList godoc
// @Summary Lấy danh sách tin tức
// @Description Lấy danh sách các tin tức mới nhất, có thể giới hạn số lượng.
// @Tags news
// @Produce json
// @Param limit query int false "Giới hạn số lượng tin tức (mặc định 10)" default(10)
// @Param offset query int false "Vị trí bắt đầu lấy (mặc định 0)" default(0)
// @Success 200 {array} model.NewsResponse "Danh sách tin tức"
// @Failure 400 {object} map[string]string "Tham số không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /news [get]
func (ctrl *NewsController) GetNewsList(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	newsList, err := ctrl.newsService.GetNewsList(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news list: " + err.Error()})
		return
	}
	if newsList == nil { // Đảm bảo trả về mảng rỗng thay vì null
		c.JSON(http.StatusOK, []model.NewsResponse{})
		return
	}

	c.JSON(http.StatusOK, model.ToListNewsResponse(newsList))
}

// UpdateNews godoc
// @Summary Cập nhật tin tức
// @Description Cập nhật thông tin của một bản ghi tin tức dựa trên ID.
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "ID của tin tức (UUID)"
// @Param news body model.UpdateNewsRequest true "Thông tin cần cập nhật"
// @Success 200 {object} model.NewsResponse "Tin tức đã được cập nhật"
// @Failure 400 {object} map[string]string "ID hoặc dữ liệu không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy tin tức"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /news/{id} [put]
func (ctrl *NewsController) UpdateNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var req model.UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNews, err := ctrl.newsService.UpdateNews(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "news not found for update" || errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		} else if err.Error() == "no fields to update" { // Ví dụ lỗi nghiệp vụ từ service
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, model.ToNewsResponse(updatedNews))
}

// DeleteNews godoc
// @Summary Xóa tin tức
// @Description Xóa một bản ghi tin tức dựa trên ID.
// @Tags news
// @Produce json
// @Param id path string true "ID của tin tức (UUID)"
// @Success 200 {object} map[string]string "Thông báo xóa thành công"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy tin tức"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /news/{id} [delete]
func (ctrl *NewsController) DeleteNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	err = ctrl.newsService.DeleteNews(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "news not found for deletion" || errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}
