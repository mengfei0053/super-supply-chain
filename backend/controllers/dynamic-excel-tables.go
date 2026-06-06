package controllers

import (
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"super-supply-chain/models"
	"super-supply-chain/utils"
	excel_template_engines "super-supply-chain/utils/excel-template-engines"
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

	sqlQuery := models.DB.Table(tableName).Where("created_at between ? and ?",
		query.Filter.Start+" 00:00:00",
		query.Filter.End+" 23:59:59",
	).Count(&total)
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
	id := c.Param("id")
	tableName := c.Param("tableName")

	sqlQuery := models.DB.Unscoped().Table(tableName).Where("id = ?", id).Delete(&models.DynamicExcelTable{})
	if sqlQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": sqlQuery.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Delete successfully"})
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

	// Retrieve the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		panic(err)
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

	var data models.ExcelData

	if tableName == "dynamic_customs_declaration_form" {
		data, err = utils.GetBaoguanDan(filePath, tableName)
	} else {
		data, err = utils.GetExcelData(filePath, tableName)
	}
	if err != nil {
		panic(err)
		return
	}
	fileUrl, err := utils.UploadToNas(filePath, newFileName)

	if err != nil {
		panic(err)
		return
	}
	query := models.DB.Table(tableName).Model(&models.DynamicExcelTable{}).Create(&models.DynamicExcelTable{
		UploadFilePath: fileUrl,
		FileName:       file.Filename,
		Datas:          data,
		NasFileName:    newFileName,
	})
	if query.Error != nil {
		panic(query.Error)
		return
	}

	c.JSON(http.StatusOK, data)
	//c.JSON(http.StatusOK, gin.H{})
}

func ExportDynamicExcel(c *gin.Context) {
	ids := c.QueryArray("ids")
	queryType := c.Query("type")
	tableName := c.Param("tableName")

	filePath, err := excel_template_engines.GetExcelExportFilePath(tableName, ids, queryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", buildExcelExportContentDisposition(queryType))
	c.File(filePath)

}

func buildExcelExportContentDisposition(queryType string) string {
	fileName := buildExcelExportFileName(queryType)
	return "attachment; filename=\"export.xlsx\"; filename*=UTF-8''" + url.PathEscape(fileName)
}

func buildExcelExportFileName(queryType string) string {
	labels := map[string]string{
		"invoice_freight":                     "导出发票-运费",
		"invoice_clearance_only":             "导出发票-清关",
		"invoice_unpacking":                  "导出发票-掏箱",
		"invoice_clearance":                  "导出发票-清关-掏箱",
		"shortHaulInvoice":                   "导出-短驳费-发票",
		"shortHaul":                          "导出-短驳费表",
		"shortHaulAndFeiChang":               "导出-短驳费表(含分厂)",
		"dynamic_Integrity_packaging_invoice": "导出发票",
	}

	label, ok := labels[queryType]
	if !ok {
		label = "导出文件"
	}

	return label + "_" + time.Now().Format("2006_01_02_15_04_05") + ".xlsx"
}
