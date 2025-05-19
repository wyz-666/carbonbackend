package response

type ResUserScore struct {
	Uuid        string  `json:uuid`
	UserID      string  `json:"userId"`
	CompanyName string  `json:"companyName"`
	UserName    string  `json:"userName"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Score       float64 `json:"score"`
}
