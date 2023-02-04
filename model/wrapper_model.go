package model

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseList struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
	Message string      `json:"message"`
}

type ResponseData struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Token struct {
	Expired string `json:"expired"`
	Token   string `json:"token"`
}
