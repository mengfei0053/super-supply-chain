package models

import (
	"gorm.io/gorm"
)

type WholeOrBulkType string
type CostType string

const (
	Whole WholeOrBulkType = "whole" // 整箱
	Bulk  WholeOrBulkType = "bulk"  // 散货
)

const (
	// 清关费
	ClearanceFee CostType = "clearance_fee"
	// 短驳费
	ShortHaulFee CostType = "short_haul_fee"
	// 掏箱费
	UnpackingFee CostType = "unpacking_fee"
	// 装卸费
	LoadingAndUnloadingFee CostType = "loading_and_unloading_fee"
)

/*
*
20'普通柜塑料粒子
40'普通柜塑料粒子
20'普通柜食品（原料、预包装）
40'普通柜食品（原料、预包装）
20'普通柜机械设备仪器
40'普通柜机械设备仪器
20'冷冻柜食品原料
40'冷冻柜食品原料
保税区散货
空运或海运拼箱散货
20'普通柜食品饮料成品出口
40'普通柜食品饮料成品出口
*/
const (
	// 20'普通柜塑料粒子
	ContainerType20GPPlasticParticles int = iota + 1
	// 40'普通柜塑料粒子
	ContainerType40GPPlasticParticles
	// 20'普通柜食品（原料、预包装）
	ContainerType20GPFood
	// 40'普通柜食品（原料、预包装）
	ContainerType40GPFood
	// 20'普通柜机械设备仪器
	ContainerType20GPMachineryEquipment
	// 40'普通柜机械设备仪器
	ContainerType40GPMachineryEquipment
	// 20'冷冻柜食品原料
	ContainerType20GPFoodRawMaterials
	// 40'冷冻柜食品原料
	ContainerType40GPFoodRawMaterials
	// 保税区散货
	ContainerTypeBondedAreaBulk
	// 空运或海运拼箱散货
	ContainerTypeAirOrSeaBulk
	// 20'普通柜食品饮料成品出口
	ContainerType20GPFoodBeverage
	// 40'普通柜食品饮料成品出口
	ContainerType40GPFoodBeverage
)

var MapContainerType = map[string]int{
	"20'普通柜塑料粒子":       ContainerType20GPPlasticParticles,
	"40'普通柜塑料粒子":       ContainerType40GPPlasticParticles,
	"20'普通柜食品（原料、预包装）": ContainerType20GPFood,
	"40'普通柜食品（原料、预包装）": ContainerType40GPFood,
	"20'普通柜机械设备仪器":     ContainerType20GPMachineryEquipment,
	"40'普通柜机械设备仪器":     ContainerType40GPMachineryEquipment,
	"20'冷冻柜食品原料":       ContainerType20GPFoodRawMaterials,
	"40'冷冻柜食品原料":       ContainerType40GPFoodRawMaterials,
	"保税区散货":            ContainerTypeBondedAreaBulk,
	"空运或海运拼箱散货":        ContainerTypeAirOrSeaBulk,
	"20'普通柜食品饮料成品出口":   ContainerType20GPFoodBeverage,
	"40'普通柜食品饮料成品出口":   ContainerType40GPFoodBeverage,
}

// Slice
var SliceContainerType = []string{
	"20'普通柜塑料粒子",
	"40'普通柜塑料粒子",
	"20'普通柜食品（原料、预包装）",
	"40'普通柜食品（原料、预包装）",
	"20'普通柜机械设备仪器",
	"40'普通柜机械设备仪器",
	"20'冷冻柜食品原料",
	"40'冷冻柜食品原料",
	"保税区散货",
	"空运或海运拼箱散货",
	"20'普通柜食品饮料成品出口",
	"40'普通柜食品饮料成品出口",
}

func UpdateDb() {
	// 生成 CASE WHEN 表达式
	caseStmt := "CASE container_type "
	args := make([]interface{}, 0)
	for _, ct := range SliceContainerType {
		caseStmt += "WHEN ? THEN ? "
		args = append(args, ct, MapContainerType[ct])
	}
	caseStmt += "ELSE container_type END"

	// 执行批量更新
	DB.Model(&ClearancePriceBase{}).
		Where("container_type IN (?)", SliceContainerType).
		Update("container_type_enum", gorm.Expr(caseStmt, args...))

}

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
	// 集装箱类型
	ContainerTypeEnum int `gorm:"type:int; comment:集装箱类型"`
	// 整箱还是散货
	WholeOrBulk string `gorm:"type:ENUM('whole','bulk'); comment:整箱还是散货"`
	// 费用描述
	CostDescription string `gorm:"type:varchar(100); comment:费用描述"`
	// 费用描述
	CostType string `gorm:"type:ENUM('clearance_fee','short_haul_fee','unpacking_fee','loading_and_unloading_fee'); comment:费用类型"`
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
