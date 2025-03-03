package excel_template_engines

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

type TestListItem struct {
	No   int    `json:"no"`
	Name string `json:"name"`
}

type TestData struct {
	Name string         `json:"name"`
	Age  int            `json:"age"`
	List []TestListItem `json:"list"`
}

type IterationInfo struct {
	Start    int    `json:"start"`
	End      int    `json:"end"`
	RangeKey string `json:"rangeKey"`
	IndexKey string `json:"indexKey"`
	ItemKey  string `json:"itemKey"`
}

func GetCompanyInfo(c string) models.BaseCompaniesInfos {
	var res models.BaseCompaniesInfos
	models.DB.Model(&models.BaseCompaniesInfos{}).Where("alias = ?", c).First(&res)

	return res
}

func GetArrivalPort(c string) string {
	var res []models.BaseDict
	query := models.DB.Model(&models.BaseDict{}).Where("type = ?", "港口字典").Find(&res)
	if query.Error != nil {
		log.Fatal(query.Error)
	}

	dictMap := make(map[string]string)
	for _, dict := range res {
		dictMap[dict.Key] = dict.Value
	}

	return dictMap[c]

}

func GetInvoiceTmpPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(pwd, "uploads/invoice_tmp.xlsx")
}

func CreateFeiChangFeiyong(data *models.ExcelData, fileName string) (string, error) {
	filePath := GetInvoiceTmpPath()
	uploadDir := utils.GetUploadTmpDir()
	outPath := filepath.Join(uploadDir, fmt.Sprintf("分厂开票模板_%s", fileName))
	// 1. 加载模板
	var err error

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatal(err)
		return outPath, err
	}
	defer f.Close()

	for j, itemData := range data.List {

		rowData := []interface{}{j, "增值税专用发票", "货物运输服务", "是"}
		carNum := strings.TrimSpace(itemData["car_num"])
		comment := fmt.Sprintf(`品名:%s    重量: %s吨
车号:%s  车船吨位:33吨  车种: 货车 汽车
订单号:%s  分订单号:%s `, data.BaseData["name"], itemData["count"], carNum, data.BaseData["sap_number"], itemData["plan_number"])
		cellNum := j + 4
		cell := fmt.Sprintf("A%d", cellNum)

		err = f.SetSheetRow("1-发票基本信息", cell, &rowData)
		if err != nil {
			log.Fatal(err)
		}

		companyInfo := GetCompanyInfo(itemData["company_name"])
		companyName := utils.If(companyInfo.Name != "", companyInfo.Name, itemData["company_name"])
		unifiedSocialCreditCode := utils.If(companyInfo.Name != "", companyInfo.UnifiedSocialCreditCode, "")
		companyAddr := utils.If(companyInfo.Name != "", companyInfo.TargetAddr, "")

		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("F%d", cellNum), &[]interface{}{
			companyName,
			unifiedSocialCreditCode,
		})
		if err != nil {
			log.Fatal(err)
		}

		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("R%d", cellNum), &[]interface{}{comment})
		if err != nil {
			log.Fatal(err)
		}

		err = f.SetSheetRow("1-发票基本信息", fmt.Sprintf("Y%d", cellNum), &[]interface{}{
			"展示开户银行、银行账号",
		})
		if err != nil {
			log.Fatal(err)
		}

		unit_price, err := strconv.ParseFloat(strings.TrimSpace(itemData["unit_price"]), 64)
		if err != nil {
			log.Fatal(err)
		}
		trucking_unit_price, err := strconv.ParseFloat(strings.TrimSpace(itemData["trucking_unit_price"]), 64)
		if err != nil {
			log.Fatal(err)
		}
		count, err := strconv.ParseFloat(strings.TrimSpace(itemData["count"]), 64)
		if err != nil {
			log.Fatal(err)
		}

		totalPrice := decimal.NewFromFloat(count).Mul(decimal.NewFromFloat(unit_price + trucking_unit_price)).Round(2)

		err = f.SetSheetRow("2-发票明细信息", fmt.Sprintf("A%d", cellNum), &[]interface{}{
			j, "公路运费", "3010102020100000000", "", "吨", itemData["count"], unit_price + trucking_unit_price, totalPrice, "0.09",
		})
		if err != nil {
			log.Fatal(err)
		}

		err = f.SetSheetRow("3-特定业务信息", fmt.Sprintf("A%d", cellNum), &[]interface{}{
			j,
		})
		if err != nil {
			log.Fatal(err)
		}

		arrivalPort := GetArrivalPort(data.BaseData["arrival_port"])

		err = f.SetSheetRow("3-特定业务信息", fmt.Sprintf("H%d", cellNum), &[]interface{}{
			arrivalPort, companyAddr, "公路运输", carNum, data.BaseData["name"],
		})
		if err != nil {
			log.Fatal(err)
		}

	}

	// 5. 保存结果
	if err = f.SaveAs(outPath); err != nil {
		log.Fatal(err)
	}
	return outPath, err
}
