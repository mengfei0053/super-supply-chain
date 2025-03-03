package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"super-supply-chain/configs"
	"super-supply-chain/controllers"
	"super-supply-chain/middleware"
	"super-supply-chain/models"
	"super-supply-chain/utils"
	"time"
)

func main() {

	configs.LoadConfigFile()
	gin.SetMode(gin.ReleaseMode)

	models.InitDB()

	r := gin.Default()
	logger := utils.InitLogger()

	r.Use(middleware.GinZapLogger(logger))
	// 替换默认 Recovery 中间件
	r.Use(middleware.GinZapRecovery(logger, true))

	controllers.LoadStatic(r)

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

		protected.GET("/yifan/cost-calculation", controllers.GetExcelMappingRules)
		protected.GET("/yifan/cost-calculation/:id", controllers.GetExcelMappingRuleDetail)
		protected.POST("/yifan/cost-calculation", controllers.CreateExcelMappingRules)
		protected.PUT("/yifan/cost-calculation/:id", controllers.UpdateExcelMappingRules)
		protected.DELETE("/yifan/cost-calculation/:id", controllers.DeleteExcelMappingRules)

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
		protected.GET("/excel-exports/:tableName", controllers.ExportDynamicExcel)
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

	srv := &http.Server{
		Addr:    ":" + configs.PORT,
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.String("address", srv.Addr), zap.Error(err))
		}
	}()
	logger.Info("服务器已启动", zap.String("address", srv.Addr))

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭异常", zap.String("reason", err.Error()))
	}
	logger.Info("服务器已关闭")
}
