package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"super-supply-chain/configs"
	"super-supply-chain/controllers"
	"super-supply-chain/middleware"
	"super-supply-chain/models"
	"time"
)

func main() {

	configs.LoadConfigFile()

	models.InitDB()

	r := gin.Default()

	// Custom CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CORS_ALLOW_ORIGIN")}, // Update with your allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "range"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count", "Content-Range", "Authorization", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	protected := r.Group(("/api/admin"))
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/settlement-form-entry", controllers.GetSettlementFormEntry)
		protected.GET("/settlement-form-entry/:id", controllers.GetSettlementFormEntryDetail)
		protected.POST("/settlement-form-entries", controllers.CreateSettlementFormEntry)
		protected.PUT("/settlement-form-entry/:id", controllers.UpdateSettlementFormEntry)
		protected.DELETE("/settlement-form-entry/:id", controllers.DeleteSettlementFormEntry)

		protected.GET("/excel-mapping-rule", controllers.GetExcelMappingRules)
		protected.GET("/excel-mapping-rule/:id", controllers.GetExcelMappingRuleDetail)
		protected.POST("/excel-mapping-rule", controllers.CreateExcelMappingRules)
		protected.PUT("/excel-mapping-rule/:id", controllers.UpdateExcelMappingRules)
		protected.DELETE("/excel-mapping-rule/:id", controllers.DeleteExcelMappingRules)

		protected.GET("/excel-read-rules", controllers.GetExcelReadRulesList)
		protected.GET("/excel-read-rules/:id", controllers.GetExcelReadRule)
		protected.POST("/excel-read-rules", controllers.CreateExcelReadRules)
		protected.PUT("/excel-read-rules/:id", controllers.UpdateExcelReadRules)
		protected.DELETE("/excel-read-rules/:id", controllers.DeleteExcelReadRules)

		protected.GET("/dict-manage", controllers.GetDicts)
		protected.GET("/dict-manage/:id", controllers.GetDictDeltail)
		protected.POST("/dict-manage", controllers.CreateDict)
		protected.PUT("/dict-manage/:id", controllers.UpdateDict)
		protected.DELETE("/dict-manage/:id", controllers.DeleteDict)
		protected.GET("/dict-manage/map/:type", controllers.GetDictMap)

		protected.GET("/excel/:tableName", controllers.GetDynamicExcelTableList)
		protected.GET("/excel/:tableName/:id", controllers.GetDynamicExcelTableDetail)
		protected.POST("/excel/:tableName", controllers.CreateDynamicExcelTable)
		protected.PUT("/excel/:tableName/:id", controllers.UpdateDynamicExcelTable)
		protected.DELETE("/excel/:tableName/:id", controllers.DeleteDynamicExcelTable)

		protected.GET("/excel-export-rule/template/:tableName", controllers.GetExcelExportRules)
		protected.GET("/excel-export-rule/template/:tableName/:id", controllers.GetExcelExportRuleDetail)
		protected.POST("/excel-export-rule/template/:tableName", controllers.CreateExcelExportRuleTemplate)
		protected.PUT("/excel-export-rule/:tableName/:id", controllers.UpdateExcelExportRule)
		protected.DELETE("/excel-export-rule/:tableName/:id", controllers.DeleteExcelExportRule)
		protected.POST("/excel-export-rule/:tableName/export", controllers.ExportExcel)
		protected.GET("/excel-export-rule/:tableName/export/:id", controllers.SingleExportExcel)

		protected.GET("/options/:key", controllers.GetOptions)
		protected.GET("/menus", controllers.GetDynamicExcelMenus)

	}

	r.Run(":" + configs.PORT)
}
