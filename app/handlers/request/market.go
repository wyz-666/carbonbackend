package request

type ReqMarket struct {
	Product      string `json:"product" form:"product"`
	Date         string `json:"date" form:"date"`
	LowerPrice   string `json:"lowerPrice" form:"lowerPrice"`
	HigherPrice  string `json:"higherPrice" form:"higherPrice"`
	ClosingPrice string `json:"closingPrice" form:"closingPrice"`
}

type ReqExpectation struct {
	Product     string `json:"product" form:"product"`
	Type        string `json:"type" form:"type"`
	Date        string `json:"date" form:"date"`
	LowerPrice  string `json:"lowerPrice" form:"lowerPrice"`
	HigherPrice string `json:"higherPrice" form:"higherPrice"`
	MidPrice    string `json:"midPrice" form:"midPrice"`
}

type ReqGECExpectation struct {
	Product    string `json:"product" form:"product"`
	Type       string `json:"type" form:"type"`
	Date       string `json:"date" form:"date"`
	Price      string `json:"price" form:"price"`
	PriceIndex string `json:"priceIndex" form:"priceIndex"`
}
