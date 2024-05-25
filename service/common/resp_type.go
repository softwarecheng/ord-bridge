package common

type BaseResp struct {
	Code int    `json:"code" example:"0"`
	Msg  string `json:"msg" example:"ok"`
}

type ListResp struct {
	Start int `json:"start" example:"0"`
	Total int `json:"total" example:"9992"`
}
