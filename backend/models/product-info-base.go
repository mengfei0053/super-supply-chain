package models

type ProductInfoBase struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// 产品名称
	ProductName string `gorm:"type:varchar(100); comment:产品名称" json:"product_name"`
	// 产品型号
	ProductModel string `gorm:"type:varchar(100); comment:产品型号" json:"product_model"`
	// 集装箱规格
	ContainerType int `gorm:"type:int; comment:集装箱规格" json:"container_type"`
	// 集装箱规格
	ContainerTypeWeight float64 `gorm:"type:float; comment:集装箱规格" json:"container_type_weight"`
	// 集装箱规格单位
	ContainerTypeUnit string `gorm:"type:varchar(100); comment:集装箱规格单位" json:"container_type_unit"`
	// 包规格
	PackingSpecification float64 `gorm:"type:float; comment:包规格" json:"packing_specification"`
	// 包规格单位
	PackingSpecificationUnit string `gorm:"type:varchar(100); comment:包规格单位" json:"packing_specification_unit"`
}
