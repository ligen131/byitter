package controllers

import (
	"byoj/utils/logs"
	"net/http"

	"github.com/labstack/echo"
)

func HealthGET(c echo.Context) error {
	logs.Debug("GET /health")
	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    "ok",
	})
}
