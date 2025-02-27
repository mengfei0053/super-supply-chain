package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type BaseAccountsInfos struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Account   string `gorm:"unique;comment:账号"`
	Realname  string `gorm:"type:varchar(50)"`
	Password  string `gorm:"type:varchar(60)"`
	Email     string `gorm:"type:varchar(100)"`
	PhoneNum  string `gorm:"type:varchar(20)"`
	Avatar    string `gorm:"type:varchar(500)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"default:null"`
}

func (u *BaseAccountsInfos) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}
