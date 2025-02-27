package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func GetExcelMappingRules(c *gin.Context) {

	filter := c.Query("filter")
	range_ := c.Query("range")
	sort := c.Query("sort")

	fmt.Print(filter, "filter")
	fmt.Print(range_, "range")
	fmt.Print(sort, "sort")

	var sliceRange []int
	err := json.Unmarshal([]byte(range_), &sliceRange)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal range"})
		return
	}
	fmt.Print(sliceRange[0], "sliceRange")
	fmt.Print(sliceRange[1], "sliceRange")

	limit := sliceRange[1] - sliceRange[0] + 1
	offset := sliceRange[0]

	var rules []models.ExcelMappingRules
	// 查询总条数
	var total int64

	query := models.DB.Model(&models.ExcelMappingRules{}).Count(&total)
	query = query.Limit(limit).Offset(offset).Order("id").Find(&rules)

	utils.SetContentRange(c, total)
	c.JSON(http.StatusOK, rules)
}

func GetExcelMappingRuleDetail(c *gin.Context) {
	var ruleItem models.ExcelMappingRules
	models.DB.Where("id = ?", c.Param("id")).First(&ruleItem)
	c.JSON(http.StatusOK, ruleItem)
}

func CreateExcelMappingRules(c *gin.Context) {
	// 从 request body 中获取数据
	var req models.ExcelMappingRules
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := models.DB.Create(&req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mapping rules"})
		return
	}

	c.JSON(http.StatusOK, req)
}
func UpdateExcelMappingRules(c *gin.Context) {
	var req models.ExcelMappingRules
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := models.DB.Updates(&req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mapping rules"})
		return
	}
	c.JSON(http.StatusOK, req)

}
func DeleteExcelMappingRules(c *gin.Context) {

	id := c.Param("id")
	result := models.DB.Delete(&models.ExcelMappingRules{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete mapping rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DeleteExcelMappingRules mapping rules successfully"})
}
