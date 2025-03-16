package excel_template_engines

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"log"
	"super-supply-chain/models"
)

func GetBaoguandanInfo(contract_num string) (models.DynamicExcelTable, error) {
	dynamicBaoguandanInfo := models.DynamicExcelTable{}

	q := models.DB.Table("dynamic_customs_declaration_form").Where("datas->>'$.baseData.contract_num' = ?", contract_num).First(&dynamicBaoguandanInfo)
	if q.Error != nil {
		panic(q.Error)
		return dynamicBaoguandanInfo, q.Error
	}
	return dynamicBaoguandanInfo, nil
}

func GetChengxinInvoiceFile(datas []models.DynamicExcelTable, newFilePath string) error {
	var err error
	invoice_tmp_path := GetInvoiceTmpPath()

	f, err := excelize.OpenFile(invoice_tmp_path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	for index, data := range datas {
		baseData := data.Datas.BaseData
		contract_num := baseData["contract_num"]
		product_name := baseData["product_name"]
		container_type := baseData["container_type"]
		port_area := baseData["port_area"]
		total_price := baseData["total_price"]
		//port_area := baseData["port_area"]
		rowIndex := index + 4
		baoguandanInfo, err := GetBaoguandanInfo(contract_num)
		companyName := baoguandanInfo.Datas.BaseData["domestic_consignee"]
		companyInfo := GetCompanyInfo(companyName)
		weight, err := decimal.NewFromString(baoguandanInfo.Datas.BaseData["weight"])

		comment := fmt.Sprintf(`%s   %s吨  %s   %s   %s`,
			product_name,
			weight.Div(decimal.NewFromInt(1000)).String(),
			port_area,
			contract_num,
			container_type,
		)

		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("A%d", rowIndex), &[]interface{}{
			contract_num,
			"增值税专用发票",
			"",
			"是",
			"",
			companyInfo.Name,
			companyInfo.UnifiedSocialCreditCode,
		})
		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("R%d", rowIndex), &[]interface{}{comment})

		detail := []string{
			contract_num,
			"代理清关费",
			"3040802020000000000",
			"",
			"",
			"",
			"",
			total_price,
			"0.06",
		}

		err = f.SetSheetRow("2-发票明细信息", fmt.Sprintf("A%d", rowIndex), &detail)

		if err != nil {
			panic(err)
		}
	}

	if err = f.SaveAs(newFilePath); err != nil {
		panic(err)
		return err
	}

	return nil
}
