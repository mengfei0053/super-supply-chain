package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
	"super-supply-chain/models"
	excelTemplateEngines "super-supply-chain/utils/excel-template-engines"
)

// 测试运费发票模板能生成基础信息、明细信息和特定业务信息。
func TestGetFreightInvoiceFileGeneratesWorkbook(t *testing.T) {
	setupInvoiceGeneratorTest(t)

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

	if err := excelTemplateEngines.GetFreightInvoiceFile(datas, outputPath); err != nil {
		t.Fatalf("GetFreightInvoiceFile returned error: %v", err)
	}

	workbook := openWorkbook(t, outputPath)
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

// 测试清关费发票模板能生成清关费和掏箱费明细。
func TestGetClearanceInvoiceFileGeneratesWorkbook(t *testing.T) {
	setupInvoiceGeneratorTest(t)

	outputPath := filepath.Join(t.TempDir(), "clearance-invoice.xlsx")
	datas := []models.DynamicExcelTable{
		{
			Datas: models.ExcelData{
				BaseData: map[string]string{
					"sap_number":               "CLR-001",
					"bus_num":                  "BUS-001",
					"invoice_company":          "杭州测试公司",
					"arrival_port":             "上海洋山",
					"product_name":             "进口果汁",
					"total_count":              "12.5",
					"clearance_fee_count":      "2",
					"clearance_fee_unit_price": "150",
					"unpacking_fee_unit_count": "1",
					"unpacking_fee_unit_price": "80",
				},
			},
		},
	}

	if err := excelTemplateEngines.GetClearanceInvoiceFile(datas, outputPath); err != nil {
		t.Fatalf("GetClearanceInvoiceFile returned error: %v", err)
	}

	workbook := openWorkbook(t, outputPath)
	defer workbook.Close()

	assertCell(t, workbook, "1-发票基本信息", "A4", "CLR-001")
	assertCell(t, workbook, "1-发票基本信息", "B4", "增值税专用发票")
	assertCell(t, workbook, "1-发票基本信息", "F4", "杭州测试公司")
	assertCell(t, workbook, "1-发票基本信息", "G4", "91330100TEST")
	assertCell(t, workbook, "1-发票基本信息", "W4", "订单号：CLR-001  业务编号: BUS-001\n货物名称: 进口果汁\n数量: 12.5 吨\n到港口岸: 上海洋山")
	assertCell(t, workbook, "1-发票基本信息", "AD4", "展示开户银行、银行账号")

	assertCell(t, workbook, "2-发票明细信息", "A4", "CLR-001")
	assertCell(t, workbook, "2-发票明细信息", "B4", "代理清关费")
	assertCell(t, workbook, "2-发票明细信息", "C4", "3040802020000000000")
	assertCell(t, workbook, "2-发票明细信息", "D4", "40'")
	assertCell(t, workbook, "2-发票明细信息", "E4", "柜")
	assertCell(t, workbook, "2-发票明细信息", "F4", "2")
	assertCell(t, workbook, "2-发票明细信息", "G4", "150")
	assertCell(t, workbook, "2-发票明细信息", "H4", "300")
	assertCell(t, workbook, "2-发票明细信息", "I4", "0.06")

	assertCell(t, workbook, "2-发票明细信息", "A5", "CLR-001")
	assertCell(t, workbook, "2-发票明细信息", "B5", "代理掏箱费")
	assertCell(t, workbook, "2-发票明细信息", "H5", "80")
	assertCell(t, workbook, "2-发票明细信息", "I5", "0.06")
}

// 测试掏箱费数量或单价缺失、为 0 时不生成掏箱费明细。
func TestGetUnpackingInvoiceFileSkipsMissingOrZeroValues(t *testing.T) {
	setupInvoiceGeneratorTest(t)

	outputPath := filepath.Join(t.TempDir(), "unpacking-invoice.xlsx")
	datas := []models.DynamicExcelTable{
		{
			Datas: models.ExcelData{BaseData: map[string]string{
				"sap_number":               "UNP-001",
				"bus_num":                  "BUS-001",
				"invoice_company":          "杭州测试公司",
				"arrival_port":             "上海洋山",
				"product_name":             "进口果汁",
				"total_count":              "12.5",
				"unpacking_fee_unit_count": "0",
				"unpacking_fee_unit_price": "80",
			}},
		},
		{
			Datas: models.ExcelData{BaseData: map[string]string{
				"sap_number":               "UNP-002",
				"bus_num":                  "BUS-002",
				"invoice_company":          "杭州测试公司",
				"arrival_port":             "上海洋山",
				"product_name":             "进口果汁",
				"total_count":              "12.5",
				"unpacking_fee_unit_count": "1",
				"unpacking_fee_unit_price": "0.00",
			}},
		},
		{
			Datas: models.ExcelData{BaseData: map[string]string{
				"sap_number":               "UNP-003",
				"bus_num":                  "BUS-003",
				"invoice_company":          "杭州测试公司",
				"arrival_port":             "上海洋山",
				"product_name":             "进口果汁",
				"total_count":              "12.5",
				"unpacking_fee_unit_count": "1",
			}},
		},
	}

	if err := excelTemplateEngines.GetUnpackingInvoiceFile(datas, outputPath); err != nil {
		t.Fatalf("GetUnpackingInvoiceFile returned error: %v", err)
	}

	workbook := openWorkbook(t, outputPath)
	defer workbook.Close()

	assertCell(t, workbook, "1-发票基本信息", "A4", "")
	assertCell(t, workbook, "2-发票明细信息", "A4", "")
	assertCell(t, workbook, "2-发票明细信息", "B4", "")
}

// 测试短驳费用模板只导出短驳费不为空且不为 0 的数据。
func TestGetShortHaulFileGeneratesWorkbook(t *testing.T) {
	setupInvoiceGeneratorTest(t)

	outputPath := filepath.Join(t.TempDir(), "short-haul.xlsx")
	datas := []models.DynamicExcelTable{
		{
			Datas: models.ExcelData{BaseData: map[string]string{
				"invoice_company":      "杭州测试公司",
				"product_name":         "进口奶粉",
				"total_count":          "8",
				"short_haul_fee_price": "128.50",
				"arrival_port":         "上海洋山",
				"sap_number":           "SH-001",
			}},
		},
		{
			Datas: models.ExcelData{BaseData: map[string]string{
				"invoice_company":      "杭州测试公司",
				"product_name":         "不应导出",
				"total_count":          "1",
				"short_haul_fee_price": "0",
				"arrival_port":         "上海洋山",
				"sap_number":           "SH-002",
			}},
		},
	}

	if err := excelTemplateEngines.GetShortHaulFile(datas, outputPath); err != nil {
		t.Fatalf("GetShortHaulFile returned error: %v", err)
	}

	workbook := openWorkbook(t, outputPath)
	defer workbook.Close()

	assertCell(t, workbook, "Sheet1", "A2", "杭州测试公司")
	assertCell(t, workbook, "Sheet1", "B2", "进口奶粉")
	assertCell(t, workbook, "Sheet1", "C2", "短驳费")
	assertCell(t, workbook, "Sheet1", "D2", "8")
	assertCell(t, workbook, "Sheet1", "E2", "128.50")
	assertCell(t, workbook, "Sheet1", "F2", "上海洋山")
	assertCell(t, workbook, "Sheet1", "G2", "SH-001")
	assertCell(t, workbook, "Sheet1", "A3", "")
}

// 测试短驳发票模板能生成发票基本信息、明细和特定业务信息。
func TestGetShortHaulInvoceGeneratesWorkbook(t *testing.T) {
	setupInvoiceGeneratorTest(t)

	outputPath := filepath.Join(t.TempDir(), "short-haul-invoice.xlsx")
	datas := []models.DynamicExcelTable{
		{
			Datas: models.ExcelData{BaseData: map[string]string{
				"invoice_company":      "杭州测试公司",
				"product_name":         "进口奶粉",
				"total_count":          "8",
				"short_haul_fee_price": "128.50",
				"arrival_port":         "上海洋山",
				"sap_number":           "SHI-001",
				"short_car_num":        "浙A00001",
			}},
		},
	}

	if err := excelTemplateEngines.GetShortHaulInvoce(datas, outputPath); err != nil {
		t.Fatalf("GetShortHaulInvoce returned error: %v", err)
	}

	workbook := openWorkbook(t, outputPath)
	defer workbook.Close()

	assertCell(t, workbook, "1-发票基本信息", "A4", "SHI-001")
	assertCell(t, workbook, "1-发票基本信息", "B4", "增值税专用发票")
	assertCell(t, workbook, "1-发票基本信息", "C4", "货物运输服务")
	assertCell(t, workbook, "1-发票基本信息", "F4", "杭州测试公司")
	assertCell(t, workbook, "1-发票基本信息", "G4", "91330100TEST")
	assertCell(t, workbook, "1-发票基本信息", "W4", "上海到上海 上海洋山 短驳费\n品名：进口奶粉 重量: 8吨\n车号: 浙A00001  车船吨位: 30吨  车种: 货车 汽车\n工作编号: SHI-001 ")
	assertCell(t, workbook, "1-发票基本信息", "AD4", "展示开户银行、银行账号")

	assertCell(t, workbook, "2-发票明细信息", "A4", "SHI-001")
	assertCell(t, workbook, "2-发票明细信息", "B4", "公路运费")
	assertCell(t, workbook, "2-发票明细信息", "F4", "8")
	assertCell(t, workbook, "2-发票明细信息", "G4", "")
	assertCell(t, workbook, "2-发票明细信息", "H4", "128.50")
	assertCell(t, workbook, "2-发票明细信息", "I4", "0.09")

	assertCell(t, workbook, "3-特定业务信息", "A4", "SHI-001")
	assertCell(t, workbook, "3-特定业务信息", "H4", "上海")
	assertCell(t, workbook, "3-特定业务信息", "I4", "上海")
	assertCell(t, workbook, "3-特定业务信息", "K4", "浙A00001")
	assertCell(t, workbook, "3-特定业务信息", "L4", "进口奶粉")
}

// 测试诚信包装发票模板会读取对应报关单信息并生成发票明细。
func TestGetChengxinInvoiceFileGeneratesWorkbook(t *testing.T) {
	setupInvoiceGeneratorTest(t)
	seedChengxinBaoguandan(t)

	outputPath := filepath.Join(t.TempDir(), "chengxin-invoice.xlsx")
	datas := []models.DynamicExcelTable{
		{
			Datas: models.ExcelData{BaseData: map[string]string{
				"contract_num":   "CON-001",
				"product_name":   "进口奶粉",
				"container_type": "40HQ",
				"port_area":      "上海洋山",
				"total_price":    "600",
			}},
		},
	}

	if err := excelTemplateEngines.GetChengxinInvoiceFile(datas, outputPath); err != nil {
		t.Fatalf("GetChengxinInvoiceFile returned error: %v", err)
	}

	workbook := openWorkbook(t, outputPath)
	defer workbook.Close()

	assertCell(t, workbook, "1-发票基本信息", "A4", "CON-001")
	assertCell(t, workbook, "1-发票基本信息", "B4", "增值税专用发票")
	assertCell(t, workbook, "1-发票基本信息", "F4", "杭州测试公司")
	assertCell(t, workbook, "1-发票基本信息", "G4", "91330100TEST")
	assertCell(t, workbook, "1-发票基本信息", "R4", "进口奶粉   2.5吨  上海洋山   CON-001   40HQ")

	assertCell(t, workbook, "2-发票明细信息", "A4", "CON-001")
	assertCell(t, workbook, "2-发票明细信息", "B4", "代理清关费")
	assertCell(t, workbook, "2-发票明细信息", "C4", "3040802020000000000")
	assertCell(t, workbook, "2-发票明细信息", "H4", "600")
	assertCell(t, workbook, "2-发票明细信息", "I4", "0.06")
}

func setupInvoiceGeneratorTest(t *testing.T) {
	t.Helper()

	setupTestDB(t, &models.BaseCompaniesInfos{}, &models.DynamicExcelTable{})
	seedInvoiceCompanies(t)
	chdirForInvoiceTemplate(t)
}

func chdirForInvoiceTemplate(t *testing.T) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if filepath.Base(previousDir) != "backend" {
		if err := os.Chdir(".."); err != nil {
			t.Fatalf("change working directory to backend: %v", err)
		}
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})
}

func seedInvoiceCompanies(t *testing.T) {
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
			t.Fatalf("seed invoice company %s: %v", company.Name, err)
		}
	}
}

func seedChengxinBaoguandan(t *testing.T) {
	t.Helper()

	if err := models.DB.Table("dynamic_customs_declaration_form").AutoMigrate(&models.DynamicExcelTable{}); err != nil {
		t.Fatalf("migrate chengxin customs declaration table: %v", err)
	}

	baoguandan := models.DynamicExcelTable{
		FileName:       "baoguandan.xlsx",
		NasFileName:    "baoguandan.xlsx",
		UploadFilePath: "/tmp/baoguandan.xlsx",
		Datas: models.ExcelData{BaseData: map[string]string{
			"contract_num":       "CON-001",
			"domestic_consignee": "杭州测试公司",
			"weight":             "2500",
		}},
	}

	if err := models.DB.Table("dynamic_customs_declaration_form").Create(&baoguandan).Error; err != nil {
		t.Fatalf("seed chengxin customs declaration: %v", err)
	}
}

func openWorkbook(t *testing.T, path string) *excelize.File {
	t.Helper()

	workbook, err := excelize.OpenFile(path)
	if err != nil {
		t.Fatalf("open generated workbook: %v", err)
	}
	return workbook
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
