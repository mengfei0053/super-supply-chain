package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
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

func TestCompaniesIndexWithoutKeywordRequiresRange(t *testing.T) {
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
	if !strings.Contains(w.Body.String(), "range") {
		t.Fatalf("body = %s; want range error", w.Body.String())
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

func TestCreateCompanyInsertsSearchableRow(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))

	w := doCompanyJSON(t, http.MethodPost, "/api/admin/companies", map[string]any{
		"name":                       "  杭州测试公司  ",
		"alias":                      " 杭州测试 ",
		"target_addr":                "杭州",
		"unified_social_credit_code": "91330100TEST",
	}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s; want 200", w.Code, w.Body.String())
	}
	created := decodeCompanyObject(t, w.Body.Bytes())
	if created["name"] != "杭州测试公司" || created["alias"] != "杭州测试" || created["target_addr"] != "杭州" {
		t.Fatalf("created = %s", w.Body.String())
	}
	if created["unified_social_credit_code"] != "91330100TEST" || created["deleted_at"] != nil {
		t.Fatalf("created = %s", w.Body.String())
	}
	if created["addr_city"] != "" || created["bank_code"] != "" || created["phone_num"] != "" {
		t.Fatalf("unexpected extra fields = %s", w.Body.String())
	}

	byName := getCompanies(t, "杭州测试公司", false)
	if len(byName) != 1 || byName[0]["alias"] != "杭州测试" {
		t.Fatalf("search by name = %s", mustJSON(t, byName))
	}
	byAlias := getCompanies(t, "杭州测试", false)
	if len(byAlias) != 1 || byAlias[0]["name"] != "杭州测试公司" {
		t.Fatalf("search by alias = %s", mustJSON(t, byAlias))
	}
}

func TestCreateCompanyRequiresNameAndCreditCode(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))

	missingName := doCompanyJSON(t, http.MethodPost, "/api/admin/companies", map[string]any{
		"alias":                      "别名",
		"target_addr":                "杭州",
		"unified_social_credit_code": "91330100TEST",
	}, token)
	if missingName.Code != http.StatusBadRequest || !strings.Contains(missingName.Body.String(), "name is required") {
		t.Fatalf("missing name = %d %s", missingName.Code, missingName.Body.String())
	}

	missingCode := doCompanyJSON(t, http.MethodPost, "/api/admin/companies", map[string]any{
		"name":        "杭州测试公司",
		"alias":       "杭州测试",
		"target_addr": "杭州",
	}, token)
	if missingCode.Code != http.StatusBadRequest || !strings.Contains(missingCode.Body.String(), "unified_social_credit_code is required") {
		t.Fatalf("missing code = %d %s", missingCode.Code, missingCode.Body.String())
	}

	if got := getCompanies(t, "杭州", false); len(got) != 0 {
		t.Fatalf("failed creates were stored: %s", mustJSON(t, got))
	}
}

func TestCreateCompanyRejectsDuplicateNameOrCreditCode(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	seedCompany(t, companySeed{
		Name:                    "南阳已删除公司",
		Alias:                   "旧南阳",
		UnifiedSocialCreditCode: "91330000DEL",
		TargetAddr:              "南阳",
		DeletedAt:               "2020-01-02 03:04:05",
	})

	dupName := doCompanyJSON(t, http.MethodPost, "/api/admin/companies", map[string]any{
		"name":                       "南阳已删除公司",
		"unified_social_credit_code": "91330000NEW",
	}, token)
	if dupName.Code != http.StatusConflict || !strings.Contains(dupName.Body.String(), "name already exists") {
		t.Fatalf("duplicate name = %d %s", dupName.Code, dupName.Body.String())
	}

	dupCode := doCompanyJSON(t, http.MethodPost, "/api/admin/companies", map[string]any{
		"name":                       "新公司",
		"unified_social_credit_code": "91330000DEL",
	}, token)
	if dupCode.Code != http.StatusConflict || !strings.Contains(dupCode.Body.String(), "unified_social_credit_code already exists") {
		t.Fatalf("duplicate code = %d %s", dupCode.Code, dupCode.Body.String())
	}
}

func TestUpdateCompanyChangesBasicInfoAndIgnoresCreditCode(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	seedCompany(t, companySeed{
		Name:                    "杭州测试公司",
		Alias:                   "仓储备注",
		AddrCity:                "杭州",
		UnifiedSocialCreditCode: "91330100TEST",
		BankCode:                "BANK001",
		PhoneNum:                "0571-0000000",
		TargetAddr:              "杭州",
	})
	found := getCompanies(t, "杭州测试公司", false)
	if len(found) != 1 {
		t.Fatalf("seed search = %s", mustJSON(t, found))
	}
	id := companyIDString(t, found[0]["id"])

	partial := doCompanyJSON(t, http.MethodPut, "/api/admin/companies/"+id, map[string]any{
		"target_addr": "海宁",
	}, token)
	if partial.Code != http.StatusOK {
		t.Fatalf("partial status = %d, body = %s", partial.Code, partial.Body.String())
	}
	partialRow := decodeCompanyObject(t, partial.Body.Bytes())
	if partialRow["name"] != "杭州测试公司" || partialRow["alias"] != "仓储备注" || partialRow["target_addr"] != "海宁" {
		t.Fatalf("partial row = %s", partial.Body.String())
	}
	if partialRow["unified_social_credit_code"] != "91330100TEST" || partialRow["addr_city"] != "杭州" || partialRow["bank_code"] != "BANK001" {
		t.Fatalf("partial row changed other fields: %s", partial.Body.String())
	}

	full := doCompanyJSON(t, http.MethodPut, "/api/admin/companies/"+id, map[string]any{
		"name":                       "杭州测试公司（更新）",
		"alias":                      "",
		"target_addr":                "嘉兴",
		"unified_social_credit_code": "SHOULD-IGNORE",
		"addr_city":                  "嘉兴市",
		"phone_num":                  "0573-1111111",
	}, token)
	if full.Code != http.StatusOK {
		t.Fatalf("update status = %d, body = %s", full.Code, full.Body.String())
	}
	updated := decodeCompanyObject(t, full.Body.Bytes())
	if updated["name"] != "杭州测试公司（更新）" || updated["alias"] != "" || updated["target_addr"] != "嘉兴" {
		t.Fatalf("updated = %s", full.Body.String())
	}
	if updated["unified_social_credit_code"] != "91330100TEST" {
		t.Fatalf("credit code changed: %s", full.Body.String())
	}
	if updated["addr_city"] != "嘉兴市" || updated["phone_num"] != "0573-1111111" || updated["bank_code"] != "BANK001" {
		t.Fatalf("updated basic fields = %s", full.Body.String())
	}

	if got := getCompanies(t, "仓储备注", false); len(got) != 0 {
		t.Fatalf("cleared alias still matches: %s", mustJSON(t, got))
	}
	byNewName := getCompanies(t, "杭州测试公司（更新）", false)
	if len(byNewName) != 1 || byNewName[0]["target_addr"] != "嘉兴" {
		t.Fatalf("search after rename = %s", mustJSON(t, byNewName))
	}
}

func TestUpdateCompanyRejectsMissingAndDeletedRows(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	seedCompany(t, companySeed{
		Name:                    "南阳已删除公司",
		Alias:                   "旧南阳",
		UnifiedSocialCreditCode: "91330000DEL",
		DeletedAt:               "2020-01-02 03:04:05",
	})
	deleted := getCompanies(t, "南阳已删除公司", true)
	if len(deleted) != 1 {
		t.Fatalf("deleted search = %s", mustJSON(t, deleted))
	}
	id := companyIDString(t, deleted[0]["id"])

	missing := doCompanyJSON(t, http.MethodPut, "/api/admin/companies/99999", map[string]any{
		"alias": "不会写入",
	}, token)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d, body = %s", missing.Code, missing.Body.String())
	}

	soft := doCompanyJSON(t, http.MethodPut, "/api/admin/companies/"+id, map[string]any{
		"name":        "改名",
		"alias":       "新别名",
		"target_addr": "海宁",
	}, token)
	if soft.Code != http.StatusNotFound {
		t.Fatalf("deleted status = %d, body = %s", soft.Code, soft.Body.String())
	}
	again := getCompanies(t, "南阳已删除公司", true)
	if len(again) != 1 || again[0]["name"] != "南阳已删除公司" || again[0]["alias"] != "旧南阳" {
		t.Fatalf("deleted row changed: %s", mustJSON(t, again))
	}

	empty := doCompanyJSON(t, http.MethodPut, "/api/admin/companies/"+id, map[string]any{}, token)
	if empty.Code != http.StatusBadRequest || !strings.Contains(empty.Body.String(), "at least one updatable field is required") {
		t.Fatalf("empty update = %d %s", empty.Code, empty.Body.String())
	}
	badID := doCompanyJSON(t, http.MethodPut, "/api/admin/companies/abc", map[string]any{"name": "x"}, token)
	if badID.Code != http.StatusBadRequest {
		t.Fatalf("bad id = %d %s", badID.Code, badID.Body.String())
	}
}

func TestListCompaniesReturnsContentRangeAndSupportsFilterQ(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	seedCompany(t, companySeed{
		Name: "杭州测试公司A", Alias: "杭测A", UnifiedSocialCreditCode: "CODE-A",
	})
	seedCompany(t, companySeed{
		Name: "南阳食品有限公司", Alias: "南阳", UnifiedSocialCreditCode: "CODE-B",
	})
	seedCompany(t, companySeed{
		Name: "南阳已删除公司", Alias: "已删", UnifiedSocialCreditCode: "CODE-C", DeletedAt: "2020-01-02 03:04:05",
	})

	router := setupProtectedAPIRouter()
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/api/admin/companies?range=%5B0%2C10%5D&filter=%7B%7D", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status = %d %s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Range"); got != "2" {
		t.Fatalf("Content-Range = %q, want 2", got)
	}
	var rows []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &rows); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("list rows = %s", mustJSON(t, rows))
	}

	filtered := httptest.NewRequest(http.MethodGet, "/api/admin/companies?range=%5B0%2C10%5D&filter="+url.QueryEscape(`{"q":"南阳"}`), nil)
	filtered.Header.Set("Authorization", "Bearer "+token)
	fw := httptest.NewRecorder()
	router.ServeHTTP(fw, filtered)
	if fw.Code != http.StatusOK {
		t.Fatalf("filtered status = %d %s", fw.Code, fw.Body.String())
	}
	if got := fw.Header().Get("Content-Range"); got != "1" {
		t.Fatalf("filtered Content-Range = %q, want 1", got)
	}
	var matched []map[string]any
	if err := json.Unmarshal(fw.Body.Bytes(), &matched); err != nil {
		t.Fatalf("decode filtered: %v", err)
	}
	if len(matched) != 1 || matched[0]["name"] != "南阳食品有限公司" {
		t.Fatalf("filtered rows = %s", mustJSON(t, matched))
	}
}

func TestGetCompanyByID(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	seedCompany(t, companySeed{
		Name: "杭州测试公司", Alias: "杭测", UnifiedSocialCreditCode: "CODE-GET", TargetAddr: "杭州",
	})
	found := getCompanies(t, "杭州测试公司", false)
	id := companyIDString(t, found[0]["id"])

	ok := doCompanyJSON(t, http.MethodGet, "/api/admin/companies/"+id, nil, token)
	if ok.Code != http.StatusOK {
		t.Fatalf("get = %d %s", ok.Code, ok.Body.String())
	}
	body := decodeCompanyObject(t, ok.Body.Bytes())
	if body["name"] != "杭州测试公司" || body["target_addr"] != "杭州" {
		t.Fatalf("get body = %s", ok.Body.String())
	}

	missing := doCompanyJSON(t, http.MethodGet, "/api/admin/companies/99999", nil, token)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing get = %d %s", missing.Code, missing.Body.String())
	}
}

func TestDeleteCompanySoftDeletes(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	seedCompany(t, companySeed{
		Name: "杭州待删公司", Alias: "待删", UnifiedSocialCreditCode: "CODE-DEL",
	})
	found := getCompanies(t, "杭州待删公司", false)
	id := companyIDString(t, found[0]["id"])

	deleted := doCompanyJSON(t, http.MethodDelete, "/api/admin/companies/"+id, nil, token)
	if deleted.Code != http.StatusOK {
		t.Fatalf("delete = %d %s", deleted.Code, deleted.Body.String())
	}
	if got := getCompanies(t, "杭州待删公司", false); len(got) != 0 {
		t.Fatalf("active after delete = %s", mustJSON(t, got))
	}
	if got := getCompanies(t, "杭州待删公司", true); len(got) != 1 {
		t.Fatalf("includeDeleted after delete = %s", mustJSON(t, got))
	}
	again := doCompanyJSON(t, http.MethodDelete, "/api/admin/companies/"+id, nil, token)
	if again.Code != http.StatusNotFound {
		t.Fatalf("second delete = %d %s", again.Code, again.Body.String())
	}
	getGone := doCompanyJSON(t, http.MethodGet, "/api/admin/companies/"+id, nil, token)
	if getGone.Code != http.StatusNotFound {
		t.Fatalf("get after delete = %d %s", getGone.Code, getGone.Body.String())
	}
}

func TestCreateCompanyAcceptsAddressBankPhone(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	token := signedToken(t, "alice", time.Now().Add(time.Hour))
	w := doCompanyJSON(t, http.MethodPost, "/api/admin/companies", map[string]any{
		"name":                         "杭州完整公司",
		"unified_social_credit_code":   "CODE-FULL",
		"alias":                        "完整",
		"target_addr":                  "海宁",
		"addr_country":                 "中国",
		"addr_province":                "浙江",
		"addr_city":                    "杭州",
		"addr_street":                  "示例路 1 号",
		"bank_code":                    "BANK001",
		"phone_num":                    "0571-0000000",
	}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("create full = %d %s", w.Code, w.Body.String())
	}
	body := decodeCompanyObject(t, w.Body.Bytes())
	if body["addr_city"] != "杭州" || body["bank_code"] != "BANK001" || body["phone_num"] != "0571-0000000" {
		t.Fatalf("create full body = %s", w.Body.String())
	}
}

func doCompanyJSON(t *testing.T, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	router := setupProtectedAPIRouter()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeCompanyObject(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("response is not a JSON object: %v; body = %s", err, raw)
	}
	return body
}

func companyIDString(t *testing.T, value any) string {
	t.Helper()
	switch id := value.(type) {
	case float64:
		if id <= 0 || id != float64(uint64(id)) {
			t.Fatalf("id = %v", value)
		}
		return strconv.FormatUint(uint64(id), 10)
	case string:
		if id == "" {
			t.Fatalf("id is empty")
		}
		return id
	default:
		t.Fatalf("id type = %T (%v)", value, value)
		return ""
	}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(raw)
}
