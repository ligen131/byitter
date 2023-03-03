package router

import (
	"byoj/controllers"

	"github.com/labstack/echo"
)

func Load(e *echo.Echo) {
	routes(e)
}

func routes(e *echo.Echo) {
	e.GET("/", controllers.IndexGET)
}
