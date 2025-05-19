package request

type ReqUser struct {
	Account         string `json:"account" form:"account"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	Type            int    `json:"type" form:"type"`
	Name            string `json:"name" form:"name"`
	Company         string `json:"company" form:"company"`
	CompanyType     string `json:"companyType" form:"companyType"`
	Phone           string `json:"phone" form:"phone"`
	Email           string `json:"email" form:"email"`
}

type ReqLogin struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
}

type ReqApproveUser struct {
	Uuid string `json:"uuid" form:"uuid"`
}

type ReqUpdateUserType struct {
	Uuid string `json:"uuid" form:"uuid"`
	Type string `json:"type" form:"type"`
}
