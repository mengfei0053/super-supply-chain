package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DSN := fmt.Sprintf("%s:%s@tcp(%s)/super_supply_chain?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_SERVER"),
	)

	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	DB.AutoMigrate(
		&BaseAccountsInfos{},
		&BaseCompaniesInfos{},
		&Order{},
		&ShippingOrder{},
		&ExcelExportTemplates{},
		&FreightBase{},
		&ClearancePriceBase{},
		&BaseDict{},
		&ProductInfoBase{},
		&ExcelReadRuleInfos{},
		&UploadFile{})

}
