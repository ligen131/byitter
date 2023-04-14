package controllers

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Err     string `json:"err"`
}
