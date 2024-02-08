package validator

import (
	"api/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IUserValidator interface {
	UserValidation(user model.User) error
}

type UserValidator struct{}

func NewUserValidator() IUserValidator {
	return &UserValidator{}
}

func (uv *UserValidator) UserValidation(user model.User) error {
	return validation.ValidateStruct(&user,
		validation.Field(
			&user.StudentNumber,
			validation.Required.Error("student number is required"),
			validation.Min(uint(10000000)).Error("student number must be greater than 20121000"),
			validation.Max(uint(39999999)).Error("student number must be less than 20130000"),
		),
		validation.Field(
			&user.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(2, 32).Error("name must be between 3 and 50 characters"),
		),
		validation.Field(
			&user.CreatedAt,
			validation.Required.Error("created at is required"),
			validation.Min(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).Error("created at must be after 2024"),
			validation.Max(time.Now()).Error("created at must be before now"),
		),
		validation.Field(
			&user.UpdatedAt,
			validation.Required.Error("updated at is required"),
			validation.Min(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).Error("updated at must be after 2024"),
			validation.Max(time.Now()).Error("updated at must be before now"),
		),
	)
}
