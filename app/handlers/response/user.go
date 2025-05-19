package response

type ResLogin struct {
	Uuid     string `json:uuid`
	UserID   string `json:"userId"`
	UserType int    `json:userType`
	Token    string `json:"token"`
}
