package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

type ExcelReadRuleItem struct {
	ID               uint    `json:"id"`
	MenuName         string  `json:"menuName"`
	SheetIndex       uint    `json:"sheetIndex"`
	DynamicTableName string  `json:"dynamicTableName"`
	Desc             string  `json:"desc"`
	MapRule          MapRule `json:"mapRule"`
	IterateRule      MapRule `json:"iterateRules"`
}

// 获取Excel读取规则
func GetExcelReadRule(c *gin.Context) {
	id := c.Param("id")
	rule := models.ExcelReadRuleInfos{}
	models.DB.Model(&models.ExcelReadRuleInfos{}).Where("id = ?", id).First(&rule)

	c.JSON(http.StatusOK, rule)

}

type MapRule struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 获取Excel读取规则列表
func GetExcelReadRulesList(c *gin.Context) {
	query, err := utils.GetListQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询总条数
	var total int64
	res := []models.ExcelReadRuleInfos{}

	sqlQuery := models.DB.Model(&models.ExcelReadRuleInfos{}).Count(&total)
	sqlQuery = sqlQuery.Limit(query.Limit).Offset(query.Offset).Order("id").Find(&res)

	utils.SetContentRange(c, total)
	c.JSON(http.StatusOK, res)
}

// 创建Excel读取规则
func CreateExcelReadRules(c *gin.Context) {
	req := models.ExcelReadRuleInfos{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error

	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err := models.DB.Table(req.DynamicTableName).Migrator().CreateTable(&models.DynamicExcelTable{})
		if err != nil {
			return err
		}

		result := models.DB.Create(&req)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 更新Excel读取规则
func UpdateExcelReadRules(c *gin.Context) {
	body := models.ExcelReadRuleInfos{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Updates(&body)

	c.JSON(200, body)
}

// 删除Excel读取规则
func DeleteExcelReadRules(c *gin.Context) {
	id := c.Param("id")
	models.DB.Delete(&models.ExcelReadRuleInfos{}, id)

	c.JSON(200, gin.H{
		"message": "Create Excel read rule successfully",
	})

}
