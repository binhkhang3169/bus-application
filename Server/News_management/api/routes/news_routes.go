// routes/news_routes.go
package routes

import (
	"news-management/api/controller"

	"github.com/gin-gonic/gin"
)

// SetupNewsRoutes thiết lập các route cho module tin tức
func SetupNewsRoutes(router *gin.RouterGroup, newsController *controller.NewsController) {
	newsRoutes := router.Group("/news")
	{
		newsRoutes.POST("", newsController.CreateNews)
		newsRoutes.GET("", newsController.GetNewsList)     // /news?limit=5&offset=0
		newsRoutes.GET("/:id", newsController.GetNewsByID) // /news/uuid-goes-here
		newsRoutes.PUT("/:id", newsController.UpdateNews)
		newsRoutes.DELETE("/:id", newsController.DeleteNews)
	}
}
