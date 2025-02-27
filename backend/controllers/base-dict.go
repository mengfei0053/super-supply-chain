package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func CreateDict(c *gin.Context) {
	dict := models.BaseDict{}
	if err := c.ShouldBindJSON(&dict); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := models.DB.Create(&dict)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	c.JSON(http.StatusOK, dict)
}

func GetDicts(c *gin.Context) {
	query, err := utils.GetListQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var total int64
	res := []models.BaseDict{}
	sqlQuery := models.DB.Model(&models.BaseDict{}).Count(&total)
	sqlQuery = sqlQuery.Limit(query.Limit).Offset(query.Offset).Find(&res)
	utils.SetContentRange(c, total)
	c.JSON(http.StatusOK, res)
}

func GetDictDeltail(c *gin.Context) {
	id := c.Param("id")
	dict := models.BaseDict{}
	query := models.DB.Where("id = ?", id).First(&dict)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	c.JSON(http.StatusOK, dict)
}

func UpdateDict(c *gin.Context) {
	id := c.Param("id")
	dict := models.BaseDict{}
	if err := c.ShouldBindJSON(&dict); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := models.DB.Model(&models.BaseDict{}).Where("id = ?", id).Updates(&dict)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	c.JSON(http.StatusOK, dict)

}

func DeleteDict(c *gin.Context) {
	id := c.Param("id")
	query := models.DB.Delete(&models.BaseDict{}, id)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Delete dict successfully"})
}

func GetDictMap(c *gin.Context) {
	typeParam := c.Param("type")
	dicts := []models.BaseDict{}
	query := models.DB.Where("type = ?", typeParam).Find(&dicts)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	dictMap := make(map[string]string)
	for _, dict := range dicts {
		dictMap[dict.Key] = dict.Value
	}
	c.JSON(http.StatusOK, dictMap)
}
