package models

import "gorm.io/gorm"

type ExcelReadRules struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	MenuName         string `gorm:"type:varchar(100); comment:规则名称; unique; not null" json:"menuName"`
	SheetIndex       uint   `gorm:"type:int; comment:表格索引" json:"sheetIndex"`
	DynamicTableName string `gorm:"type:varchar(100); comment:动态表名称 unique; not null" json:"dynamicTableName"`
	Desc             string `gorm:"type:varchar(100); comment:规则描述" json:"desc"`
	// MapRule  映射 ExcelMappingRules
	MapRuleId     uint              `json:"mapRuleId"`
	IterateRuleId uint              `json:"iterateRuleId"`
	MapRule       ExcelMappingRules `gorm:"foreignKey:MapRuleId;references:ID" json:"mapRule"`
	IterateRule   ExcelMappingRules `gorm:"foreignKey:IterateRuleId;references:ID" json:"iterateRule"`
}
