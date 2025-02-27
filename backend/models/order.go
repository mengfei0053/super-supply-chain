package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber string  `gorm:"type:varchar(100); comment:订单号; unique; not null" json:"orderNumber"`
	ArrivalPort string  `gorm:"type:varchar(100); comment:到达港口" json:"arrivalPort"`
	ArrivalDate string  `gorm:"type:varchar(100); comment:到达日期" json:"arrivalDate"`
	Count       float64 `gorm:"type:int; comment:数量" json:"count"`
	// 单位
	Unit string `gorm:"type:varchar(100); comment:单位" json:"unit"`
	// 货值USD
	ValueUsd float64 `gorm:"type:float; comment:货值USD" json:"valueUsd"`
	// 货值RMB
	ValueRmb float64 `gorm:"type:float; comment:货值RMB" json:"valueRmb"`
	// 品名
	GoodsName string `gorm:"type:varchar(100); comment:品名" json:"goodsName"`

	// 清关费单价
	ClearanceFeeUnitPrice float64 `gorm:"type:float; comment:清关费单价" json:"clearanceFeeUnitPrice"`
	// 短驳费单价
	ShortHaulFeeUnitPrice float64 `gorm:"type:float; comment:短驳费单价" json:"shortHaulFeeUnitPrice"`
	// 掏箱费单价
	UnpackingFeeUnitPrice float64 `gorm:"type:float; comment:掏箱费单价" json:"unpackingFeeUnitPrice"`
	// 仓储费单价
	StorageFeeUnitPrice float64 `gorm:"type:float; comment:仓储费单价" json:"storageFeeUnitPrice"`
	// 出入库费单价
	InOutFeeUnitPrice float64 `gorm:"type:float; comment:出入库费单价" json:"inOutFeeUnitPrice"`
	// 其他费用单价
	OtherFeeUnitPrice float64 `gorm:"type:float; comment:其他费用单价" json:"otherFeeUnitPrice"`
	// 清关费单位
	ClearanceFeeUnit string `gorm:"type:varchar(100); comment:清关费单位" json:"clearanceFeeUnit"`
	// 短驳费单位
	ShortHaulFeeUnit string `gorm:"type:varchar(100); comment:短驳费单位" json:"shortHaulFeeUnit"`
	// 掏箱费单位
	UnpackingFeeUnit string `gorm:"type:varchar(100); comment:掏箱费单位" json:"unpackingFeeUnit"`
	// 仓储费单位
	StorageFeeUnit string `gorm:"type:varchar(100); comment:仓储费单位" json:"storageFeeUnit"`
	// 出入库费单位
	InOutFeeUnit string `gorm:"type:varchar(100); comment:出入库费单位" json:"inOutFeeUnit"`
	// 其他费用单位
	OtherFeeUnit string `gorm:"type:varchar(100); comment:其他费用单位" json:"otherFeeUnit"`
	// 清关费
	ClearanceFee float64 `gorm:"type:float; comment:清关费" json:"clearanceFee"`
	// 短驳费
	ShortHaulFee float64 `gorm:"type:float; comment:短驳费" json:"shortHaulFee"`
	// 掏箱费
	UnpackingFee float64 `gorm:"type:float; comment:掏箱费" json:"unpackingFee"`
	// 仓储费
	StorageFee float64 `gorm:"type:float; comment:仓储费" json:"storageFee"`
	// 出入库费
	InOutFee float64 `gorm:"type:float; comment:出入库费" json:"inOutFee"`
	// 其他费用
	OtherFee float64 `gorm:"type:float; comment:其他费用" json:"otherFee"`
	// 关联文件Id
	UploadFileId uint `gorm:"type:int; comment:关联文件Id; not null" json:"uploadFileId"`
}
