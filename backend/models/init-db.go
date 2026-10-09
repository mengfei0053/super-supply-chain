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

	// 已有业务表继续手工管理。个人访问令牌是新表，启动时单独迁移。
	// 数据库账号没有建表权限时，执行 ssc-sqls/personal_access_tokens.sql。
	if err = DB.AutoMigrate(&PersonalAccessToken{}); err != nil {
		fmt.Println("warning: personal_access_tokens AutoMigrate failed:", err.Error())
		fmt.Println("apply ssc-sqls/personal_access_tokens.sql if the table is missing")
	}

	// 迁移 schema
	//DB.AutoMigrate(
	//	&BaseAccountsInfos{},
	//	&BaseCompaniesInfos{},
	//	&Order{},
	//	&ShippingOrder{},
	//	&ExcelExportTemplates{},
	//	&FreightBase{},
	//	&ClearancePriceBase{},
	//	&BaseDict{},
	//	&ProductInfoBase{},
	//	&ExcelReadRuleInfos{},
	//	&UploadFile{})
}
