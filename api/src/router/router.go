package router

import (
	"api/controller"

	"github.com/labstack/echo"
)

func NewRouter(uc controller.IApiController) *echo.Echo {
	e := echo.New()
	e.POST("/", uc.RootController)
	return e
}
