package models

import "gorm.io/gorm"

// FreightBase 运费基础表
type FreightBase struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement"`
	//进出口类型
	ImportExportType string `gorm:"type:varchar(100); comment:进出口类型"`
	// 港口
	Port string `gorm:"type:varchar(100); comment:港口"`
	// 送货基地
	TargetBase string `gorm:"type:varchar(100); comment:送货基地"`
	// 货物名称
	GoodsName string `gorm:"type:varchar(100); comment:货物名称"`
	// 里程
	Mileage float64 `gorm:"type:float; comment:里程"`
	// 运输形式
	TransportationMode string `gorm:"type:varchar(100); comment:运输形式"`
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
	// 年份
	Year string `gorm:"type:varchar(100); comment:年份"`
}
