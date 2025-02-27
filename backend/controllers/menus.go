package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"super-supply-chain/models"
)

type DynamicExcelMenu struct {
	ID               uint   `json:"id"`
	MenuName         string `json:"menuName"`
	DynamicTableName string `json:"dynamicTableName"`
}

func GetDynamicExcelMenus(c *gin.Context) {
	result := []DynamicExcelMenu{}

	err := models.DB.Model(&models.ExcelReadRules{}).Select("id, menu_name, dynamic_table_name").Find(&result)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}
