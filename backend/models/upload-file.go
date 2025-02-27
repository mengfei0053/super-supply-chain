package models

import "gorm.io/gorm"

type UploadFile struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FilePath       string `gorm:"type:varchar(100); comment:文件路径; unique; not null" json:"filePath"`
	FileName       string `gorm:"type:varchar(100); comment:文件名; unique; not null" json:"fileName"`
	OriginFileName string `gorm:"type:varchar(100); comment:原文件名; unique; not null" json:"originFileName"`
	FileSize       int64  `gorm:"type:int; comment:文件大小" json:"fileSize"`
}
