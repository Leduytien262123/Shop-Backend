package router

import (
	"backend/internal/handle"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.Engine) {
	productHandler := handle.NewProductHandler()

	// Routes công khai
	publicRoutes := r.Group("/api/products")
	{
		publicRoutes.GET("/", productHandler.GetProducts)
		publicRoutes.GET("/:id", productHandler.GetProductByID)
		publicRoutes.GET("/sku/:sku", productHandler.GetProductBySKU)
	}

	// Routes được bảo vệ (admin và owner)
	adminRoutes := r.Group("/api/admin/products")
	adminRoutes.Use(utils.AuthMiddleware())
	adminRoutes.Use(utils.AdminMiddleware()) // Bây giờ cho phép cả admin và owner
	{
		adminRoutes.GET("/", productHandler.GetProducts)
		adminRoutes.GET("/:id", productHandler.GetProductByID)
		adminRoutes.POST("/", productHandler.CreateProduct)
		adminRoutes.PUT("/:id", productHandler.UpdateProduct)
		adminRoutes.DELETE("/:id", productHandler.DeleteProduct)
		adminRoutes.PATCH("/:id/stock", productHandler.UpdateProductStock)
	}
}
