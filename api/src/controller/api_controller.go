package controller

import (
	"api/model"
	"api/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type IUserAndEntryController interface {
	HandleUserAndEntry(c echo.Context) error
}

type UserAndEntryController struct {
	uu usecase.IUserUsecase
	eu usecase.IEntryUsecase
}

type Response struct {
	UserMessage  string `json:"user_message"`
	EntryMessage string `json:"entry_message"`
}

type Request struct {
	StudentNumber uint    `json:"student_number"`
	Name          string  `json:"name"`
	Timestamp     float64 `json:"timestamp"`
}

func NewUserAndEntryController(uu usecase.IUserUsecase, eu usecase.IEntryUsecase) IUserAndEntryController {
	return &UserAndEntryController{uu, eu}
}

func (ac *UserAndEntryController) HandleUserAndEntry(c echo.Context) error {
	request := Request{}
	if err := c.Bind(&request); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if request.StudentNumber == 0 {
		fmt.Println("student number is required")
		return c.JSON(http.StatusBadRequest, "student number is required")
	}
	if request.Name == "" {
		return c.JSON(http.StatusBadRequest, "name is required")
	}
	if request.Timestamp == 0 {
		return c.JSON(http.StatusBadRequest, "timestamp is required")
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
