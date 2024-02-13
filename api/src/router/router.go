package router

import (
	"api/controller"

	"github.com/labstack/echo"
)

func NewRouter(c controller.IUserAndEntryController) *echo.Echo {
	e := echo.New()
	e.POST("/", c.HandleUserAndEntry)
	return e
}
