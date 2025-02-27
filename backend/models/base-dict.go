package models

import "gorm.io/gorm"

type BaseDict struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// 字典名称
	Key string `gorm:"type:varchar(100); comment:字典名称" json:"key"`
	// 字典值
	Value string `gorm:"type:varchar(100); comment:字典值" json:"value"`
	// 字典类型
	Type string `gorm:"type:varchar(100); comment:字典类型" json:"type"`
}
