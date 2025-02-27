package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	DB.AutoMigrate(
		&BaseAccountsInfos{},
		&BaseCompaniesInfos{},
		&Order{},
		&ShippingOrder{},
		&ExcelMappingRules{},
		&ExcelReadRules{},
		&ExcelExportTemplates{},
		&FreightBase{},
		&ClearancePriceBase{},
		&BaseDict{},
		&UploadFile{})

}
