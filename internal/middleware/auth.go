package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"meetingroomreservation/internal/utils"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}
		tokenStr := parts[1]
		claims, err := utils.ParseToken(tokenStr, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			c.Abort()
			return
		}

		idRaw, ok := claims["id"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		idFloat, ok := idRaw.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id type"})
			c.Abort()
			return
		}

		userID := uint64(idFloat)

		c.Set("user_id", userID)

		c.Next()
	}
}
