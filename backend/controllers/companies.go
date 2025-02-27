package controllers

import "github.com/gin-gonic/gin"

func GetCompanies(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(200, gin.H{"username": username})
}
