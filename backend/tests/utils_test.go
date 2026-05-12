package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"super-supply-chain/utils"
	excel_template_engines "super-supply-chain/utils/excel-template-engines"
)

// 测试列表查询参数能解析 React Admin 默认传入的 JSON 格式 range/filter/sort。
func TestGetListQueryParamsParsesReactAdminQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/items?range=%5B3,8%5D&sort=%5B%22id%22,%22DESC%22%5D&filter=%7B%22start%22:%222026-01-01%22,%22end%22:%222026-01-31%22%7D", nil)
	c.Request = req

	got, err := utils.GetListQueryParams(c)
	if err != nil {
		t.Fatalf("GetListQueryParams returned error: %v", err)
	}
	if got.Offset != 3 || got.Limit != 5 {
		t.Fatalf("GetListQueryParams offset/limit = %d/%d, want 3/5", got.Offset, got.Limit)
	}
	if got.Sort != `["id","DESC"]` {
		t.Fatalf("Sort = %q, want raw sort query", got.Sort)
	}
	if got.Filter.Start != "2026-01-01" || got.Filter.End != "2026-01-31" {
		t.Fatalf("Filter = %#v, want parsed start/end", got.Filter)
	}
}

// 测试列表查询参数仍兼容重复 range 参数和 filter.start/filter.end 形式。
func TestGetListQueryParamsKeepsRepeatedRangeCompatibility(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/items?range=3&range=8&filter.start=2026-01-01&filter.end=2026-01-31", nil)
	c.Request = req

	got, err := utils.GetListQueryParams(c)
	if err != nil {
		t.Fatalf("GetListQueryParams returned error: %v", err)
	}
	if got.Offset != 3 || got.Limit != 5 {
		t.Fatalf("GetListQueryParams offset/limit = %d/%d, want 3/5", got.Offset, got.Limit)
	}
	if got.Filter.Start != "2026-01-01" || got.Filter.End != "2026-01-31" {
		t.Fatalf("Filter = %#v, want parsed start/end", got.Filter)
	}
}

// 测试列表接口响应会写入 React Admin 读取总数所需的 Content-Range header。
func TestSetContentRangeSetsTotalHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.SetContentRange(c, 42)

	if got := w.Header().Get("Content-Range"); got != "42" {
		t.Fatalf("Content-Range = %q, want %q", got, "42")
	}
}

// 测试报关单文本中能提取集装箱标箱数。
func TestGetPkgCountExtractsContainerCount(t *testing.T) {
	got, err := utils.GetPkgCount("集装箱标箱数及号码：12\nABCD1234567")
	if err != nil {
		t.Fatalf("GetPkgCount returned error: %v", err)
	}
	if got != "12" {
		t.Fatalf("GetPkgCount() = %q, want %q", got, "12")
	}
}

// 测试报关单三行商品单元格能提取最后一个商品名称字段。
func TestGetProductNameExtractsLastFieldFromThreeLineCell(t *testing.T) {
	got, err := utils.GetProductName("1 22029900 全脂奶粉\n规格型号\n其他")
	if err != nil {
		t.Fatalf("GetProductName returned error: %v", err)
	}
	if got != "全脂奶粉" {
		t.Fatalf("GetProductName() = %q, want %q", got, "全脂奶粉")
	}
}

// 测试报关单价格单元格能同时提取价格和币种代码。
func TestGetPriceExtractsCurrencyCode(t *testing.T) {
	price, unit, err := utils.GetPrice("成交方式\n123.45\n美元")
	if err != nil {
		t.Fatalf("GetPrice returned error: %v", err)
	}
	if price != "123.45" || unit != "USD" {
		t.Fatalf("GetPrice() = %q, %q; want %q, %q", price, unit, "123.45", "USD")
	}
}

// 测试字符串切片匹配使用包含关系，而不是完全相等。
func TestSliceContainsStringUsesSubstringMatch(t *testing.T) {
	if !utils.SliceContainsString([]string{"商品编号", "商品名称及规格型号"}, "规格型号") {
		t.Fatal("SliceContainsString should match substrings")
	}
}

// 测试 Excel 行扫描能记录每个“合计”行的行号。
func TestGetToTalRowIndexsRecordsEachTotalRow(t *testing.T) {
	base := map[string]string{}
	rows := [][]string{
		{"货物", "数量"},
		{"合计", "10"},
		{"备注", ""},
		{"", "合计"},
	}

	utils.GetToTalRowIndexs(rows, base)

	if base["total_1"] != "1" || base["total_2"] != "3" {
		t.Fatalf("total rows = %#v, want total_1=1 and total_2=3", base)
	}
}

// 测试通用泛型 helper 的基础行为。
func TestGenericHelpers(t *testing.T) {
	mapped := utils.Map([]int{1, 2, 3}, func(v int) string {
		return string(rune('a' + v - 1))
	})
	if mapped[0] != "a" || mapped[1] != "b" || mapped[2] != "c" {
		t.Fatalf("Map() = %#v, want [a b c]", mapped)
	}

	if safe := utils.GetSafeArray[string](nil); safe == nil || len(safe) != 0 {
		t.Fatalf("GetSafeArray(nil) = %#v, want empty non-nil slice", safe)
	}

	if got := utils.If(true, "yes", "no"); got != "yes" {
		t.Fatalf("If(true) = %q, want yes", got)
	}
}

// 测试一帆费用模板中商品名称的业务归一化规则。
func TestYifanFeiyongProductNameNormalization(t *testing.T) {
	if got := excel_template_engines.GetProductName("进口全脂奶粉(>=26%)"); got != "全脂奶粉" {
		t.Fatalf("GetProductName() = %q, want 全脂奶粉", got)
	}
	if got := excel_template_engines.GetProductName("进口全脂乳粉(>=26%)"); got != "全脂奶粉" {
		t.Fatalf("GetProductName() = %q, want 全脂奶粉", got)
	}
	if got := excel_template_engines.GetProductName("塑料粒子XX"); got != "" {
		t.Fatalf("GetProductName() = %q, want empty string", got)
	}
}
