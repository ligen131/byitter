package controllers

import (
	"byoj/utils/logs"

	"github.com/labstack/echo"
)

type documentLink struct {
	Doc string `json:"document"`
}

type link struct {
	Link documentLink `json:"link"`
}

func IndexGET(c echo.Context) error {
	logs.Debug("GET /")

	return ResponseOK(c, link{
		Link: documentLink{
			Doc: "// Document link here.",
		},
	})
}
