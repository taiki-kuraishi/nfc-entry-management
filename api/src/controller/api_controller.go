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
	uu usecase.IUserUsecase
	eu usecase.IEntryUsecase
}

type Response struct {
	UserMessage  string `json:"user_message"`
	EntryMessage string `json:"entry_message"`
}

type EntryRequest struct {
	StudentNumber uint    `json:"student_number"`
	Name          string  `json:"name"`
	Timestamp     float64 `json:"timestamp"`
}

func NewApiController(uu usecase.IUserUsecase, eu usecase.IEntryUsecase) IApiController {
	return &ApiController{uu, eu}
}

func (ac *ApiController) RootController(c echo.Context) error {
	request := EntryRequest{}
	if err := c.Bind(&request); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// convert float64 to time.Time
	seconds := int64(request.Timestamp)
	nanoseconds := int64((request.Timestamp - float64(seconds)) * 1e9)
	timestamp := time.Unix(seconds, nanoseconds)

	user := model.User{
		StudentNumber: request.StudentNumber,
		Name:          request.Name,
		CreatedAt:     timestamp,
		UpdatedAt:     timestamp,
	}

	userMessage, err := ac.uu.CreateOrUpdateUser(user)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	entryMessage, err := ac.eu.EntryOrExit(request.StudentNumber, timestamp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := Response{
		UserMessage:  userMessage,
		EntryMessage: entryMessage,
	}

	return c.JSON(http.StatusOK, response)
}
