package model

import "gorm.io/gorm"

type SemiMonthQuotation struct {
	gorm.Model
	QID            string `gorm:"column:qid;type:varchar(64);uniqueIndex" json:"qid"`
	Uuid           string `gorm:"column:uuid" json:"uuid"`
	Product        string `gorm:"column:product" json:"product"`
	Type           string `gorm:"column:type" json:"type"`
	LowerPrice     string `gorm:"column:lowerPrice" json:"lowerPrice"`
	HigherPrice    string `gorm:"column:higherPrice" json:"higherPrice"`
	Price          string `gorm:"column:price" json:"price"`
	TxVolume       string `gorm:"column:txVolume" json:"txVolume"`
	ApplicableTime string `gorm:"column:applicableTime" json:"applicableTime"`
	Approved       bool   `gorm:"column:approved" json:"approved"`
}

type MonthQuotation struct {
	gorm.Model
	QID            string `gorm:"column:qid;type:varchar(64);uniqueIndex" json:"qid"`
	Uuid           string `gorm:"column:uuid" json:"uuid"`
	Product        string `gorm:"column:product" json:"product"`
	Type           string `gorm:"column:type" json:"type"`
	LowerPrice     string `gorm:"column:lowerPrice" json:"lowerPrice"`
	HigherPrice    string `gorm:"column:higherPrice" json:"higherPrice"`
	Price          string `gorm:"column:price" json:"price"`
	TxVolume       string `gorm:"column:txVolume" json:"txVolume"`
	ApplicableTime string `gorm:"column:applicableTime" json:"applicableTime"`
	Approved       bool   `gorm:"column:approved" json:"approved"`
}

type YearQuotation struct {
	gorm.Model
	QID            string `gorm:"column:qid;type:varchar(64);uniqueIndex" json:"qid"`
	Uuid           string `gorm:"column:uuid" json:"uuid"`
	Product        string `gorm:"column:product" json:"product"`
	Type           string `gorm:"column:type" json:"type"`
	LowerPrice     string `gorm:"column:lowerPrice" json:"lowerPrice"`
	HigherPrice    string `gorm:"column:higherPrice" json:"higherPrice"`
	Price          string `gorm:"column:price" json:"price"`
	TxVolume       string `gorm:"column:txVolume" json:"txVolume"`
	ApplicableTime string `gorm:"column:applicableTime" json:"applicableTime"`
	SubmitTime     string `gorm:"column:submitTime" json:"submitTime"`
	Approved       bool   `gorm:"column:approved" json:"approved"`
}

type MonthQuotationWithUser struct {
	MonthQuotation
	UserName string `json:"userName"`
}

type YearQuotationWithUser struct {
	YearQuotation
	UserName string `json:"userName"`
}
