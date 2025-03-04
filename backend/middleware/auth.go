package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"super-supply-chain/configs"
	"super-supply-chain/controllers"
	"time"
)

func handleNoAuth(c *gin.Context) {
	// 重定向
	cookie := &http.Cookie{
		Name:     configs.AuthKey,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24,
	}
	http.SetCookie(c.Writer, cookie)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		// 从cookie 获取 Authorization
		Authorization, err := c.Cookie("Authorization")
		if err != nil {
			handleNoAuth(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
		}
		authHeader = Authorization

		if authHeader == "" {
			handleNoAuth(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &controllers.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return controllers.JwtKey, nil
		})

		if err != nil || !token.Valid {
			handleNoAuth(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			handleNoAuth(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
