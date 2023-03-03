package controllers

import (
	"byoj/utils/logs"
	"net/http"

	"github.com/labstack/echo"
)

func IndexGET(c echo.Context) error {
	logs.Debug("GET /")
	return c.String(http.StatusOK, "Hello, the world!")
}
