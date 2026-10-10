package models

import "gorm.io/gorm"

type BaseDict struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// 字典名称（港口原文，如 上海洋山）
	Key string `gorm:"type:varchar(100); comment:字典名称" json:"key"`
	// 字典值（起始地城市，如 上海）
	Value string `gorm:"type:varchar(100); comment:字典值" json:"value"`
	// 字典类型（港口字典）
	Type string `gorm:"type:varchar(100); comment:字典类型" json:"type"`
	// 口岸名（运费/清关价查询用，如 上海口岸）
	PortName string `gorm:"column:port_name;type:varchar(100); comment:口岸名称" json:"port_name"`
	// 补差项名称（运费补差 target_addr，空表示无补差）
	ExtraPay string `gorm:"column:extra_pay;type:varchar(100); comment:补差项" json:"extra_pay"`
}
