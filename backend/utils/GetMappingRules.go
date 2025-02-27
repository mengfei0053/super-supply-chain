package utils

import (
	"fmt"
	"super-supply-chain/models"
)

type MapRuleInfo struct {
	MapRule     models.ExcelMappingRules
	IterateRule models.ExcelMappingRules
	SheetIndex  uint
}

func GetMappingRules(tableName string) (MapRuleInfo, error) {
	result := models.ExcelReadRules{}
	query := models.DB.Preload("MapRule").Preload("IterateRule").Model(&models.ExcelReadRules{}).Where("dynamic_table_name = ?", tableName).First(&result)
	if query.Error != nil {
		fmt.Println(query.Error)
		return MapRuleInfo{}, query.Error
	}

	return MapRuleInfo{
		MapRule:     result.MapRule,
		IterateRule: result.IterateRule,
		SheetIndex:  result.SheetIndex,
	}, nil
}
