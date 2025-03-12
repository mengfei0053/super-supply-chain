package utils

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"super-supply-chain/models"
)

func GetUploadTmpDir() string {
	HOME := os.Getenv("HOME")
	tmpPath := HOME + "/tmp"

	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.Mkdir(tmpPath, os.ModePerm)
	}
	return tmpPath
}

func GetToTalRowIndexs(rows [][]string, base map[string]string) {
	totalCount := 1

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]
		if len(row) > 0 {
			for j := 0; j < len(row); j++ {
				cell := row[j]
				if cell == "合计" {
					base[fmt.Sprint("total_", totalCount)] = strconv.Itoa(rowIndex)
					totalCount += 1
				}
			}
		}
	}
}

func GetExcelData(path string, mapRules MapRuleInfo, tableName string) (models.ExcelData, error) {
	var data = models.ExcelData{}
	var base = make(map[string]string)

	f, err := excelize.OpenFile(path)
	sheetIndex := mapRules.SheetIndex
	mapRule := mapRules.MapRule
	IterateRule := mapRules.IterateRule

	if err != nil {
		log.Fatal(err)
		return data, err
	}
	defer func() (map[string]string, error) {
		if err := f.Close(); err != nil {
			log.Fatal(err)
			return base, err
		}
		return base, nil
	}()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return data, errors.New("No sheets found")
	}

	rows, err := f.GetRows(sheets[sheetIndex])
	if err != nil {
		log.Fatal(err)
		return data, err
	}

	targetCols := make(map[int]string)
	for i := 0; i < 26; i++ {
		targetCols[i] = string('A' + i)
	}
	targetColIndex := make(map[string]int)
	for i := 0; i < 26; i++ {
		targetColIndex[string('A'+i)] = i
	}
	isCustomDeclation := tableName == "dynamic_customs_declaration_form"
	isIntegrity := tableName == "dynamic_Integrity_packaging_invoice"

	var res = make([]map[string]string, 0)
	startRow := IterateRule.StartRow
	for rowIndex := startRow - 1; rowIndex < len(rows); rowIndex++ {
		rowData := rows[rowIndex]
		if len(rowData) == 0 {
			break
		}

		if rowIndex >= startRow-1 {
			var chooseRowData = make(map[string]string)
			for _, mappingRule := range IterateRule.Rules {
				ExcelKey := mappingRule.ExcelKey
				JsonKey := mappingRule.JsonKey
				cellIndex := targetColIndex[ExcelKey]
				if cellIndex <= len(rowData)-1 {
					cellValue := rowData[cellIndex]
					chooseRowData[JsonKey] = cellValue
				} else {
					chooseRowData[JsonKey] = ""
				}
			}
			chooseRowData["__ROW_INDEX__"] = strconv.Itoa(rowIndex)
			res = append(res, chooseRowData)
		}

	}

	for _, rule := range mapRule.Rules {
		value, err := f.GetCellValue(sheets[sheetIndex], rule.ExcelKey)
		if err != nil {
			log.Fatal(err)
			return data, err
		}
		if isCustomDeclation {
			if value != "" {
				parts := strings.Split(value, "\n")
				base[rule.JsonKey] = parts[1]
			} else {
				base[rule.JsonKey] = value
			}

		} else if isIntegrity {
			if value != "" {
				parts := strings.Split(value, ":")
				parts2 := strings.Split(value, "：")
				if len(parts) == 2 {
					base[rule.JsonKey] = strings.TrimSpace(parts[1])
				} else if len(parts2) == 2 {
					base[rule.JsonKey] = strings.TrimSpace(parts2[1])
				} else {
					base[rule.JsonKey] = value
				}

			} else {
				base[rule.JsonKey] = value
			}
		} else {
			base[rule.JsonKey] = value
		}

	}

	GetToTalRowIndexs(rows, base)

	data.BaseData = base
	data.List = res

	return data, nil
}
