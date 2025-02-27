package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func CreateExcelExportRuleTemplate(c *gin.Context) {

	tableName := c.Param("tableName")
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	alias := c.PostForm("alias")

	uploadDir := utils.GetUploadTmpDir()

	uuidFileName := uuid.New().String()
	extension := filepath.Ext(file.Filename)
	newFileName := uuidFileName + extension
	filePath := filepath.Join(uploadDir, newFileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	nasFileUrl, err := utils.UploadToNas(filePath, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to NAS"})
		return
	}

	data := models.ExcelExportTemplates{
		UploadFilePath:  nasFileUrl,
		FileName:        newFileName,
		AssociatedTable: tableName,
		Alias:           alias,
	}

	query := models.DB.Model(&models.ExcelExportTemplates{}).Create(&data)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetExcelExportRules(c *gin.Context) {
	query, err := utils.GetListQueryParams(c)
	tableName := c.Param("tableName")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var total int64
	res := []models.ExcelExportTemplates{}
	sqlQuery := models.DB.Model(&models.ExcelExportTemplates{}).Where("associated_table = ?", tableName).Count(&total)
	sqlQuery = sqlQuery.Limit(query.Limit).Offset(query.Offset).Find(&res)
	utils.SetContentRange(c, total)
	c.JSON(http.StatusOK, res)
}

func GetExcelExportRuleDetail(c *gin.Context) {
	id := c.Param("id")
	var res models.ExcelExportTemplates
	query := models.DB.Model(&models.ExcelExportTemplates{}).Where("id = ?", id).First(&res)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	c.JSON(http.StatusOK, res)
}

func UpdateExcelExportRule(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update excel export rule successfully"})
}

func DeleteExcelExportRule(c *gin.Context) {
	id := c.Param("id")
	query := models.DB.Where("id = ?", id).Delete(&models.ExcelExportTemplates{})
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Delete excel export rule successfully"})
}
