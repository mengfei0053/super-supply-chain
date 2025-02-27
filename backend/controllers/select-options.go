package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

type ExcelMappingRuleOptions struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetOptions(c *gin.Context) {
	key := c.Param("key")
	switch key {
	case "excel-mapping-rule":
		GetExcelMappingRuleOptions(c)
	case "export-templates":
		GetExportTemplates(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	}
}

func GetExportTemplates(c *gin.Context) {
	result := []models.ExcelExportTemplates{}
	associated_table := c.Query("associated_table")

	err := models.DB.Model(&models.ExcelExportTemplates{}).Where("associated_table = ?", associated_table).Select("id, alias").Find(&result)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, utils.Map(result, func(i models.ExcelExportTemplates) map[string]any {
		return map[string]any{
			"id":   i.ID,
			"name": i.Alias,
		}
	}))
}

func GetExcelMappingRuleOptions(c *gin.Context) {

	result := []ExcelMappingRuleOptions{}

	err := models.DB.Model(&models.ExcelMappingRules{}).Select("id, name").Find(&result)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}
