package request

type ReqMarket struct {
	Product      string `json:"product" form:"product"`
	Date         string `json:"date" form:"date"`
	LowerPrice   string `json:"lowerPrice" form:"lowerPrice"`
	HigherPrice  string `json:"higherPrice" form:"higherPrice"`
	ClosingPrice string `json:"closingPrice" form:"closingPrice"`
}

type ReqExpectation struct {
	Product     string  `json:"product" form:"product"`
	Type        string  `json:"type" form:"type"`
	Date        string  `json:"date" form:"date"`
	LowerPrice  float64 `json:"lowerPrice" form:"lowerPrice"`
	HigherPrice float64 `json:"higherPrice" form:"higherPrice"`
	MidPrice    float64 `json:"midPrice" form:"midPrice"`
}

type ReqGECExpectation struct {
	Product    string  `json:"product" form:"product"`
	Type       string  `json:"type" form:"type"`
	Date       string  `json:"date" form:"date"`
	Price      float64 `json:"price" form:"price"`
	PriceIndex float64 `json:"priceIndex" form:"priceIndex"`
}
