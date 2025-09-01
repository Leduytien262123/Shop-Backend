package router

import (
	"backend/internal/handle"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupCartRoutes(r *gin.Engine) {
	cartHandler := handle.NewCartHandler()

	// Routes được bảo vệ (chỉ người dùng)
	cartRoutes := r.Group("/api/cart")
	cartRoutes.Use(utils.AuthMiddleware())
	{
		cartRoutes.GET("/", cartHandler.GetCart)
		cartRoutes.POST("/add", cartHandler.AddToCart)
		cartRoutes.PUT("/items/:product_id", cartHandler.UpdateCartItem)
		cartRoutes.DELETE("/items/:product_id", cartHandler.RemoveFromCart)
		cartRoutes.DELETE("/clear", cartHandler.ClearCart)
	}
}

func SetupNewsRoutes(r *gin.Engine) {
	newsHandler := handle.NewNewsHandler()

	// Routes công khai
	publicRoutes := r.Group("/api/news")
	{
		publicRoutes.GET("/", newsHandler.GetNews)
		publicRoutes.GET("/:id", newsHandler.GetNewsByID)
		publicRoutes.GET("/slug/:slug", newsHandler.GetNewsBySlug)
	}

	// Routes admin
	adminRoutes := r.Group("/api/admin/news")
	adminRoutes.Use(utils.AuthMiddleware())
	adminRoutes.Use(utils.AdminMiddleware())
	{
		adminRoutes.GET("/", newsHandler.GetNews) // Tất cả tin tức bao gồm cả chưa xuất bản
		adminRoutes.GET("/:id", newsHandler.GetNewsByID)
		adminRoutes.POST("/", newsHandler.CreateNews)
		adminRoutes.PUT("/:id", newsHandler.UpdateNews)
		adminRoutes.DELETE("/:id", newsHandler.DeleteNews)
	}
}