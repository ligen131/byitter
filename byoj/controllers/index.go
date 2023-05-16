package controllers

import (
	"byoj/utils/logs"
	"net/http"

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
	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data: link{
			Link: documentLink{
				Doc: "// Document link here.",
			},
		},
	})
}
