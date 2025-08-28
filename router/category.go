package router

import (
	"backend/internal/handle"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(r *gin.Engine) {
	categoryHandler := handle.NewCategoryHandler()

	// Public routes
	publicRoutes := r.Group("/api/categories")
	{
		publicRoutes.GET("/", categoryHandler.GetCategories)
		publicRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		publicRoutes.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
	}

	// Protected routes (admin only) - TEMPORARILY DISABLED AUTH FOR TESTING
	adminRoutes := r.Group("/api/admin/categories")
	// adminRoutes.Use(utils.AuthMiddleware()) // Commented for testing
	// adminRoutes.Use(utils.AdminMiddleware()) // Commented for testing
	{
		adminRoutes.GET("/", categoryHandler.GetCategories) 
		adminRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		adminRoutes.POST("/", categoryHandler.CreateCategory)
		adminRoutes.PUT("/:id", categoryHandler.UpdateCategory)
		adminRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
