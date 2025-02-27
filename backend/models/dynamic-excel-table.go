package models

import "gorm.io/gorm"

type ExcelData struct {
	BaseData map[string]string   `json:"baseData"`
	List     []map[string]string `json:"list"`
}

type DynamicExcelTable struct {
	gorm.Model
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UploadFilePath string    `gorm:"type:varchar(100); comment:文件路径; unique; not null" json:"uploadFilePath"`
	FileName       string    `gorm:"type:varchar(100); comment:文件名称; unique; not null" json:"fileName"`
	Datas          ExcelData `gorm:"type:varchar(10000); serializer:json; comment:数据" json:"datas"`
}
