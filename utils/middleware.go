package utils

import (
	"backend/internal/consts"
	"backend/internal/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.UnauthorizedResponse(c, consts.MSG_UNAUTHORIZED)
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			helpers.UnauthorizedResponse(c, consts.MSG_UNAUTHORIZED)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		token, err := helpers.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			helpers.UnauthorizedResponse(c, consts.MSG_UNAUTHORIZED)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.UnauthorizedResponse(c, consts.MSG_UNAUTHORIZED)
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("role", claims["role"].(string))

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != consts.ROLE_ADMIN {
			helpers.ForbiddenResponse(c, consts.MSG_FORBIDDEN)
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}
