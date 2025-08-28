package router

import (
	"backend/app"
	"backend/internal/handle"
	"backend/internal/repo"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.Engine) {
	userRepo := repo.NewUserRepository(app.GetDB())
	adminHandler := handle.NewAdminHandler(userRepo)

	admin := router.Group("/api/admin")
	admin.Use(utils.AuthMiddleware())
	admin.Use(utils.AdminMiddleware())
	{
		admin.GET("/users", adminHandler.GetAllUsers)
		admin.GET("/users/:id", adminHandler.GetUserByID)
		admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
		admin.PUT("/users/:id/status", adminHandler.ToggleUserStatus)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
	}
}