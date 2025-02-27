package models

import "gorm.io/gorm"

type ExcelExportTemplates struct {
	gorm.Model
	ID              uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Alias           string `gorm:"type:varchar(100); comment:模板别名; unique; not null" json:"alias"`
	FileName        string `gorm:"type:varchar(100); comment:文件名; unique; not null" json:"fileName"`
	UploadFilePath  string `gorm:"type:varchar(400); comment:文件路径; unique; not null" json:"uploadFilePath"`
	AssociatedTable string `gorm:"type:varchar(100); comment:关联表" json:"associatedTable"`
}
