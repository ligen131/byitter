package controllers

type ErrorMessage struct {
	Message string `json:"msg"`
	Err     string `json:"err"`
}

type StatusMessage struct {
	Status string `json:"status"`
}

type ResponseStruct struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}
