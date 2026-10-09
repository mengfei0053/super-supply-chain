package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"super-supply-chain/controllers"
	"super-supply-chain/middleware"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func setupPATRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CaptureAndRedactQueryToken())
	protected := router.Group("/api/admin")
	protected.Use(middleware.AuthMiddleware())
	sessionOnly := middleware.RequireSessionJWT()
	protected.GET("/personal-access-tokens", sessionOnly, controllers.ListPersonalAccessTokens)
	protected.POST("/personal-access-tokens", sessionOnly, controllers.CreatePersonalAccessToken)
	protected.DELETE("/personal-access-tokens/:id", sessionOnly, controllers.RevokePersonalAccessToken)
	protected.GET("/dict-manage", controllers.GetDicts)
	protected.GET("/excel/ping", func(c *gin.Context) {
		if strings.Contains(c.Request.URL.RawQuery, "ssc_pat_") || strings.Contains(c.Request.RequestURI, "ssc_pat_") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token leaked into request URL"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"username": c.GetString("username"), "authMethod": c.GetString("authMethod")})
	})
	protected.DELETE("/excel/orders/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
	})
	return router
}

func createAccount(t *testing.T, account string) models.BaseAccountsInfos {
	t.Helper()
	user := models.BaseAccountsInfos{Account: account, Realname: account}
	if err := user.SetPassword("secret"); err != nil {
		t.Fatalf("set password: %v", err)
	}
	if err := models.DB.Create(&user).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}
	return user
}

func doJSON(router http.Handler, method, path, token string, body any, extraHeader map[string]string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for key, value := range extraHeader {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func createToken(t *testing.T, router http.Handler, jwt string, body map[string]string) (map[string]any, string) {
	t.Helper()
	w := doJSON(router, http.MethodPost, "/api/admin/personal-access-tokens", jwt, body, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("create status = %d, body = %s", w.Code, w.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("create body: %v", err)
	}
	plaintext, _ := payload["token"].(string)
	if plaintext == "" {
		t.Fatalf("create response did not include token: %s", w.Body.String())
	}
	if _, ok := payload["tokenHash"]; ok {
		t.Fatal("create response included tokenHash")
	}
	if _, ok := payload["token_hash"]; ok {
		t.Fatal("create response included token_hash")
	}
	if strings.Contains(w.Body.String(), models.HashPersonalAccessToken(plaintext)) {
		t.Fatal("create response included the token hash")
	}
	return payload, plaintext
}

func TestPersonalAccessTokenCreateListAndHash(t *testing.T) {
	setupTestDB(t, &models.BaseAccountsInfos{}, &models.PersonalAccessToken{})
	alice := createAccount(t, "alice")
	router := setupPATRouter()
	jwt := signedToken(t, "alice", time.Now().Add(time.Hour))

	created, plaintext := createToken(t, router, jwt, map[string]string{
		"name":      "cursor-mcp",
		"scopes":    "excel，dict",
		"expiresAt": "2099-12-31",
	})
	if !models.IsPersonalAccessToken(plaintext) {
		t.Fatalf("plaintext %q is not a PAT", plaintext)
	}
	if created["name"] != "cursor-mcp" {
		t.Fatalf("name = %#v", created["name"])
	}
	if created["scopes"] != "excel,dict" {
		t.Fatalf("scopes = %#v", created["scopes"])
	}
	prefix, _ := created["prefix"].(string)
	if prefix == "" || !strings.HasPrefix(plaintext, prefix) || prefix == plaintext {
		t.Fatalf("prefix = %q, plaintext = %q", prefix, plaintext)
	}

	var stored models.PersonalAccessToken
	if err := models.DB.Where("user_id = ?", alice.ID).First(&stored).Error; err != nil {
		t.Fatalf("stored token: %v", err)
	}
	if stored.TokenHash == plaintext || strings.Contains(stored.TokenHash, "ssc_pat_") {
		t.Fatalf("plaintext stored in token_hash: %s", stored.TokenHash)
	}
	if stored.TokenHash != models.HashPersonalAccessToken(plaintext) {
		t.Fatalf("hash = %s", stored.TokenHash)
	}

	list := doJSON(router, http.MethodGet, "/api/admin/personal-access-tokens?range=%5B0,10%5D", jwt, nil, nil)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", list.Code, list.Body.String())
	}
	if list.Header().Get("Content-Range") != "1" {
		t.Fatalf("Content-Range = %q", list.Header().Get("Content-Range"))
	}
	if strings.Contains(list.Body.String(), plaintext) || strings.Contains(list.Body.String(), stored.TokenHash) {
		t.Fatalf("list revealed secret material: %s", list.Body.String())
	}
	var rows []map[string]any
	if err := json.Unmarshal(list.Body.Bytes(), &rows); err != nil {
		t.Fatalf("list json: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("list len = %d", len(rows))
	}
	if _, ok := rows[0]["token"]; ok {
		t.Fatal("list row included token")
	}
	if rows[0]["prefix"] != prefix {
		t.Fatalf("list prefix = %#v", rows[0]["prefix"])
	}
}

func TestPersonalAccessTokenRejectsBadInput(t *testing.T) {
	setupTestDB(t, &models.BaseAccountsInfos{}, &models.PersonalAccessToken{})
	createAccount(t, "alice")
	router := setupPATRouter()
	jwt := signedToken(t, "alice", time.Now().Add(time.Hour))

	missing := doJSON(router, http.MethodPost, "/api/admin/personal-access-tokens", jwt, map[string]string{}, nil)
	if missing.Code != http.StatusBadRequest {
		t.Fatalf("missing name status = %d", missing.Code)
	}
	unknown := doJSON(router, http.MethodPost, "/api/admin/personal-access-tokens", jwt, map[string]string{
		"name":   "bad",
		"scopes": "admin",
	}, nil)
	if unknown.Code != http.StatusBadRequest {
		t.Fatalf("unknown scope status = %d, body = %s", unknown.Code, unknown.Body.String())
	}
	past := doJSON(router, http.MethodPost, "/api/admin/personal-access-tokens", jwt, map[string]string{
		"name":      "old",
		"expiresAt": "2000-01-01",
	}, nil)
	if past.Code != http.StatusBadRequest {
		t.Fatalf("past expiry status = %d, body = %s", past.Code, past.Body.String())
	}
}

func TestPersonalAccessTokenAuthMethodsAndRevoke(t *testing.T) {
	setupTestDB(t, &models.BaseAccountsInfos{}, &models.PersonalAccessToken{}, &models.BaseDict{})
	createAccount(t, "alice")
	createAccount(t, "bob")
	router := setupPATRouter()
	aliceJWT := signedToken(t, "alice", time.Now().Add(time.Hour))
	bobJWT := signedToken(t, "bob", time.Now().Add(time.Hour))

	created, plaintext := createToken(t, router, aliceJWT, map[string]string{"name": "full"})
	tokenID, _ := created["id"].(float64)
	tokenPath := fmt.Sprintf("/api/admin/personal-access-tokens/%.0f", tokenID)
	models.DB.Create(&models.BaseDict{Key: "enabled", Value: "启用", Type: "status"})

	header := doJSON(router, http.MethodGet, "/api/admin/dict-manage?range=%5B0,10%5D", plaintext, nil, nil)
	if header.Code != http.StatusOK {
		t.Fatalf("bearer PAT status = %d, body = %s", header.Code, header.Body.String())
	}

	apiKey := httptest.NewRequest(http.MethodGet, "/api/admin/excel/ping", nil)
	apiKey.Header.Set("X-API-Key", plaintext)
	apiKeyW := httptest.NewRecorder()
	router.ServeHTTP(apiKeyW, apiKey)
	if apiKeyW.Code != http.StatusOK || !strings.Contains(apiKeyW.Body.String(), `"username":"alice"`) {
		t.Fatalf("X-API-Key status = %d, body = %s", apiKeyW.Code, apiKeyW.Body.String())
	}

	query := httptest.NewRequest(http.MethodGet, "/api/admin/excel/ping?ids=1&token="+plaintext, nil)
	queryW := httptest.NewRecorder()
	router.ServeHTTP(queryW, query)
	if queryW.Code != http.StatusOK || strings.Contains(queryW.Body.String(), plaintext) {
		t.Fatalf("query token status = %d, body = %s", queryW.Code, queryW.Body.String())
	}

	var used models.PersonalAccessToken
	if err := models.DB.First(&used).Error; err != nil {
		t.Fatal(err)
	}
	if used.LastUsedAt == nil {
		t.Fatal("last_used_at was not set")
	}

	// A PAT must not be able to mint another token.
	mint := doJSON(router, http.MethodPost, "/api/admin/personal-access-tokens", plaintext, map[string]string{"name": "nested"}, nil)
	if mint.Code != http.StatusForbidden {
		t.Fatalf("PAT create status = %d, body = %s", mint.Code, mint.Body.String())
	}
	var count int64
	models.DB.Model(&models.PersonalAccessToken{}).Count(&count)
	if count != 1 {
		t.Fatalf("token count = %d, want 1", count)
	}

	// Another user cannot revoke it.
	foreign := doJSON(router, http.MethodDelete, tokenPath, bobJWT, nil, nil)
	if foreign.Code != http.StatusNotFound {
		t.Fatalf("foreign revoke status = %d, body = %s", foreign.Code, foreign.Body.String())
	}

	revoked := doJSON(router, http.MethodDelete, tokenPath, aliceJWT, nil, nil)
	if revoked.Code != http.StatusOK {
		t.Fatalf("revoke status = %d, body = %s", revoked.Code, revoked.Body.String())
	}
	after := doJSON(router, http.MethodGet, "/api/admin/excel/ping", plaintext, nil, nil)
	if after.Code != http.StatusUnauthorized || !strings.Contains(after.Body.String(), "Token revoked") {
		t.Fatalf("revoked token status = %d, body = %s", after.Code, after.Body.String())
	}
	listed := doJSON(router, http.MethodGet, "/api/admin/personal-access-tokens?range=%5B0,10%5D", aliceJWT, nil, nil)
	if strings.Contains(listed.Body.String(), plaintext) || listed.Header().Get("Content-Range") != "0" {
		t.Fatalf("revoked token still listed: %s range=%s", listed.Body.String(), listed.Header().Get("Content-Range"))
	}
}

func TestPersonalAccessTokenExpiryScopeAndJWTCoexistence(t *testing.T) {
	setupTestDB(t, &models.BaseAccountsInfos{}, &models.PersonalAccessToken{})
	createAccount(t, "alice")
	createAccount(t, "bob")
	router := setupPATRouter()
	aliceJWT := signedToken(t, "alice", time.Now().Add(time.Hour))
	bobJWT := signedToken(t, "bob", time.Now().Add(time.Hour))

	_, excelToken := createToken(t, router, aliceJWT, map[string]string{"name": "excel", "scopes": "excel"})
	allowed := doJSON(router, http.MethodDelete, "/api/admin/excel/orders/9", excelToken, nil, nil)
	if allowed.Code != http.StatusOK {
		t.Fatalf("excel scope status = %d, body = %s", allowed.Code, allowed.Body.String())
	}
	denied := doJSON(router, http.MethodGet, "/api/admin/dict-manage?range=%5B0,10%5D", excelToken, nil, nil)
	if denied.Code != http.StatusForbidden {
		t.Fatalf("dict with excel scope status = %d, body = %s", denied.Code, denied.Body.String())
	}

	_, expiring := createToken(t, router, aliceJWT, map[string]string{"name": "soon"})
	past := time.Now().Add(-time.Minute)
	if err := models.DB.Model(&models.PersonalAccessToken{}).Where("name = ?", "soon").Update("expires_at", past).Error; err != nil {
		t.Fatal(err)
	}
	expired := doJSON(router, http.MethodGet, "/api/admin/excel/ping", expiring, nil, nil)
	if expired.Code != http.StatusUnauthorized || !strings.Contains(expired.Body.String(), "Token expired") {
		t.Fatalf("expired PAT status = %d, body = %s", expired.Code, expired.Body.String())
	}

	// Existing session JWTs still work, and a cookie still wins over a different JWT header.
	cookieReq := httptest.NewRequest(http.MethodGet, "/api/admin/excel/ping", nil)
	cookieReq.AddCookie(&http.Cookie{Name: "Authorization", Value: "Bearer " + aliceJWT})
	cookieReq.Header.Set("Authorization", "Bearer "+bobJWT)
	cookieW := httptest.NewRecorder()
	router.ServeHTTP(cookieW, cookieReq)
	if cookieW.Code != http.StatusOK || !strings.Contains(cookieW.Body.String(), `"username":"alice"`) {
		t.Fatalf("cookie jwt status = %d, body = %s", cookieW.Code, cookieW.Body.String())
	}

	// An API client sending a PAT is not blocked by a leftover session cookie.
	patReq := httptest.NewRequest(http.MethodGet, "/api/admin/excel/ping", nil)
	patReq.AddCookie(&http.Cookie{Name: "Authorization", Value: "Bearer " + aliceJWT})
	patReq.Header.Set("Authorization", "Bearer "+excelToken)
	patW := httptest.NewRecorder()
	router.ServeHTTP(patW, patReq)
	if patW.Code != http.StatusOK || !strings.Contains(patW.Body.String(), `"authMethod":"pat"`) {
		t.Fatalf("PAT over cookie status = %d, body = %s", patW.Code, patW.Body.String())
	}

	all, err := utils.NormalizeTokenScopes("all, excel")
	if err != nil || all != "*" {
		t.Fatalf("normalize all = %q, err = %v", all, err)
	}
	if !utils.TokenScopeAllows("excel", "/api/admin/excel-exports/orders") {
		t.Fatal("excel scope should allow exports")
	}
	if utils.TokenScopeAllows("dict", "/api/admin/excel/orders") {
		t.Fatal("dict scope should not allow excel")
	}
	if !utils.TokenScopeAllows("", "/api/admin/personal-access-tokens") {
		t.Fatal("empty scope should allow every path")
	}
}

func TestLoggerRedactsTokenMaterial(t *testing.T) {
	if got := middleware.RedactSensitiveQuery("ids=1&token=ssc_pat_supersecret&type=a"); strings.Contains(got, "supersecret") || strings.Contains(got, "ssc_pat_") {
		t.Fatalf("query leaked: %s", got)
	}

	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	router := gin.New()
	router.Use(middleware.GinZapLogger(logger))
	router.GET("/api/admin/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/api/admin/ping?ids=1&token=ssc_pat_supersecret", nil)
	router.ServeHTTP(httptest.NewRecorder(), req)
	if strings.Contains(observerText(recorded.All()), "supersecret") {
		t.Fatalf("access log leaked token: %s", observerText(recorded.All()))
	}

	errorCore, errorLogs := observer.New(zapcore.ErrorLevel)
	recovery := gin.New()
	recovery.Use(middleware.GinZapRecovery(zap.New(errorCore), true))
	recovery.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})
	panicReq := httptest.NewRequest(http.MethodGet, "/panic?token=ssc_pat_supersecret", nil)
	panicReq.Header.Set("Authorization", "Bearer ssc_pat_supersecret")
	panicReq.Header.Set("X-API-Key", "ssc_pat_supersecret")
	panicReq.AddCookie(&http.Cookie{Name: "Authorization", Value: "Bearer ssc_pat_supersecret"})
	recovery.ServeHTTP(httptest.NewRecorder(), panicReq)
	if strings.Contains(observerText(errorLogs.All()), "supersecret") {
		t.Fatalf("panic log leaked token: %s", observerText(errorLogs.All()))
	}
}

func observerText(entries []observer.LoggedEntry) string {
	var b strings.Builder
	for _, entry := range entries {
		b.WriteString(entry.Message)
		for _, field := range entry.Context {
			b.WriteString(field.Key)
			b.WriteString("=")
			b.WriteString(field.String)
			b.WriteString(" ")
		}
	}
	return b.String()
}
