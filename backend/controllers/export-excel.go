package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

type ExportExcelReq struct {
	TemplateId int `json:"templateId"`
}

func ExportExcel(c *gin.Context) {
	tableName := c.Param("tableName")
	var req ExportExcelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templateFileInfo := models.ExcelExportTemplates{}
	query := models.DB.Model(models.ExcelExportTemplates{}).Where("id = ?", req.TemplateId).First(&templateFileInfo)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}
	downloadFile, err := utils.DownloadFromNas(templateFileInfo.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(tableName)

	fmt.Println(downloadFile)
	c.JSON(http.StatusOK, gin.H{"message": downloadFile})
}

func SingleExportExcel(c *gin.Context) {
	tableName := c.Param("tableName")
	id := c.Param("id")

	excelData := models.DynamicExcelTable{}

	query := models.DB.Table(tableName).Where("id = ?", id).First(&excelData)
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error})
		return
	}

	filePath, err := utils.CreateFile(&excelData.Datas, excelData.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=example.zip")
	c.File(filePath)
}
