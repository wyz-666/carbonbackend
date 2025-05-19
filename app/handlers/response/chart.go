package response

type MonthlyPriceStats struct {
	Month     string  `json:"month"`
	AvgLow    float64 `json:"avgLow"`
	AvgHigh   float64 `json:"avgHigh"`
	AvgPrice  float64 `json:"avgPrice"`
	LowIndex  float64 `json:"lowIndex"`
	HighIndex float64 `json:"highIndex"`
	MidIndex  float64 `json:"midIndex"`
	FitPrice  float64 `json:"fitPrice"`
}

type GECMonthlyPriceStats struct {
	Month    string  `json:"month"`
	Type     string  `json:"type"`
	AvgPrice float64 `json:"avgPrice"`
}

// CategoryStats 聚合结果结构体
type CategoryStats struct {
	Type        int     `json:"type"`        // 2 或 3
	CompanyType string  `json:"companyType"` // "控排企业" 或 "非控排企业"
	AvgLower    float64 `json:"avgLower"`
	AvgHigher   float64 `json:"avgHigher"`
	AvgBoth     float64 `json:"avgBoth"` // (AvgLower+AvgHigher)/2
}
