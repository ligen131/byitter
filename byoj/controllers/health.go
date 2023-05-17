package controllers

import (
	"byoj/utils/logs"

	"github.com/labstack/echo"
)

func HealthGET(c echo.Context) error {
	logs.Debug("GET /health")

	return ResponseOK(c, "ok")
}
