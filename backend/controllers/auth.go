package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"super-supply-chain/configs"
	"super-supply-chain/models"
	"time"
)

var JwtKey = []byte("mengfei_super_supply_chain_admin")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Register(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.BaseAccountsInfos
	user.Account = creds.Username
	if err := user.SetPassword(creds.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set password"})
		return
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.BaseAccountsInfos
	if err := models.DB.Where("account = ?", creds.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	type UserResponse struct {
		ID       uint   `json:"id"`
		UserName string `json:"username"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		Token    string `json:"token"`
		Avatar   string `json:"avatar"`
	}

	userRes := UserResponse{
		ID:       user.ID,
		UserName: user.Account,
		FullName: user.Realname,
		Email:    user.Email,
		Token:    tokenString,
		Avatar:   user.Avatar,
	}

	Authorization := "Bearer " + tokenString

	c.SetCookie(configs.AuthKey, Authorization, 60*60*24, "/", configs.HOST, false, true)

	c.JSON(http.StatusOK, userRes)

	//c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
