package router

import (
	"byoj/controllers"
	"byoj/controllers/middleware"

	"github.com/labstack/echo"
	echomw "github.com/labstack/echo/middleware"
)

func Load(e *echo.Echo) {
	routes(e)
}

func routes(e *echo.Echo) {
	e.Use(echomw.Recover())

	e.GET("/", controllers.IndexGET)

	e.GET("/health", controllers.HealthGET)

	userGroup := e.Group("/user")
	{
		userGroup.GET("", controllers.UserGET)
		userGroup.GET("/", controllers.UserGET)
		userGroup.POST("/register", controllers.UserRegisterPOST)
		userGroup.POST("/login", controllers.UserLoginPOST)
		userGroup.GET("/isauth", controllers.UserIsAuthGET, middleware.TokenVerificationMiddleware)
	}
}
