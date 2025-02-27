package models

import (
	"gorm.io/gorm"
)

type MappingRule struct {
	ExcelKey string `json:"excelKey"`
	JsonKey  string `json:"jsonKey"`
	Desc     string `json:"desc"`
}

type ExcelMappingRules struct {
	gorm.Model
	ID       uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string        `gorm:"type:varchar(100); comment:规则名称; unique; not null" json:"name"`
	Type     string        `gorm:"type:varchar(100); comment:规则类型" json:"type"`
	StartRow int           `gorm:"type:int; comment:开始行" json:"startRow"`
	Rules    []MappingRule `gorm:"type:varchar(1000); serializer:json; comment:规则" json:"rules"`
}
