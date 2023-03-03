package server

import (
	"byoj/router"
	"byoj/utils/logs"
	"strconv"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Server struct {
	Port int `json:"port"`
}

func Run(s Server) error {
	e := echo.New()

	// Router and middleware
	router.Load(e)

	err := e.Start(":" + strconv.Itoa(s.Port))
	if err != nil {
		logs.Error("Server run failed at port "+strconv.Itoa(s.Port), zap.Error(err))
		return err
	}
	return nil
}
