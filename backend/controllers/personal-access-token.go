package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"super-supply-chain/models"
	"super-supply-chain/utils"
)

type personalAccessTokenView struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"userId"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Scopes     string     `json:"scopes"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	RevokedAt  *time.Time `json:"revokedAt"`
}

type personalAccessTokenCreated struct {
	personalAccessTokenView
	Token string `json:"token"`
}

type createPersonalAccessTokenBody struct {
	Name      string `json:"name"`
	Scopes    string `json:"scopes"`
	ExpiresAt string `json:"expiresAt"`
}

func tokenView(row models.PersonalAccessToken) personalAccessTokenView {
	return personalAccessTokenView{
		ID:         row.ID,
		UserID:     row.UserID,
		Name:       row.Name,
		Prefix:     row.Prefix,
		Scopes:     row.Scopes,
		CreatedAt:  row.CreatedAt,
		LastUsedAt: row.LastUsedAt,
		ExpiresAt:  row.ExpiresAt,
		RevokedAt:  row.RevokedAt,
	}
}

func writeAPIError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message, "message": message})
}

func currentAccount(c *gin.Context) (models.BaseAccountsInfos, bool) {
	username := c.GetString("username")
	var user models.BaseAccountsInfos
	if username == "" || models.DB == nil {
		writeAPIError(c, http.StatusUnauthorized, "Unauthorized")
		return user, false
	}
	if err := models.DB.Where("account = ?", username).First(&user).Error; err != nil {
		writeAPIError(c, http.StatusUnauthorized, "Unauthorized")
		return user, false
	}
	return user, true
}

func ListPersonalAccessTokens(c *gin.Context) {
	user, ok := currentAccount(c)
	if !ok {
		return
	}
	query, err := utils.GetListQueryParams(c)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, err.Error())
		return
	}

	var total int64
	if err := models.DB.Model(&models.PersonalAccessToken{}).
		Where("user_id = ? AND revoked_at IS NULL", user.ID).
		Count(&total).Error; err != nil {
		writeAPIError(c, http.StatusInternalServerError, "Failed to list tokens")
		return
	}

	rows := []models.PersonalAccessToken{}
	if err := models.DB.
		Where("user_id = ? AND revoked_at IS NULL", user.ID).
		Order("id desc").
		Limit(query.Limit).
		Offset(query.Offset).
		Find(&rows).Error; err != nil {
		writeAPIError(c, http.StatusInternalServerError, "Failed to list tokens")
		return
	}

	views := make([]personalAccessTokenView, 0, len(rows))
	for _, row := range rows {
		views = append(views, tokenView(row))
	}
	utils.SetContentRange(c, total)
	c.JSON(http.StatusOK, views)
}

func CreatePersonalAccessToken(c *gin.Context) {
	user, ok := currentAccount(c)
	if !ok {
		return
	}

	var body createPersonalAccessTokenBody
	if err := c.ShouldBindJSON(&body); err != nil {
		writeAPIError(c, http.StatusBadRequest, "Invalid request")
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" || len([]rune(name)) > 100 {
		writeAPIError(c, http.StatusBadRequest, "name is required and must be at most 100 characters")
		return
	}
	scopes, err := utils.NormalizeTokenScopes(body.Scopes)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, err.Error())
		return
	}
	expiresAt, err := parseOptionalExpiry(body.ExpiresAt)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, err.Error())
		return
	}

	plaintext, hash, prefix, err := models.NewPersonalAccessTokenSecret()
	if err != nil {
		writeAPIError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	record := models.PersonalAccessToken{
		UserID:    user.ID,
		Name:      name,
		TokenHash: hash,
		Prefix:    prefix,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	}
	if err := models.DB.Create(&record).Error; err != nil {
		writeAPIError(c, http.StatusInternalServerError, "Failed to create token")
		return
	}

	// plaintext is returned only in this response. Do not log record or plaintext.
	c.JSON(http.StatusOK, personalAccessTokenCreated{
		personalAccessTokenView: tokenView(record),
		Token:                   plaintext,
	})
}

func RevokePersonalAccessToken(c *gin.Context) {
	user, ok := currentAccount(c)
	if !ok {
		return
	}

	var record models.PersonalAccessToken
	err := models.DB.
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", c.Param("id"), user.ID).
		First(&record).Error
	if err != nil {
		writeAPIError(c, http.StatusNotFound, "Token not found")
		return
	}

	now := time.Now()
	if err := models.DB.Model(&record).Update("revoked_at", now).Error; err != nil {
		writeAPIError(c, http.StatusInternalServerError, "Failed to revoke token")
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": record.ID})
}

func parseOptionalExpiry(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		if !parsed.After(time.Now()) {
			return nil, errors.New("expiresAt must be in the future")
		}
		return &parsed, nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		return nil, errors.New("invalid expiresAt")
	}
	end := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 0, time.Local)
	if !end.After(time.Now()) {
		return nil, errors.New("expiresAt must be in the future")
	}
	return &end, nil
}
