package models

import (
	"gorm.io/gorm"
)

type ShippingOrder struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement"`
	// 订单号
	OrderNumber string `gorm:"type:varchar(100); comment:订单号; unique; not null"`
	// 目标公司
	TargetCompany string `gorm:"type:varchar(100); comment:目标公司"`
	// 数量
	Count int `gorm:"type:int; comment:数量"`
	// 单位
	Unit string `gorm:"type:varchar(100); comment:单位"`
	// 单价
	Price float64 `gorm:"type:float; comment:单价"`
	// 短驳费
	ShortHaulFee float64 `gorm:"type:float; comment:短驳费"`
	// 箱数
	BoxCount int `gorm:"type:int; comment:箱数"`
	// 运输方式
	Transportation string `gorm:"type:varchar(100); comment:运输方式"`
	// 金额
	Amount float64 `gorm:"type:float; comment:金额"`
}
