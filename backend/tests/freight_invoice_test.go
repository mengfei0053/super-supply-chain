package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
	"super-supply-chain/models"
	excel_template_engines "super-supply-chain/utils/excel-template-engines"
)

// 测试运费发票模板能生成基础信息、明细信息和特定业务信息。
func TestGetFreightInvoiceFileGeneratesWorkbook(t *testing.T) {
	setupTestDB(t, &models.BaseCompaniesInfos{})
	seedFreightInvoiceCompanies(t)
	chdirForInvoiceTemplate(t)

	outputPath := filepath.Join(t.TempDir(), "freight-invoice.xlsx")
	datas := []models.DynamicExcelTable{
		{
			Datas: models.ExcelData{
				BaseData: map[string]string{
					"sap_number":   "SAP-001",
					"arrival_port": "上海洋山",
					"product_name": "进口奶粉",
				},
				List: []map[string]string{
					{
						"count":        "2",
						"unit_price":   "100",
						"price":        "200",
						"pkg_num":      "1",
						"company_name": "杭州测试公司",
						"car_num":      "浙A12345",
					},
					{
						"count":               "3",
						"unit_price":          "100",
						"trucking_unit_price": "2",
						"price":               "300",
						"pkg_num":             "",
						"company_name":        "杭州测试公司",
						"car_num":             "浙A67890",
					},
				},
			},
		},
	}

	if err := excel_template_engines.GetFreightInvoiceFile(datas, outputPath); err != nil {
		t.Fatalf("GetFreightInvoiceFile returned error: %v", err)
	}

	workbook, err := excelize.OpenFile(outputPath)
	if err != nil {
		t.Fatalf("open generated workbook: %v", err)
	}
	defer workbook.Close()

	assertCell(t, workbook, "1-发票基本信息", "A4", "SAP-001")
	assertCell(t, workbook, "1-发票基本信息", "B4", "增值税专用发票")
	assertCell(t, workbook, "1-发票基本信息", "F4", "浙江迅尔智链货运有限公司")
	assertCell(t, workbook, "1-发票基本信息", "G4", "91330000TEST")
	assertCell(t, workbook, "1-发票基本信息", "W4", "SAP-001\n进口原料运费")
	assertCell(t, workbook, "1-发票基本信息", "AD4", "展示开户银行、银行账号")

	assertCell(t, workbook, "2-发票明细信息", "A4", "SAP-001")
	assertCell(t, workbook, "2-发票明细信息", "B4", "公路运输")
	assertCell(t, workbook, "2-发票明细信息", "F4", "2")
	assertCell(t, workbook, "2-发票明细信息", "G4", "")
	assertCell(t, workbook, "2-发票明细信息", "H4", "200")
	assertCell(t, workbook, "2-发票明细信息", "I4", "0.09")

	assertCell(t, workbook, "2-发票明细信息", "A5", "SAP-001")
	assertCell(t, workbook, "2-发票明细信息", "F5", "3")
	assertCell(t, workbook, "2-发票明细信息", "G5", "102")
	assertCell(t, workbook, "2-发票明细信息", "H5", "306")

	assertCell(t, workbook, "3-特定业务信息", "A4", "SAP-001")
	assertCell(t, workbook, "3-特定业务信息", "H4", "上海")
	assertCell(t, workbook, "3-特定业务信息", "I4", "杭州")
	assertCell(t, workbook, "3-特定业务信息", "J4", "公路运输")
	assertCell(t, workbook, "3-特定业务信息", "K4", "浙A12345")
	assertCell(t, workbook, "3-特定业务信息", "L4", "进口奶粉")
}

func chdirForInvoiceTemplate(t *testing.T) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatalf("change working directory to backend: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
}

func seedFreightInvoiceCompanies(t *testing.T) {
	t.Helper()

	companies := []models.BaseCompaniesInfos{
		{
			Name:                    "浙江迅尔智链货运有限公司",
			TargetAddr:              "杭州",
			Alias:                   "浙江迅尔智链货运",
			UnifiedSocialCreditCode: "91330000TEST",
		},
		{
			Name:                    "杭州测试公司",
			TargetAddr:              "杭州",
			Alias:                   "杭州测试",
			UnifiedSocialCreditCode: "91330100TEST",
		},
	}

	for _, company := range companies {
		if err := models.DB.Exec(
			`INSERT INTO base_companies_infos (name, target_addr, alias, unified_social_credit_code, deleted_at) VALUES (?, ?, ?, ?, NULL)`,
			company.Name,
			company.TargetAddr,
			company.Alias,
			company.UnifiedSocialCreditCode,
		).Error; err != nil {
			t.Fatalf("seed freight invoice company %s: %v", company.Name, err)
		}
	}
}

func assertCell(t *testing.T, workbook *excelize.File, sheet, cell, want string) {
	t.Helper()

	got, err := workbook.GetCellValue(sheet, cell)
	if err != nil {
		t.Fatalf("read %s!%s: %v", sheet, cell, err)
	}
	if got != want {
		t.Fatalf("%s!%s = %q, want %q", sheet, cell, got, want)
	}
}
