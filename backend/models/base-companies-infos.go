package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseCompaniesInfos struct {
	gorm.Model
	ID                      uint   `gorm:"primaryKey;autoIncrement"`
	Name                    string `gorm:"type:varchar(255);unique;not null;comment:公司名称"`
	TargetAddr              string `gorm:"type:varchar(255);comment:发票目标地址"`
	Alias                   string `gorm:"type:varchar(255);comment:公司别名"`
	AddrCountry             string `gorm:"type:varchar(100)"`
	AddrProvince            string `gorm:"type:varchar(100)"`
	AddrCity                string `gorm:"type:varchar(100)"`
	AddrStreet              string `gorm:"type:varchar(255)"`
	UnifiedSocialCreditCode string `gorm:"type:varchar(100);unique;not null;comment:统一社会信用代码"`
	BankCode                string `gorm:"type:varchar(100);comment:银行代码"`
	PhoneNum                string `gorm:"type:varchar(20)"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               time.Time
}
