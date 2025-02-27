package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"super-supply-chain/models"
)

func ReadTemplate() {

}

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

type Product struct {
	Name  string
	Price float64
}

type Data struct {
	User     string
	Age      int
	Products []Product // 结构体列表
}

// 替换简单占位符（如 {{ .User }}）
func replaceSimplePlaceholder(text string, data Data) string {
	re := regexp.MustCompile(`\{\{\s*\.(\w+)\s*\}\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		field := re.ReplaceAllString(match, "$1")
		switch field {
		case "User":
			return data.User
		case "Age":
			return fmt.Sprintf("%d", data.Age)
		default:
			return match
		}
	})
}

// 判断是否为循环区块（如 {{ range .Products }}）
func isRangeBlock(text string) bool {
	return regexp.MustCompile(`\{\{\s*range\s*\.(\w+)\s*\}\}`).MatchString(text)
}

// 处理结构体列表的循环
func processStructRange(f *excelize.File, sheet, cell string, list interface{}, startRow, startCol int) {
	val := reflect.ValueOf(list)
	if val.Kind() != reflect.Slice {
		panic("list must be a slice")
	}

	// 1. 定位模板行（占位符的下一行）
	templateRow := startRow + 1

	// 2. 插入足够多的空行（使用 InsertRows）
	length := val.Len()
	if length > 1 {
		// 插入 (length-1) 行（因为模板行已存在）
		if err := f.InsertRows(sheet, templateRow+1, length-1); err != nil { // Excel行号从1开始
			panic(err)
		}
	}

	// 3. 复制模板行并填充数据
	for i := 0; i < length; i++ {
		elem := val.Index(i)
		currentRow := templateRow + i // 当前操作的行号

		// 复制模板行的样式和值
		cols, _ := f.GetCols(sheet)
		for colIdx := 0; colIdx < len(cols); colIdx++ {
			srcCell, _ := excelize.CoordinatesToCellName(colIdx+1, templateRow)
			destCell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)

			// 复制值
			value, _ := f.GetCellValue(sheet, srcCell)
			replaced := replaceStructFields(value, elem)
			f.SetCellValue(sheet, destCell, replaced)

			// 复制样式（可选）
			styleID, _ := f.GetCellStyle(sheet, srcCell)
			f.SetCellStyle(sheet, destCell, destCell, styleID)
		}
	}

	// 4. 删除原始模板行（若需要）
	if length > 0 {
		f.RemoveRow(sheet, templateRow)
	}
}

// 替换结构体字段（如 {{ .Name }}）
func replaceStructFields(text string, elem reflect.Value) string {
	re := regexp.MustCompile(`\{\{\s*\.(\w+)\s*\}\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		fieldName := re.ReplaceAllString(match, "$1")
		// 处理指针类型
		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return match
		}
		field := elem.FieldByName(fieldName)
		if !field.IsValid() {
			return match
		}
		return fmt.Sprintf("%v", field.Interface())
	})
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

func CreateFile(data *models.ExcelData, fileName string) (string, error) {
	filePath := "/Users/menghongfei/Downloads/test_template.xlsx"
	uploadDir := GetUploadTmpDir()
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
		comment := fmt.Sprintf(`品名:%s    重量: %s 
车号:%s  车船吨位:33吨  车种: 货车 汽车
订单号:%s  分订单号:%s `, data.BaseData["name"], itemData["count"], carNum, data.BaseData["sap_number"], itemData["plan_number"])
		cellNum := j + 4
		cell := fmt.Sprintf("A%d", cellNum)

		err = f.SetSheetRow("1-发票基本信息", cell, &rowData)
		if err != nil {
			log.Fatal(err)
		}

		companyInfo := GetCompanyInfo(itemData["company_name"])
		companyName := If(companyInfo.Name != "", companyInfo.Name, itemData["company_name"])
		unifiedSocialCreditCode := If(companyInfo.Name != "", companyInfo.UnifiedSocialCreditCode, "")
		companyAddr := If(companyInfo.Name != "", companyInfo.TargetAddr, "")

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
