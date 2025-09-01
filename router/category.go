package router

import (
	"backend/internal/handle"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(r *gin.Engine) {
	categoryHandler := handle.NewCategoryHandler()

	// Routes công khai
	publicRoutes := r.Group("/api/categories")
	{
		publicRoutes.GET("/", categoryHandler.GetCategories)
		publicRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		publicRoutes.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
	}

	// Routes được bảo vệ (admin và owner)
	adminRoutes := r.Group("/api/admin/categories")
	adminRoutes.Use(utils.AuthMiddleware())
	adminRoutes.Use(utils.AdminMiddleware()) // Bây giờ cho phép cả admin và owner
	{
		adminRoutes.GET("/", categoryHandler.GetCategories) 
		adminRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		adminRoutes.POST("/", categoryHandler.CreateCategory)
		adminRoutes.PUT("/:id", categoryHandler.UpdateCategory)
		adminRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
