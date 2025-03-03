package excel_template_engines

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"log"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func GetProductName(name string) string {
	re := regexp.MustCompile(`.*全脂(奶|乳)粉.*`)

	if re.MatchString(name) {
		return re.ReplaceAllString(name, "全脂奶粉")
	}

	return ""
}

var PortMap = map[string]string{
	"上海洋山": "上海口岸",
	"天津东疆": "天津口岸",
	"天津新港": "天津口岸",
	"广州黄埔": "广州口岸",
}

func GetClearanceFeeByPort(port string, costType models.CostType, containerType int) float64 {
	portName := PortMap[port]
	var clearancePrice models.ClearancePriceBase

	fmt.Println("port", port)
	fmt.Println("portName", portName)

	models.DB.Model(&models.ClearancePriceBase{}).Where(
		"port = ? and cost_type = ? and container_type_enum = ?",
		portName, costType, containerType).
		First(&clearancePrice)

	return clearancePrice.Price

}

func GetContainerCount(num float64, base float64) string {
	const (
		epsilon = 1e-9 // 精度容差
	)

	// 处理特殊值（NaN/Inf）
	if math.IsNaN(num) || math.IsInf(num, 0) {
		return ""
	}

	// 计算商并四舍五入到最近整数
	quotient := num / base
	rounded := math.Round(quotient)

	// 双重验证：商接近整数，且反向计算值接近原数
	if math.Abs(quotient-rounded) < epsilon &&
		math.Abs(num-rounded*base) < epsilon {
		return fmt.Sprintf("%.0f", rounded) // 返回整数倍数字符串
	}
	return ""
}

func GetProductInfo(productName string) (models.ProductInfoBase, error) {
	name := GetProductName(productName)
	var productInfo models.ProductInfoBase
	if name != "" {
		q := models.DB.Model(&models.ProductInfoBase{}).Where("product_name = ?", name).First(&productInfo)
		if q.Error != nil {
			log.Fatal(q.Error)
			return productInfo, q.Error
		}
		return productInfo, nil
	}
	return productInfo, errors.New("Product not found")
}

func GetZhengXiangTotal(data *models.DynamicExcelTable, productInfo models.ProductInfoBase) (int, error) {
	var total int
	for _, item := range data.Datas.List {
		var err error
		var containerCount int64
		planned_count, err := strconv.ParseFloat(item["planned_count"], 64)
		containerCountStr := GetContainerCount(planned_count, productInfo.ContainerTypeWeight)

		if containerCountStr != "" {
			containerCount, err = strconv.ParseInt(containerCountStr, 10, 64)
		}
		if err != nil {
			log.Fatal(err)
			return total, err
		}
		total += int(containerCount)

	}
	return total, nil
}

/*
*
普通散货
冷冻散货
普通整柜
冷冻整柜
普通整柜出口拖车费
框架箱整柜
普通散货内陆运输费
*/
const (
	transportationMode                     = "普通散货"
	coldTransportationMode                 = "冷冻散货"
	wholeCabinet                           = "普通整柜"
	coldWholeCabinet                       = "冷冻整柜"
	wholeCabinetExportTrailerFee           = "普通整柜出口拖车费"
	frameBoxWholeCabinet                   = "框架箱整柜"
	ordinaryBulkCargoLandTransportationFee = "普通散货内陆运输费"
)

type FreightInfo struct {
	Freight  float64
	ExtraPay float64
}

type PortInfo struct {
	PortName string
	ExtraPay string
	Addr     string
}

var PortInfoMap = map[string]PortInfo{
	"上海洋山": {
		PortName: "上海口岸",
		Addr:     "上海",
		ExtraPay: "洋山补差",
	},
	"上海外高桥": {
		PortName: "上海口岸",
		ExtraPay: "",
		Addr:     "上海",
	},
	"天津东疆": {
		PortName: "天津口岸",
		ExtraPay: "东疆补差",
		Addr:     "天津",
	},
	"天津新港": {
		PortName: "天津口岸",
		ExtraPay: "",
		Addr:     "天津",
	},
	"广州黄埔": {
		PortName: "广州口岸",
		ExtraPay: "南沙补差",
		Addr:     "广州",
	},
}

func GetUnitFreight(port string, containerCount string, company_name string) FreightInfo {
	var freight models.FreightBase
	var extraPay models.FreightBase
	var trans string
	var companyInfo models.BaseCompaniesInfos

	portInfo := PortInfoMap[port]
	if containerCount != "" {
		trans = wholeCabinet
	} else {
		trans = transportationMode
	}

	const (
		YEAR = "2024"
	)

	models.DB.Model(&models.BaseCompaniesInfos{}).Where("alias = ?", company_name).First(&companyInfo)

	models.DB.Model(&models.FreightBase{}).Where(
		"year = ? and port = ? and transportation_mode = ? and target_addr = ?",
		YEAR, portInfo.PortName, trans, companyInfo.TargetAddr).First(&freight)

	if portInfo.ExtraPay != "" {
		models.DB.Model(&models.FreightBase{}).Where(
			"year = ? and port = ? and transportation_mode = ? and target_addr = ?",
			YEAR, portInfo.PortName, trans, portInfo.ExtraPay).First(&extraPay)
	}

	return FreightInfo{
		Freight:  freight.Price,
		ExtraPay: extraPay.Price,
	}

}

func CreateCostCalculation(data *models.DynamicExcelTable, tableName string) (string, error) {

	var err error

	tempFile, err := utils.DownloadFromNas(data.NasFileName)
	uploadDir := utils.GetUploadTmpDir()
	ext := filepath.Ext(data.NasFileName)
	newFilePath := filepath.Join(uploadDir, uuid.New().String()+ext)

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	f, err := excelize.OpenFile(tempFile)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	port := data.Datas.BaseData["port"]
	count, err := strconv.ParseFloat(strings.TrimSpace(data.Datas.BaseData["count"]), 64)
	productName := data.Datas.BaseData["product_name"]
	totalRow, err := strconv.ParseInt(data.Datas.BaseData["total_1"], 10, 64)

	startRow, err := strconv.ParseInt(data.Datas.List[0]["__ROW_INDEX__"], 10, 64)
	endRow, err := strconv.ParseInt(data.Datas.List[len(data.Datas.List)-1]["__ROW_INDEX__"], 10, 64)

	productInfo, err := GetProductInfo(productName)
	// 箱量
	totalContainers := decimal.NewFromFloat(count).Div(decimal.NewFromFloat(productInfo.ContainerTypeWeight)).Round(0).IntPart()

	totalZhengXiang, err := GetZhengXiangTotal(data, productInfo)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	sheets := f.GetSheetList()
	sheetName := sheets[0]

	clearanceFeeUnitPrice := GetClearanceFeeByPort(port, models.ClearanceFee, productInfo.ContainerType)
	shortHaulFeeUnitPrice := GetClearanceFeeByPort(port, models.ShortHaulFee, productInfo.ContainerType)
	unpackingFeeUnitPrice := GetClearanceFeeByPort(port, models.UnpackingFee, productInfo.ContainerType)

	err = f.SetCellValue(sheetName, "E5", "上海翊帆")
	err = f.SetCellFormula(sheetName, "I6", fmt.Sprintf("=%d*%f", totalContainers, clearanceFeeUnitPrice))
	err = f.SetCellFormula(sheetName, "K6", fmt.Sprintf("=(%d-%d)*%f", totalContainers, totalZhengXiang, shortHaulFeeUnitPrice))
	err = f.SetCellFormula(sheetName, "M6", fmt.Sprintf("=(%d-%d)*%f", totalContainers, totalZhengXiang, unpackingFeeUnitPrice))

	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("J%d", totalRow+1),
		"",
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("L%d", totalRow+1),
		"",
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("M%d", totalRow+1),
		"",
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("N%d", totalRow+1),
		"",
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("O%d", totalRow+1),
		"",
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("P%d", totalRow+1),
		"",
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("Q%d", totalRow+1),
		"",
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("R%d", totalRow+1),
		"",
	)

	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("J%d", totalRow+1),
		fmt.Sprintf("=SUM(J%d:J%d)", startRow+1, endRow+1),
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("L%d", totalRow+1),
		fmt.Sprintf("=SUM(L%d:L%d)", startRow+1, endRow+1),
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("M%d", totalRow+1),
		fmt.Sprintf("=SUM(M%d:M%d)", startRow+1, endRow+1),
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("N%d", totalRow+1),
		fmt.Sprintf("=SUM(N%d:N%d)", startRow+1, endRow+1),
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("O%d", totalRow+1),
		fmt.Sprintf("=SUM(O%d:O%d)", startRow+1, endRow+1),
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("P%d", totalRow+1),
		fmt.Sprintf("=SUM(P%d:P%d)", startRow+1, endRow+1),
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("Q%d", totalRow+1),
		fmt.Sprintf("=SUM(Q%d:Q%d)", startRow+1, endRow+1),
	)
	err = f.SetCellFormula(sheetName,
		fmt.Sprintf("R%d", totalRow+1),
		fmt.Sprintf("=SUM(R%d:R%d)", startRow+1, endRow+1),
	)

	f.UpdateLinkedValue()

	var readRule models.ExcelReadRules

	q := models.DB.Preload("IterateRule").Model(&models.ExcelReadRules{}).Where("dynamic_table_name = ?", tableName).First(&readRule)

	if q.Error != nil {
		log.Fatal(q.Error)
		return "", q.Error
	}

	for i := 0; i < len(data.Datas.List); i++ {
		item := data.Datas.List[i]
		var planned_count float64
		var rowIndex int64
		var fewer_packages float64
		var pkgCount string
		var __ROW_INDEX__ int64

		planned_count, err = strconv.ParseFloat(item["planned_count"], 64)
		__ROW_INDEX__, err = strconv.ParseInt(item["__ROW_INDEX__"], 10, 64)
		rowIndex = __ROW_INDEX__ + 1

		if item["fewer_packages"] != "" {
			fewer_packages, err = strconv.ParseFloat(item["fewer_packages"], 64)
			pkgCount = fmt.Sprintf("少%d包", int(math.Abs(fewer_packages)))
		}

		// 箱数
		containerCount := GetContainerCount(planned_count, productInfo.ContainerTypeWeight)
		realCount := (planned_count*1000 - math.Abs(fewer_packages)*productInfo.PackingSpecification) / 1000

		freightInfo := GetUnitFreight(port, containerCount, item["company_name"])

		if containerCount != "" {
			f.SetCellFormula(sheetName,
				fmt.Sprintf("P%d", rowIndex),
				fmt.Sprintf("=(%f+%f)*E%d", freightInfo.Freight, freightInfo.ExtraPay, rowIndex),
			)
		} else {
			f.SetCellFormula(sheetName,
				fmt.Sprintf("P%d", rowIndex),
				fmt.Sprintf("=(%f+%f)*J%d", freightInfo.Freight, freightInfo.ExtraPay, rowIndex),
			)
		}
		f.SetCellFormula(sheetName,
			fmt.Sprintf("Q%d", rowIndex),
			fmt.Sprintf("=P%d/1.09", rowIndex),
		)

		f.SetCellFormula(sheetName,
			fmt.Sprintf("R%d", rowIndex),
			fmt.Sprintf("=Q%d/J%d", rowIndex, rowIndex),
		)

		f.SetSheetRow(sheetName, fmt.Sprintf("E%d", rowIndex), &[]interface{}{
			containerCount,
		})
		f.SetSheetRow(sheetName, fmt.Sprintf("I%d", rowIndex), &[]interface{}{
			pkgCount,
			realCount,
		})

	}

	totalNoTaxClearanceFee, err := f.CalcCellValue(sheetName, "I7")
	totalNoTaxShortHaulFee, err := f.CalcCellValue(sheetName, "K7")
	totalNoTaxUnpackingFee, err := f.CalcCellValue(sheetName, "M7")
	realSendCount, err := f.CalcCellValue(sheetName, fmt.Sprintf("J%d", totalRow+1))

	fmt.Println("totalNoTaxClearanceFee", totalNoTaxClearanceFee)
	fmt.Println("totalNoTaxShortHaulFee", totalNoTaxShortHaulFee)
	fmt.Println("totalNoTaxUnpackingFee", totalNoTaxUnpackingFee)
	fmt.Println("realSendCount", realSendCount)

	for i := 0; i < len(data.Datas.List); i++ {
		item := data.Datas.List[i]
		//var planned_count float64
		var __ROW_INDEX__ int64
		//var fewer_packages float64
		//var pkgCount string

		//planned_count, err = strconv.ParseFloat(item["planned_count"], 64)
		__ROW_INDEX__, err = strconv.ParseInt(item["__ROW_INDEX__"], 10, 64)

		rowIndex := __ROW_INDEX__ + 1

		f.SetCellFormula(sheetName,
			fmt.Sprintf("L%d", rowIndex),
			fmt.Sprintf("=%s/%s*J%d", totalNoTaxClearanceFee, realSendCount, rowIndex),
		)
		f.SetCellFormula(sheetName,
			fmt.Sprintf("M%d", rowIndex),
			fmt.Sprintf("=%s/%s*J%d", totalNoTaxShortHaulFee, realSendCount, rowIndex),
		)
		f.SetCellFormula(sheetName,
			fmt.Sprintf("N%d", rowIndex),
			fmt.Sprintf("=%s/%s*J%d", totalNoTaxUnpackingFee, realSendCount, rowIndex),
		)

	}

	f.UpdateLinkedValue()

	defer f.Close()

	if err = f.SaveAs(newFilePath); err != nil {
		log.Fatal(err)
		return "", err
	}
	return newFilePath, nil
}
