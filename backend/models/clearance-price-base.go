package models

import "gorm.io/gorm"

// 清关价格基础表
type ClearancePriceBase struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement"`
	// 进出口类型
	ImportExportType string `gorm:"type:varchar(100); comment:进出口类型"`
	// 招标口岸
	Port string `gorm:"type:varchar(100); comment:招标口岸"`
	// 集装箱类型
	ContainerType string `gorm:"type:varchar(100); comment:集装箱类型"`
	// 费用描述
	CostDescription string `gorm:"type:varchar(100); comment:费用描述"`
	// 单位
	Unit string `gorm:"type:varchar(100); comment:单位"`
	// 含税单价
	Price float64 `gorm:"type:float; comment:含税单价"`
	// 不含税单价
	PriceWithoutTax float64 `gorm:"type:float; comment:不含税单价"`
	// 税金
	Tax float64 `gorm:"type:float; comment:税金"`
	// 税率
	TaxRate float64 `gorm:"type:float; comment:税率"`
}
