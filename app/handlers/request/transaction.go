package request

type ReqTransaction struct {
	Uuid     string `json:"uuid" form:"uuid"`
	Project  string `json:"project" form:"project"`
	Type     string `json:"type" form:"type"`
	Price    string `json:"price" form:"price"`
	TxVolume string `json:"txVolume" form:"txVolume"`
}

type ReqNotition struct {
	Uuid string `json:"uuid" form:"uuid"`
	Tid  string `json:"tid" form:"tid"`
	Type string `json:"type" form:"type"`
}

type ReqBoard struct {
	Date    string `json:"date" form:"date"`
	Content string `json:"content" form:"content"`
}
