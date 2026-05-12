package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"super-supply-chain/controllers"
	"super-supply-chain/middleware"
)

func signedToken(t *testing.T, username string, expiresAt time.Time) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &controllers.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})
	tokenString, err := token.SignedString(controllers.JwtKey)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}
	return tokenString
}

func setupProtectedAPIRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	protected := r.Group("/api/admin")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/settlement-form-entry", controllers.GetSettlementFormEntry)
		protected.GET("/settlement-form-entry/:id", controllers.GetSettlementFormEntryDetail)
		protected.POST("/settlement-form-entries", controllers.CreateSettlementFormEntry)
		protected.PUT("/settlement-form-entry/:id", controllers.UpdateSettlementFormEntry)
		protected.DELETE("/settlement-form-entry/:id", controllers.DeleteSettlementFormEntry)

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

	return r
}

func TestProtectedAPIRequiresAuthorization(t *testing.T) {
	router := setupProtectedAPIRouter()
	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/admin/settlement-form-entry"},
		{http.MethodGet, "/api/admin/settlement-form-entry/1"},
		{http.MethodPost, "/api/admin/settlement-form-entries"},
		{http.MethodPut, "/api/admin/settlement-form-entry/1"},
		{http.MethodDelete, "/api/admin/settlement-form-entry/1"},
		{http.MethodGet, "/api/admin/excel-read-rules"},
		{http.MethodGet, "/api/admin/excel-read-rules/1"},
		{http.MethodPost, "/api/admin/excel-read-rules"},
		{http.MethodPut, "/api/admin/excel-read-rules/1"},
		{http.MethodDelete, "/api/admin/excel-read-rules/1"},
		{http.MethodGet, "/api/admin/dict-manage"},
		{http.MethodGet, "/api/admin/dict-manage/1"},
		{http.MethodPost, "/api/admin/dict-manage"},
		{http.MethodPut, "/api/admin/dict-manage/1"},
		{http.MethodDelete, "/api/admin/dict-manage/1"},
		{http.MethodGet, "/api/admin/dict-manage/map/status"},
		{http.MethodGet, "/api/admin/excel/orders"},
		{http.MethodGet, "/api/admin/excel-exports/orders"},
		{http.MethodGet, "/api/admin/excel/orders/1"},
		{http.MethodPost, "/api/admin/excel/orders"},
		{http.MethodPut, "/api/admin/excel/orders/1"},
		{http.MethodDelete, "/api/admin/excel/orders/1"},
		{http.MethodGet, "/api/admin/excel-export-rule/template/orders"},
		{http.MethodGet, "/api/admin/excel-export-rule/template/orders/1"},
		{http.MethodPost, "/api/admin/excel-export-rule/template/orders"},
		{http.MethodPut, "/api/admin/excel-export-rule/orders/1"},
		{http.MethodDelete, "/api/admin/excel-export-rule/orders/1"},
		{http.MethodPost, "/api/admin/excel-export-rule/orders/export"},
		{http.MethodGet, "/api/admin/excel-export-rule/orders/export/1"},
		{http.MethodGet, "/api/admin/options/status"},
		{http.MethodGet, "/api/admin/menus"},
	}

	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, body = %s; want 401", w.Code, w.Body.String())
			}
		})
	}
}

func TestAuthMiddlewareAcceptsAuthorizationBearerHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/admin/ping", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"username": c.GetString("username")})
	})

	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/api/admin/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s; want 200", w.Code, w.Body.String())
	}
}

func TestAuthMiddlewareRejectsExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/admin/ping", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"username": c.GetString("username")})
	})

	token := signedToken(t, "alice", time.Now().Add(-time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/api/admin/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s; want 401", w.Code, w.Body.String())
	}
}
