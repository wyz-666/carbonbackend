package model

import "gorm.io/gorm"

type BuyerTx struct {
	gorm.Model
	TID      string `gorm:"column:tid;type:varchar(64)" json:"tid"`
	Uuid     string `gorm:"column:uuid" json:"uuid"`
	Project  string `gorm:"column:project" json:"project"`
	Type     string `gorm:"column:type" json:"type"`
	Price    string `gorm:"column:price" json:"price"`
	TxVolume string `gorm:"column:txVolume" json:"txVolume"`
	State    int    `gorm:"column:state" json:"state"`
}

type SellerTx struct {
	gorm.Model
	TID      string `gorm:"column:tid;type:varchar(64)" json:"tid"`
	Uuid     string `gorm:"column:uuid" json:"uuid"`
	Project  string `gorm:"column:project" json:"project"`
	Type     string `gorm:"column:type" json:"type"`
	Price    string `gorm:"column:price" json:"price"`
	TxVolume string `gorm:"column:txVolume" json:"txVolume"`
	State    int    `gorm:"column:state" json:"state"`
}

type Notition struct {
	gorm.Model
	NID         string `gorm:"column:nid;type:varchar(64)" json:"nid"`
	TID         string `gorm:"column:tid" json:"tid"`
	Uuid        string `gorm:"column:uuid" json:"uuid"`
	Type        string `gorm:"column:type" json:"type"`
	UserName    string `gorm:"column:userName" json:"userName"`
	CompanyName string `gorm:"column:companyName" json:"companyName"`
	Phone       string `gorm:"column:phone" json:"phone"`
	Email       string `gorm:"column:email" json:"email"`
	State       int    `gorm:"column:state" json:"state"`
}

type Board struct {
	gorm.Model
	Date    string `gorm:"column:date" json:"date"`
	Content string `gorm:"column:content" json:"content"`
}
