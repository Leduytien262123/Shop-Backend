package router

import (
	"backend/internal/handle"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(r *gin.Engine) {
	orderHandler := handle.NewOrderHandler()

	// Routes công khai - KHÔNG YÊU CẦU XÁC THỰC
	publicRoutes := r.Group("/api/public/orders")
	{
		// Tạo đơn hàng khách (không cần xác thực)
		publicRoutes.POST("/", orderHandler.CreateOrder)
		
		// Theo dõi đơn hàng bằng số đơn hàng (không cần xác thực)
		publicRoutes.GET("/track/:order_number", orderHandler.TrackOrderByNumber)
		
		// Tra cứu đơn hàng bằng email/số điện thoại (không cần xác thực)
		publicRoutes.POST("/lookup", orderHandler.LookupGuestOrders)
	}

	// Routes đơn hàng - Yêu cầu xác thực hỗn hợp
	orderRoutes := r.Group("/api/orders")
	{
		// Tạo đơn hàng - KHÔNG YÊU CẦU XÁC THỰC (cả khách và người dùng đã đăng nhập đều có thể sử dụng)
		orderRoutes.POST("/", orderHandler.CreateOrder)
		
		// Lấy đơn hàng cụ thể theo ID - KHÔNG YÊU CẦU XÁC THỰC
		orderRoutes.GET("/:id", orderHandler.GetOrderByID)
	}

	// Routes được bảo vệ chỉ dành cho người dùng đã đăng nhập
	userRoutes := r.Group("/api/orders")
	userRoutes.Use(utils.AuthMiddleware())
	{
		// Lấy đơn hàng của tôi (chỉ dành cho người dùng đã đăng nhập)
		userRoutes.GET("/my", orderHandler.GetMyOrders)
	}

	// Routes admin
	adminRoutes := r.Group("/api/admin/orders")
	adminRoutes.Use(utils.AuthMiddleware())
	adminRoutes.Use(utils.AdminMiddleware())
	{
		adminRoutes.GET("/", orderHandler.GetOrders)
		adminRoutes.GET("/stats", orderHandler.GetOrderStats)
		adminRoutes.GET("/guest-stats", orderHandler.GetGuestOrderStats)
		adminRoutes.GET("/:id", orderHandler.GetOrderByID)
		adminRoutes.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		adminRoutes.PUT("/:id/payment", orderHandler.UpdatePaymentStatus)
	}
}