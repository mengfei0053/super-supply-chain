package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

type SettlementFormEntryItem struct {
	ID          uint   `json:"id"`
	OrderNumber string `json:"orderNumber"`
	ArrivalDate string `json:"arrivalDate"`
	ArrivalPort string `json:"arrivalPort"`
}

func GetSettlementFormEntry(c *gin.Context) {

	var orders []models.Order
	var res []SettlementFormEntryItem

	models.DB.Order("created_at desc").Find(&orders)

	for _, num := range orders {
		res = append(res, SettlementFormEntryItem{
			ID:          num.ID,
			OrderNumber: num.OrderNumber,
			ArrivalDate: num.ArrivalDate,
			ArrivalPort: num.ArrivalPort,
		})
	}

	utils.SetContentRange(c, 0)
	if res == nil {
		res = []SettlementFormEntryItem{}
	}
	c.JSON(http.StatusOK, res)
}

func GetSettlementFormEntryDetail(c *gin.Context) {
	var orders models.Order

	models.DB.Where("id = ?", c.Param("id")).First(&orders)

	res := SettlementFormEntryItem{
		ID:          orders.ID,
		OrderNumber: orders.OrderNumber,
		ArrivalDate: orders.ArrivalDate,
		ArrivalPort: orders.ArrivalPort,
	}

	c.JSON(http.StatusOK, res)
}

func CreateSettlementFormEntry(c *gin.Context) {
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

	fmt.Println(uploadDir)

	filePath := filepath.Join(uploadDir, newFileName)

	// Save the file to the specified directory
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	fileUrl, err := utils.UploadToNas(filePath, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to NAS"})
		return
	}

	models.DB.Save(&models.UploadFile{
		FileName:       newFileName,
		FilePath:       fileUrl,
		OriginFileName: file.Filename,
		FileSize:       file.Size,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Upload settlement form successfully", "file_path": filePath})

}

func UpdateSettlementFormEntry(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create settlement form entry successfully"})
}

func DeleteSettlementFormEntry(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create settlement form entry successfully"})
}
