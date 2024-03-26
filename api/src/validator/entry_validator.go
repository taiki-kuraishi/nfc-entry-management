package validator

import (
	"api/model"
	"os"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IEntryValidator interface {
	StudentNumberValidation(studentNumber uint) error
	EntryValidation(entry model.Entry) error
}

type EntryValidator struct{}

func NewEntryValidator() IEntryValidator {
	return &EntryValidator{}
}

func (ev *EntryValidator) StudentNumberValidation(studentNumber uint) error {
	studentNumberMin, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MIN"), 10, 64)
	if err != nil {
		return err
	}

	studentNumberMax, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MAX"), 10, 64)
	if err != nil {
		return err
	}

	return validation.Validate(studentNumber,
		validation.Required.Error("student number is required"),
		validation.Min(studentNumberMin).Error("student number must be greater than 10000000"),
		validation.Max(studentNumberMax).Error("student number must be less than 39999999"),
	)
}

type TimeAfter struct {
	required time.Time
}

func (t TimeAfter) Validate(value interface{}) error {
	if timeValue, ok := value.(*time.Time); ok {
		if timeValue.Before(t.required) {
			return validation.NewError("validation_min", "must be after "+t.required.String())
		}
	}
	return nil
}

type TimeBefore struct {
	required time.Time
}

func (t TimeBefore) Validate(value interface{}) error {
	if timeValue, ok := value.(*time.Time); ok {
		if timeValue.After(t.required) {
			return validation.NewError("validation_max", "must be before "+t.required.String())
		}
	}
	return nil
}

func (ev *EntryValidator) EntryValidation(entry model.Entry) error {
	studentNumberMin, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MIN"), 10, 64)
	if err != nil {
		return err
	}

	studentNumberMax, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MAX"), 10, 64)
	if err != nil {
		return err
	}

	// TimeValidationMin, err := strconv.ParseInt(os.Getenv("TIME_VALIDATION_MIN"), 10, 64)
	// if err != nil {
	// 	return err
	// }

	err = validation.ValidateStruct(&entry,
		// validation.Field(
		// 	&entry.EntryTime,
		// 	validation.Required.Error("entry time is required"),
		// 	validation.Min(time.Unix(TimeValidationMin, 0)).Error("must be after "+time.Unix(TimeValidationMin, 0).String()),
		// 	validation.Max(time.Now()).Error("must be before "+time.Now().Round(time.Second).String()),
		// ),
		validation.Field(
			&entry.StudentNumber,
			validation.Required.Error("student number is required"),
			validation.Min(studentNumberMin).Error("student number must be greater than 10000000"),
			validation.Max(studentNumberMax).Error("student number must be less than 39999999"),
		),
	)
	if err != nil {
		return err
	}

	// if entry.ExitTime != nil {
	// 	err = validation.ValidateStruct(&entry,
	// 		validation.Field(
	// 			&entry.ExitTime,
	// 			TimeAfter{required: entry.EntryTime},
	// 			TimeBefore{required: time.Now().Round(time.Second)},
	// 		),
	// 	)
	// }

	return err
}
