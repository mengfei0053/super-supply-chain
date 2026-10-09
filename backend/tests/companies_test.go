package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"super-supply-chain/models"
)

func TestSearchCompaniesMatchesNameOrAliasAndSkipsDeleted(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	seedCompany(t, companySeed{
		Name:                    "南阳食品有限公司",
		Alias:                   "食品",
		AddrCountry:             "中国",
		AddrProvince:            "河南",
		AddrCity:                "南阳",
		AddrStreet:              "示例路 1 号",
		UnifiedSocialCreditCode: "91330000NAME",
		BankCode:                "BANK001",
		PhoneNum:                "0377-0000000",
		TargetAddr:              "南阳",
	})
	seedCompany(t, companySeed{
		Name:                    "某某物流有限公司",
		Alias:                   "南阳仓",
		UnifiedSocialCreditCode: "91330000ALIAS",
		TargetAddr:              "南阳",
	})
	seedCompany(t, companySeed{
		Name:                    "北京公司",
		Alias:                   "北京",
		UnifiedSocialCreditCode: "91330000OTHER",
		TargetAddr:              "北京",
	})
	seedCompany(t, companySeed{
		Name:                    "南阳已删除公司",
		Alias:                   "旧南阳",
		UnifiedSocialCreditCode: "91330000DEL",
		TargetAddr:              "南阳",
		DeletedAt:               "2020-01-02 03:04:05",
	})

	body := getCompanies(t, "南阳", false)
	if len(body) != 2 {
		t.Fatalf("active matches = %d, want 2; body = %s", len(body), mustJSON(t, body))
	}
	if body[0]["name"] != "南阳食品有限公司" || body[1]["alias"] != "南阳仓" {
		t.Fatalf("matches = %s", mustJSON(t, body))
	}

	keys := []string{
		"id", "created_at", "updated_at", "deleted_at", "name",
		"addr_country", "addr_province", "addr_city", "addr_street",
		"unified_social_credit_code", "bank_code", "phone_num", "alias", "target_addr",
	}
	for _, key := range keys {
		if _, ok := body[0][key]; !ok {
			t.Fatalf("first row missing %s: %s", key, mustJSON(t, body[0]))
		}
	}
	if body[0]["deleted_at"] != nil {
		t.Fatalf("active deleted_at = %#v, want null", body[0]["deleted_at"])
	}
	if body[0]["addr_city"] != "南阳" || body[0]["unified_social_credit_code"] != "91330000NAME" {
		t.Fatalf("full row = %s", mustJSON(t, body[0]))
	}

	withDeleted := getCompanies(t, "南阳", true)
	if len(withDeleted) != 3 {
		t.Fatalf("includeDeleted matches = %d, want 3; body = %s", len(withDeleted), mustJSON(t, withDeleted))
	}
	if withDeleted[2]["name"] != "南阳已删除公司" || withDeleted[2]["deleted_at"] == nil {
		t.Fatalf("deleted row = %s", mustJSON(t, withDeleted[2]))
	}
}

func TestSearchCompaniesTreatsLikeWildcardsAsLiterals(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	seedCompany(t, companySeed{
		Name:                    "北京公司",
		Alias:                   "北京",
		UnifiedSocialCreditCode: "91330000BJ",
	})
	seedCompany(t, companySeed{
		Name:                    "百分百公司",
		Alias:                   "100%",
		UnifiedSocialCreditCode: "91330000PCT",
	})

	if got := getCompanies(t, "%", false); len(got) != 1 || got[0]["name"] != "百分百公司" {
		t.Fatalf("literal percent matches = %s, want only 百分百公司", mustJSON(t, got))
	}
	if got := getCompanies(t, "100%", false); len(got) != 1 || got[0]["alias"] != "100%" {
		t.Fatalf("alias percent matches = %s", mustJSON(t, got))
	}
}

func TestSearchCompaniesRequiresKeyword(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	router := setupProtectedAPIRouter()
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/api/admin/companies", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s; want 400", w.Code, w.Body.String())
	}
}

type companySeed struct {
	Name                    string
	Alias                   string
	AddrCountry             string
	AddrProvince            string
	AddrCity                string
	AddrStreet              string
	UnifiedSocialCreditCode string
	BankCode                string
	PhoneNum                string
	TargetAddr              string
	DeletedAt               string
}

func seedCompany(t *testing.T, row companySeed) {
	t.Helper()
	err := models.DB.Exec(
		`INSERT INTO base_companies_infos (
			name, alias, addr_country, addr_province, addr_city, addr_street,
			unified_social_credit_code, bank_code, phone_num, target_addr, deleted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		row.Name,
		row.Alias,
		row.AddrCountry,
		row.AddrProvince,
		row.AddrCity,
		row.AddrStreet,
		row.UnifiedSocialCreditCode,
		row.BankCode,
		row.PhoneNum,
		row.TargetAddr,
		nullIfEmpty(row.DeletedAt),
	).Error
	if err != nil {
		t.Fatalf("seed company %s: %v", row.Name, err)
	}
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func getCompanies(t *testing.T, keyword string, includeDeleted bool) []map[string]any {
	t.Helper()
	router := setupProtectedAPIRouter()
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	path := "/api/admin/companies?keyword=" + url.QueryEscape(keyword)
	if includeDeleted {
		path += "&includeDeleted=true"
	}
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s; want 200", w.Code, w.Body.String())
	}
	var body []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not a JSON array: %v; body = %s", err, w.Body.String())
	}
	return body
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(raw)
}
