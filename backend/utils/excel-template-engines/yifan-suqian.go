package excel_template_engines

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func GetCompanyInfoByName(name string) models.BaseCompaniesInfos {
	companyInfo := models.BaseCompaniesInfos{}
	q := models.DB.Model(models.BaseCompaniesInfos{}).Where("alias = ?", name).First(&companyInfo)
	if q.Error != nil {
		utils.Logger.Error(q.Error.Error())
	}
	return companyInfo
}

// 运费发票
func GetFreightInvoiceFile(datas []models.DynamicExcelTable, newFilePath string) error {
	invoice_tmp_path := GetInvoiceTmpPath()
	var err error
	f, err := excelize.OpenFile(invoice_tmp_path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()
	var invoce_details [][]string
	var extra_details [][]string

	for index, data := range datas {
		excelData := data.Datas
		companyInfo := GetCompanyInfoByName("浙江迅尔智链货运有限公司")
		rowIndex := index + 4
		arrivalPort := data.Datas.BaseData["arrival_port"]
		comment := fmt.Sprintf(`%s
进口原料费`, excelData.BaseData["sap_number"])
		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("A%d", rowIndex), &[]interface{}{
			strconv.Itoa(int(data.ID)),
			"增值税专用发票",
			"货物运输服务",
			"是",
			companyInfo.Name,
			companyInfo.UnifiedSocialCreditCode,
		})
		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("R%d", rowIndex), &[]interface{}{comment})
		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("Y%d", rowIndex), &[]interface{}{
			"展示开户银行、银行账号",
		})
		productName := excelData.BaseData["product_name"]

		for _, item := range data.Datas.List {
			var err error
			count, err := decimal.NewFromString(item["count"])
			unit_price, err := decimal.NewFromString(item["unit_price"])
			itemCompanyInifo := GetCompanyInfoByName(item["company_name"])
			car_num := item["car_num"]

			invoce_details = append(invoce_details, []string{
				strconv.Itoa(int(data.ID)),
				"公路运输",
				"3010102020100000000",
				"",
				"吨",
				count.String(),
				unit_price.String(),
				count.Mul(unit_price).Round(2).String(),
				"0.09",
			})
			extra_details = append(extra_details, []string{
				strconv.Itoa(int(data.ID)),
				"",
				"",
				"",
				"",
				"",
				"",
				PortInfoMap[arrivalPort].Addr,
				itemCompanyInifo.TargetAddr,
				"公路运输",
				car_num,
				productName,
			})
			if err != nil {
				return err
			}
		}

	}

	for index, detail := range invoce_details {
		rowIndex := index + 4
		err = f.SetSheetRow("2-发票明细信息", fmt.Sprintf("A%d", rowIndex), &detail)
	}
	for i, detail := range extra_details {
		rowIndex := i + 4
		err = f.SetSheetRow("3-特定业务信息", fmt.Sprintf("A%d", rowIndex), &detail)
	}

	if err = f.SaveAs(newFilePath); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetChangejiuTmpPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(pwd, "uploads/changjiu_tmp.xlsx")
}

// 清关费发票
func GetClearanceInvoiceFile(datas []models.DynamicExcelTable, newFilePath string) error {
	invoice_tmp_path := GetInvoiceTmpPath()

	var err error
	f, err := excelize.OpenFile(invoice_tmp_path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()
	var invoce_details [][]string

	for index, data := range datas {
		excelData := data.Datas
		companyInfo := GetCompanyInfoByName(excelData.BaseData["invoice_company"])
		rowIndex := index + 4
		arrivalPort := data.Datas.BaseData["arrival_port"]
		comment := fmt.Sprintf(`业务编号: %s
货物名称: %s
数量: %s
到港口岸: %s`, excelData.BaseData["sap_number"],
			excelData.BaseData["product_name"],
			excelData.BaseData["total_count"],
			arrivalPort,
		)
		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("A%d", rowIndex), &[]interface{}{
			strconv.Itoa(int(data.ID)),
			"增值税专用发票",
			"货物运输服务",
			"是",
			companyInfo.Name,
			companyInfo.UnifiedSocialCreditCode,
		})
		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("R%d", rowIndex), &[]interface{}{comment})
		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("Y%d", rowIndex), &[]interface{}{
			"展示开户银行、银行账号",
		})
		var err error

		unpacking_fee_unit_count, err := decimal.NewFromString(excelData.BaseData["unpacking_fee_unit_count"])
		unpacking_fee_unit_price, err := decimal.NewFromString(excelData.BaseData["unpacking_fee_unit_price"])

		invoce_details = append(invoce_details, []string{
			strconv.Itoa(int(data.ID)),
			"代理清关",
			"3010102020100000000",
			"40'",
			"柜",
			unpacking_fee_unit_count.String(),
			unpacking_fee_unit_price.String(),
			unpacking_fee_unit_count.Mul(unpacking_fee_unit_price).Round(2).String(),
			"0.09",
		})
		if excelData.BaseData["clearance_fee_count"] != "" {
			var err error
			clearance_fee_count, err := decimal.NewFromString(excelData.BaseData["clearance_fee_count"])
			clearance_fee_unit_price, err := decimal.NewFromString(excelData.BaseData["clearance_fee_unit_price"])
			invoce_details = append(invoce_details, []string{
				strconv.Itoa(int(data.ID)),
				"代理掏箱",
				"3010102020100000000",
				"40'",
				"柜",
				clearance_fee_count.String(),
				clearance_fee_unit_price.String(),
				clearance_fee_count.Mul(clearance_fee_unit_price).Round(2).String(),
				"0.09",
			})
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

	}

	for index, detail := range invoce_details {
		rowIndex := index + 4
		err = f.SetSheetRow("2-发票明细信息", fmt.Sprintf("A%d", rowIndex), &detail)
	}

	if err = f.SaveAs(newFilePath); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetShortHaulFile(datas []models.DynamicExcelTable, newFilePath string) error {
	invoice_tmp_path := GetChangejiuTmpPath()

	f, err := excelize.OpenFile(invoice_tmp_path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	for i, data := range datas {
		var err error
		err = f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), &[]interface{}{
			data.Datas.BaseData["invoice_company"],
			data.Datas.BaseData["product_name"],
			"短驳费",
			data.Datas.BaseData["total_count"],
			data.Datas.BaseData["short_haul_fee_price"],
			data.Datas.BaseData["arrival_port"],
			data.Datas.BaseData["sap_number"],
		})
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	if err = f.SaveAs(newFilePath); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func GetExcelExportFilePath(tableName string, ids []string, queryType string) (string, error) {

	var datas []models.DynamicExcelTable
	uploadDir := utils.GetUploadTmpDir()
	newFilePath := filepath.Join(uploadDir, uuid.New().String()+".xlsx")
	var err error

	query := models.DB.Table(tableName).Where("id in (?)", ids).Find(&datas)
	if query.Error != nil {
		return "", query.Error
	}
	switch queryType {
	case "invoice_freight":
		err = GetFreightInvoiceFile(datas, newFilePath)
	case "invoice_clearance":
		err = GetClearanceInvoiceFile(datas, newFilePath)
	case "shortHaul":
		err = GetShortHaulFile(datas, newFilePath)
	}
	if err != nil {
		return "", err
	}

	return newFilePath, nil
}
