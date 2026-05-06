package middleware

import (
	"net/http"
	"strings"

	"github.com/enzo959/forumium/config"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token manquant"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "format token invalide"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		claims, err := config.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalide ou expiré"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
