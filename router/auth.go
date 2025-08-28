package router

import (
	"backend/app"
	"backend/internal/handle"
	"backend/internal/repo"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	// Initialize repository and handler
	userRepo := repo.NewUserRepository(app.GetDB())
	authHandler := handle.NewAuthHandler(userRepo)

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/auth")
	protected.Use(utils.AuthMiddleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.PUT("/profile", authHandler.UpdateProfile)
	}
}