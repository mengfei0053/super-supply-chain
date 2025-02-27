package utils

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
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

func getRuleItem(key string, rules []models.MappingRule) *models.MappingRule {
	for _, itemRule := range rules {
		if itemRule.ExcelKey == key {
			return &itemRule
		}
	}
	return nil
}

func GetExcelData(path string, mapRules MapRuleInfo) (models.ExcelData, error) {
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

	var res = make([]map[string]string, 0)
	startRow := IterateRule.StartRow
	for rowIndex := startRow - 1; rowIndex < len(rows); rowIndex++ {
		rowData := rows[rowIndex]
		if len(rowData) == 0 {
			break
		}
		if len(rowData) > 0 && rowData[0] == "" {
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
				}
			}
			res = append(res, chooseRowData)
		}

	}

	for _, rule := range mapRule.Rules {
		value, err := f.GetCellValue(sheets[sheetIndex], rule.ExcelKey)
		if err != nil {
			log.Fatal(err)
			return data, err
		}
		base[rule.JsonKey] = value
	}
	data.BaseData = base
	data.List = res

	return data, nil
}
