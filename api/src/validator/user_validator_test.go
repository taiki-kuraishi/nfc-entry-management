package validator

import (
	"api/model"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserValidator(t *testing.T) {
	TimeValidationMin, err := strconv.ParseInt(os.Getenv("TIME_VALIDATION_MIN"), 10, 64)
	assert.NoError(t, err)

	sampleStudentNumber := uint(20122027)
	sampleName := "カイシ　タロウ"
	sampleTime := time.Now()

	uv := &UserValidator{}

	//Test case 1 valid user
	user := model.User{
		StudentNumber: 20122027,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 2 invalid StudentNumber (not required)
	user = model.User{
		Name:      sampleName,
		CreatedAt: sampleTime,
		UpdatedAt: sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages := "student number is required"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 3 invalid StudentNumber (just below minimum value)
	user = model.User{
		StudentNumber: 9999999,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "student number must be greater than 10000000"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 4 invalid StudentNumber (just above maximum value)
	user = model.User{
		StudentNumber: 40000000,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "student number must be less than 39999999"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 5 valid StudentNumber (minimum allowed value)
	user = model.User{
		StudentNumber: 10000000,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 6 valid StudentNumber (maximum allowed value)
	user = model.User{
		StudentNumber: 39999999,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 7 invalid Name (not required)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "name is required"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 8 invalid Name (just below minimum length)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          "カ",
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "name must be between 2 and 32 characters"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 9 invalid Name (just above maximum length)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          "あああああああああああああああああああああああああああああああああ", //33 characters
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "name must be between 2 and 32 characters"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 10 valid Name (minimum allowed length)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          "カイ",
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 11 valid Name (maximum allowed length)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          "ああああああああああああああああああああああああああああああああ", //31 characters
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 12 invalid CreatedAt (not required)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "created at is required"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 13 invalid CreatedAt (just below minimum value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     time.Unix(TimeValidationMin -1 , 0),
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "created at must be after " + time.Unix(TimeValidationMin, 0).String()
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 14 invalid CreatedAt (just above maximum value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     time.Now().Add(time.Second),
		UpdatedAt:     time.Now().Add(time.Second),
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "created at must be before " + time.Now().Round(time.Second).String()
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 15 valid CreatedAt (minimum allowed value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     time.Unix(TimeValidationMin, 0),
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 16 valid CreatedAt (maximum allowed value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     time.Now().Add(-time.Second),
		UpdatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 17 invalid UpdatedAt (not required)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     sampleTime,
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "updated at is required"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 18 invalid UpdatedAt (just below minimum value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     time.Unix(TimeValidationMin -1 , 0),
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "updated at must be after " + time.Unix(TimeValidationMin, 0).String()
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 19 invalid UpdatedAt (just above maximum value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     time.Now().Add(time.Second),
	}
	err = uv.UserValidation(user)
	exceptedErrorMessages = "updated at must be before " + time.Now().Round(time.Second).String()
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 20 valid UpdatedAt (minimum allowed value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     time.Unix(TimeValidationMin, 0),
		UpdatedAt:     time.Unix(TimeValidationMin, 0),
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)

	//Test case 21 valid UpdatedAt (maximum allowed value)
	user = model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     time.Now().Add(-time.Second),
		UpdatedAt:     time.Now().Add(-time.Second),
	}
	err = uv.UserValidation(user)
	assert.NoError(t, err)
}
