package controller

import (
	"api/model"
	"api/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type IApiController interface {
	RootController(c echo.Context) error
}

type ApiController struct {
	uu       usecase.IUserUsecase
	location *time.Location
}

func NewApiController(uu usecase.IUserUsecase, location *time.Location) IApiController {
	return &ApiController{uu, location}
}

func (ac *ApiController) RootController(c echo.Context) error {
	request := model.EntryRequest{}
	if err := c.Bind(&request); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// convert float64 to time.Time
	seconds := int64(request.Timestamp)
	nanoseconds := int64((request.Timestamp - float64(seconds)) * 1e9)
	timestamp := time.Unix(seconds, nanoseconds).In(ac.location)

	user := model.User{
		StudentNumber: request.StudentNumber,
		Name:          request.Name,
		CreatedAt:     timestamp,
		UpdatedAt:     timestamp,
	}

	if err := ac.uu.CreateOrUpdateUser(user); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
