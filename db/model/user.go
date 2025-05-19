package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uuid         string `gorm:"column:uuid;type:varchar(64);uniqueIndex" json:"uuid"`
	UserID       string `gorm:"column:userId" json:"userId"`
	UserName     string `gorm:"column:userName" json:"userName"`
	CompanyName  string `gorm:"column:companyName" json:"companyName"`
	CompanyType  string `gorm:"column:companyType" json:"companyType"`
	Type         int    `gorm:"column:type" json:"type"`
	Phone        string `gorm:"column:phone" json:"phone"`
	Email        string `gorm:"column:email" json:"email"`
	Account      string `gorm:"column:account" json:"account"`
	PasswordHash string `gorm:"column:passwordhash" json:"passwordhash"`
	Approved     bool   `gorm:"column:approved" json:"approved"`
}
