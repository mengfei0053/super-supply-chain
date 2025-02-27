package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"
	"path/filepath"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func GetDynamicExcelTableList(c *gin.Context) {
	query, err := utils.GetListQueryParams(c)
	tableName := c.Param("tableName")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var total int64

	res := []models.DynamicExcelTable{}

	sqlQuery := models.DB.Table(tableName).Count(&total)
	sqlQuery = sqlQuery.Limit(query.Limit).Offset(query.Offset).Find(&res)

	utils.SetContentRange(c, total)
	c.JSON(http.StatusOK, res)

}

func GetDynamicExcelTableDetail(c *gin.Context) {
	id := c.Param("id")
	tableName := c.Param("tableName")
	res := models.DynamicExcelTable{}

	sqlQuery := models.DB.Table(tableName).Where("id = ?", id).First(&res)
	if sqlQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": sqlQuery.Error})
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeleteDynamicExcelTable(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func UpdateDynamicExcelTable(c *gin.Context) {
	id := c.Param("id")
	tableName := c.Param("tableName")

	reqBody := models.DynamicExcelTable{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	sqlQuery := models.DB.Table(tableName).Where("id = ?", id).Updates(&reqBody)
	if sqlQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": sqlQuery.Error})
		return
	}

	c.JSON(http.StatusOK, reqBody)
}

func CreateDynamicExcelTable(c *gin.Context) {

	tableName := c.Param("tableName")
	mapRules, err := utils.GetMappingRules(tableName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}

	// Define the path where the file will be saved
	uploadDir := utils.GetUploadTmpDir()
	uuidFileName := uuid.New().String()
	extension := filepath.Ext(file.Filename)

	newFileName := uuidFileName + extension

	filePath := filepath.Join(uploadDir, newFileName)

	// Save the file to the specified directory
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	data, err := utils.GetExcelData(filePath, mapRules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileUrl, err := utils.UploadToNas(filePath, newFileName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to NAS"})
		return
	}
	query := models.DB.Table(tableName).Model(&models.DynamicExcelTable{}).Create(&models.DynamicExcelTable{
		UploadFilePath: fileUrl,
		FileName:       file.Filename,
		Datas:          data,
	})
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}

	c.JSON(http.StatusOK, data)
}
