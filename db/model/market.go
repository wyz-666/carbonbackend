package model

import "gorm.io/gorm"

type CEAMarket struct {
	gorm.Model
	Date         string `gorm:"column:date;" json:"date"`
	LowerPrice   string `gorm:"column:lowerPrice" json:"lowerPrice"`
	HigherPrice  string `gorm:"column:higherPrice" json:"higherPrice"`
	ClosingPrice string `gorm:"column:closingPrice" json:"closingPrice"`
}

type CCERMarket struct {
	gorm.Model
	Date         string `gorm:"column:date;" json:"date"`
	ClosingPrice string `gorm:"column:closingPrice" json:"closingPrice"`
}

type CEAMonthExpectation struct {
	Date             string  `gorm:"column:date;type:varchar(64);uniqueIndex" json:"date"`
	LowerPrice       float64 `gorm:"column:lowerPrice" json:"lowerPrice"`
	HigherPrice      float64 `gorm:"column:higherPrice" json:"higherPrice"`
	MidPrice         float64 `gorm:"column:midPrice" json:"midPrice"`
	LowerPriceIndex  float64 `gorm:"column:lowerPriceIndex" json:"lowerPriceIndex"`
	HigherPriceIndex float64 `gorm:"column:higherPriceIndex" json:"higherPriceIndex"`
	MidPriceIndex    float64 `gorm:"column:midPriceIndex" json:"midPriceIndex"`
}
type CCERMonthExpectation struct {
	Date             string  `gorm:"column:date;type:varchar(64);uniqueIndex" json:"date"`
	LowerPrice       float64 `gorm:"column:lowerPrice" json:"lowerPrice"`
	HigherPrice      float64 `gorm:"column:higherPrice" json:"higherPrice"`
	MidPrice         float64 `gorm:"column:midPrice" json:"midPrice"`
	LowerPriceIndex  float64 `gorm:"column:lowerPriceIndex" json:"lowerPriceIndex"`
	HigherPriceIndex float64 `gorm:"column:higherPriceIndex" json:"higherPriceIndex"`
	MidPriceIndex    float64 `gorm:"column:midPriceIndex" json:"midPriceIndex"`
}
type CEAYearExpectation struct {
	Date             string  `gorm:"column:date;type:varchar(64);uniqueIndex" json:"date"`
	LowerPrice       float64 `gorm:"column:lowerPrice" json:"lowerPrice"`
	HigherPrice      float64 `gorm:"column:higherPrice" json:"higherPrice"`
	MidPrice         float64 `gorm:"column:midPrice" json:"midPrice"`
	LowerPriceIndex  float64 `gorm:"column:lowerPriceIndex" json:"lowerPriceIndex"`
	HigherPriceIndex float64 `gorm:"column:higherPriceIndex" json:"higherPriceIndex"`
	MidPriceIndex    float64 `gorm:"column:midPriceIndex" json:"midPriceIndex"`
}

type GECMonthExpectation struct {
	Product    string  `gorm:"column:product" json:"product"`
	Type       string  `gorm:"column:type" json:"type"`
	Date       string  `gorm:"column:date" json:"date"`
	Price      float64 `gorm:"column:price" json:"price"`
	PriceIndex float64 `gorm:"column:priceIndex" json:"priceIndex"`
}
