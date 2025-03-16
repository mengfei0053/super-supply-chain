package models

import "gorm.io/gorm"

type MappingRule struct {
	ExcelKey string `json:"excelKey"`
	JsonKey  string `json:"jsonKey"`
	Desc     string `json:"desc"`
}

type IterateRule struct {
	StartRow int           `json:"startRow"`
	Rules    []MappingRule `json:"rules"`
}

type Rules struct {
	MapRule     []MappingRule `json:"mapRule"`
	IterateRule IterateRule   `json:"iterateRule"`
}

type ExcelReadRuleInfos struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SheetIndex       uint   `gorm:"type:int; comment:表格索引" json:"sheetIndex"`
	MenuName         string `gorm:"type:varchar(100); comment:规则名称; unique; not null" json:"menuName"`
	DynamicTableName string `gorm:"type:varchar(100); comment:动态表名称 unique; not null" json:"dynamicTableName"`
	Rules            Rules  `gorm:"type:text; serializer:json; comment:规则" json:"rules"`
}
