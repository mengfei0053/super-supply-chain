package utils

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"log"
	"regexp"
	"strconv"
	"strings"
	"super-supply-chain/models"
)

func GetPkgCount(input string) (string, error) {
	// 正则匹配模式：提取冒号后的数字
	re := regexp.MustCompile(`集装箱标箱数及号码：(\d+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", errors.New("未找到目标字段GetPkgCount")
	}

	return matches[1], nil
}

func GetProductName(input string) (string, error) {
	parts := strings.Split(input, "\n")
	if len(parts) == 3 {
		productName := strings.TrimSpace(parts[0])
		if productName != "" {
			fields := strings.Fields(productName)
			if len(fields) > 0 {
				return fields[len(fields)-1], nil
			}
		}
	}
	return "", errors.New("未找到目标字段GetProductName")
}

func SliceContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if strings.Contains(s, str) {
			return true
		}
	}
	return false
}

func isRowDataEmpty(slice []string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			return false
		}
	}

	return true
}

func GetPrice(input string) (string, string, error) {
	var units = map[string]string{
		"美元": "USD",
		"欧元": "EUR",
	}
	parts := strings.Split(input, "\n")
	if len(parts) == 3 {
		price := strings.TrimSpace(parts[1])
		unit := strings.TrimSpace(parts[2])

		if price != "" && units[unit] != "" {
			return price, units[unit], nil
		}
	}
	return "", "", errors.New("未找到目标字段GetPrice")
}

func GetBaoguanDan(path string, tableName string) (models.ExcelData, error) {
	var data = models.ExcelData{}
	var err error
	var base = make(map[string]string)
	mapRules := models.ExcelReadRuleInfos{}

	q := models.DB.Model(&models.ExcelReadRuleInfos{}).Where("dynamic_table_name = ?", tableName).First(&mapRules)
	if q.Error != nil {
		log.Fatal(q.Error)
		return data, q.Error
	}

	f, err := excelize.OpenFile(path)

	if err != nil {
		log.Fatal(err)
		return data, err
	}
	defer func() error {
		if err := f.Close(); err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return data, err
	}

	mapRule := mapRules.Rules.MapRule

	var xiangmuIndex = -1

	for index, row := range rows {
		for _, cell := range row {
			if cell != "" {
				parts := strings.Split(cell, "\n")
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])

					for _, rule := range mapRule {
						// 如果key 包含 rule.ExcelKey
						if strings.Contains(key, rule.ExcelKey) {
							if rule.ExcelKey == "净重" {
								base["weight_unit"] = key
							}
							base[rule.JsonKey] = value
						}
					}
					if strings.Contains(key, "集装箱标箱数及号码") {
						var pkg_count string
						pkg_count, err = GetPkgCount(value)
						base["pkg_count"] = pkg_count
					}
				}
				if len(parts) == 3 {
					p3 := strings.TrimSpace(parts[2])
					if p3 == "美元" || p3 == "欧元" {
						var price string
						var price_unit string
						price, price_unit, err = GetPrice(cell)
						base["price"] = price
						base["price_unit"] = price_unit
					}
				}

			}
		}
		if SliceContainsString(row, "项号") && SliceContainsString(row, "商品编号") && SliceContainsString(row, "商品名称及规格型号") {
			xiangmuIndex = index
		}
	}

	var info string
	var product_name string

	if xiangmuIndex != -1 {
		log.Println(rows[xiangmuIndex])
		infoIndex := xiangmuIndex + 2
		log.Println("A" + strconv.Itoa(infoIndex))
		info, err = f.GetCellValue("Sheet1", "A"+strconv.Itoa(infoIndex))
		product_name, err = GetProductName(info)
	} else {

		return data, errors.New("未找到目标字段GetProductName")
	}

	base["product_name"] = product_name

	data.BaseData = base

	return data, err
}
