package router

import (
	"backend/app"
	"backend/internal/handle"
	"backend/internal/repo"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	// Khởi tạo repository và handler
	userRepo := repo.NewUserRepository(app.GetDB())
	authHandler := handle.NewAuthHandler(userRepo)

	// Routes công khai
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Routes được bảo vệ
	protected := router.Group("/api/auth")
	protected.Use(utils.AuthMiddleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.PUT("/profile", authHandler.UpdateProfile)
	}
}