package validator

import (
	"api/model"
	"os"
	"strconv"
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
    studentNumberMin, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MIN"), 10, 64)
    if err != nil {
        return err
    }

    studentNumberMax, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MAX"), 10, 64)
    if err != nil {
        return err
    }

    nameMinLength, err := strconv.Atoi(os.Getenv("NAME_MIN_LENGTH"))
    if err != nil {
        return err
    }

    nameMaxLength, err := strconv.Atoi(os.Getenv("NAME_MAX_LENGTH"))
    if err != nil {
        return err
    }

    TimeValidationMin, err := strconv.ParseInt(os.Getenv("TIME_VALIDATION_MIN"), 10, 64)
    if err != nil {
        return err
    }

    err = validation.Validate(
        &user.StudentNumber,
        validation.Required.Error("student number is required"),
        validation.Min(studentNumberMin).Error("student number must be greater than 10000000"),
        validation.Max(studentNumberMax).Error("student number must be less than 39999999"),
    )
    if err != nil {
        return err
    }

    err = validation.Validate(
        &user.Name,
        validation.Required.Error("name is required"),
        validation.RuneLength(nameMinLength, nameMaxLength).Error("name must be between 2 and 32 characters"),
    )
    if err != nil {
        return err
    }

    err = validation.Validate(
        &user.CreatedAt,
        validation.Required.Error("created at is required"),
        validation.Min(time.Unix(TimeValidationMin, 0)).Error("created at must be after " + time.Unix(TimeValidationMin, 0).String()),
        validation.Max(time.Now()).Error("created at must be before "+time.Now().Round(time.Second).String()),
    )
    if err != nil {
        return err
    }

    err = validation.Validate(
        &user.UpdatedAt,
        validation.Required.Error("updated at is required"),
        validation.Min(user.CreatedAt).Error("updated at must be after " + time.Unix(TimeValidationMin, 0).String()),
        validation.Max(time.Now()).Error("updated at must be before "+time.Now().Round(time.Second).String()),
    )
    return err
}
