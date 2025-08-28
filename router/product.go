package router

import (
	"backend/internal/handle"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.Engine) {
	productHandler := handle.NewProductHandler()

	// Public routes
	publicRoutes := r.Group("/api/products")
	{
		publicRoutes.GET("/", productHandler.GetProducts)
		publicRoutes.GET("/:id", productHandler.GetProductByID)
		publicRoutes.GET("/sku/:sku", productHandler.GetProductBySKU)
	}

	// Protected routes (admin only) - TEMPORARILY DISABLED AUTH FOR TESTING
	adminRoutes := r.Group("/api/admin/products")
	// adminRoutes.Use(utils.AuthMiddleware()) // Commented for testing
	// adminRoutes.Use(utils.AdminMiddleware()) // Commented for testing
	{
		adminRoutes.GET("/", productHandler.GetProducts)             // GET all
		adminRoutes.GET("/:id", productHandler.GetProductByID)       // GET by ID
		adminRoutes.POST("/", productHandler.CreateProduct)
		adminRoutes.PUT("/:id", productHandler.UpdateProduct)
		adminRoutes.DELETE("/:id", productHandler.DeleteProduct)
		adminRoutes.PATCH("/:id/stock", productHandler.UpdateProductStock)
	}
}
