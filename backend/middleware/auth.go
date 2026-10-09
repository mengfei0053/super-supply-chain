package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"super-supply-chain/configs"
	"super-supply-chain/controllers"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

const (
	contextUsername   = "username"
	contextUserID     = "userID"
	contextAuthMethod = "authMethod"
	contextQueryToken = "queryToken"
	authMethodJWT     = "jwt"
	authMethodPAT     = "pat"
)

func authCookieName() string {
	if configs.AuthKey != "" {
		return configs.AuthKey
	}
	return "Authorization"
}

func handleNoAuth(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     authCookieName(),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24,
	}
	http.SetCookie(c.Writer, cookie)
}

func rejectAuth(c *gin.Context, message string) {
	handleNoAuth(c)
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
	c.Abort()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := credentialFromRequest(c)
		if raw == "" {
			rejectAuth(c, "Authorization header required")
			return
		}
		if models.IsPersonalAccessToken(raw) {
			authenticatePersonalAccessToken(c, raw)
			return
		}
		authenticateJWT(c, raw)
	}
}

// RequireSessionJWT blocks personal access tokens from minting or revoking tokens.
// It must run after AuthMiddleware.
func RequireSessionJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(contextAuthMethod) != authMethodJWT {
			c.JSON(http.StatusForbidden, gin.H{"error": "Personal access token management requires a session login"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func bearerValue(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	}
	return raw
}

func queryToken(c *gin.Context) string {
	if value, ok := c.Get(contextQueryToken); ok {
		if token, ok := value.(string); ok {
			return strings.TrimSpace(token)
		}
	}
	return strings.TrimSpace(c.Query("token"))
}

func credentialFromRequest(c *gin.Context) string {
	if key := bearerValue(c.GetHeader("X-API-Key")); key != "" {
		return key
	}

	headerToken := bearerValue(c.GetHeader("Authorization"))
	// A PAT in Authorization is an API client. Prefer it over a browser session cookie.
	if models.IsPersonalAccessToken(headerToken) {
		return headerToken
	}

	if cookie, err := c.Cookie(authCookieName()); err == nil {
		if cookieToken := bearerValue(cookie); cookieToken != "" {
			return cookieToken
		}
	}
	if headerToken != "" {
		return headerToken
	}
	return queryToken(c)
}

func authenticateJWT(c *gin.Context, tokenString string) {
	claims := &controllers.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return controllers.JwtKey, nil
	})
	if err != nil || !token.Valid {
		rejectAuth(c, "Invalid token")
		return
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		rejectAuth(c, "Token expired")
		return
	}

	c.Set(contextUsername, claims.Username)
	c.Set(contextAuthMethod, authMethodJWT)
	c.Next()
}

func authenticatePersonalAccessToken(c *gin.Context, plaintext string) {
	if models.DB == nil {
		rejectAuth(c, "Invalid token")
		return
	}

	var pat models.PersonalAccessToken
	err := models.DB.Where("token_hash = ?", models.HashPersonalAccessToken(plaintext)).First(&pat).Error
	if err != nil {
		rejectAuth(c, "Invalid token")
		return
	}
	if pat.RevokedAt != nil {
		rejectAuth(c, "Token revoked")
		return
	}
	if pat.ExpiresAt != nil && !pat.ExpiresAt.After(time.Now()) {
		rejectAuth(c, "Token expired")
		return
	}

	var user models.BaseAccountsInfos
	if err := models.DB.First(&user, pat.UserID).Error; err != nil || user.Account == "" {
		rejectAuth(c, "Invalid token")
		return
	}
	if !utils.TokenScopeAllows(pat.Scopes, c.Request.URL.Path) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token scope does not allow this request"})
		c.Abort()
		return
	}

	now := time.Now()
	_ = models.DB.Model(&models.PersonalAccessToken{}).Where("id = ?", pat.ID).Update("last_used_at", now).Error

	c.Set(contextUsername, user.Account)
	c.Set(contextUserID, user.ID)
	c.Set(contextAuthMethod, authMethodPAT)
	c.Next()
}
