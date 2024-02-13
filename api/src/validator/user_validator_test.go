package validator_test

import (
	"api/model"
	"api/validator"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserValidator_UserValidation(t *testing.T) {

	var err error
	TimeValidationMin := int64(1704034800)

	t.Setenv("STUDENT_NUMBER_MIN", "10000000")
	t.Setenv("STUDENT_NUMBER_MAX", "39999999")
	t.Setenv("NAME_MIN_LENGTH", "2")
	t.Setenv("NAME_MAX_LENGTH", "32")
	t.Setenv("TIME_VALIDATION_MIN", strconv.FormatInt(TimeValidationMin, 10))

	validUser := model.User{
		StudentNumber: uint(20122027),
		Name:          "カイシ　タロウ",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	uv := validator.IUserValidator(validator.NewUserValidator())

	//Test case 1 正しいケース
	t.Run("valid", func(t *testing.T) {
		err = uv.UserValidation(validUser)
		assert.NoError(t, err)
	})

	//Test case 2 STUDENT_NUMBER_MINの値が無効
	t.Run("invalid_STUDENT_NUMBER_MIN", func(t *testing.T) {
		t.Setenv("STUDENT_NUMBER_MIN", "")
		err = uv.UserValidation(validUser)
		assert.Equal(t, "strconv.ParseUint: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 3 STUDENT_NUMBER_MAXの値が無効
	t.Run("invalid_STUDENT_NUMBER_MAX", func(t *testing.T) {
		t.Setenv("STUDENT_NUMBER_MAX", "")
		err = uv.UserValidation(validUser)
		assert.Equal(t, "strconv.ParseUint: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 4 NAME_MIN_LENGTHの値が無効
	t.Run("invalid_NAME_MIN_LENGTH", func(t *testing.T) {
		t.Setenv("NAME_MIN_LENGTH", "")
		err = uv.UserValidation(validUser)
		assert.Equal(t, "strconv.Atoi: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 5 NAME_MAX_LENGTHの値が無効
	t.Run("invalid_NAME_MAX_LENGTH", func(t *testing.T) {
		t.Setenv("NAME_MAX_LENGTH", "")
		err = uv.UserValidation(validUser)
		assert.Equal(t, "strconv.Atoi: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 6 TIME_VALIDATION_MINの値が無効
	t.Run("invalid_TIME_VALIDATION_MIN", func(t *testing.T) {
		t.Setenv("TIME_VALIDATION_MIN", "")
		err = uv.UserValidation(validUser)
		assert.Equal(t, "strconv.ParseInt: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 7 StudentNumberが欠損している
	t.Run("not_required_StudentNumber", func(t *testing.T) {
		user := model.User{
			Name:      validUser.Name,
			CreatedAt: validUser.CreatedAt,
			UpdatedAt: validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "student number is required"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 8 StudentNumberが最小値よりも小さい
	t.Run("below_minimum_StudentNumber", func(t *testing.T) {
		user := model.User{
			StudentNumber: 9999999,
			Name:          validUser.Name,
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "student number must be greater than 10000000"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 9 StudentNumberが最大値よりも大きい
	t.Run("above_maximum_StudentNumber", func(t *testing.T) {
		user := model.User{
			StudentNumber: 40000000,
			Name:          validUser.Name,
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "student number must be less than 39999999"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	// Test case 10 StudentNumberが最小値
	t.Run("minimum_StudentNumber", func(t *testing.T) {
		user := model.User{
			StudentNumber: 10000000,
			Name:          validUser.Name,
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})

	// Test case 11 StudentNumberが最大値
	t.Run("maximum_StudentNumber", func(t *testing.T) {
		user := model.User{
			StudentNumber: 39999999,
			Name:          validUser.Name,
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})

	//Test case 12 Nameが欠損している
	t.Run("not_required_Name", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "name is required"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 13 Nameが最小値よりも小さい
	t.Run("below_minimum_Name", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          "カ",
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "name must be between 2 and 32 characters"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 14 Nameが最大値よりも大きい
	t.Run("above_maximum_Name", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          "あああああああああああああああああああああああああああああああああ", //33 characters
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "name must be between 2 and 32 characters"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 15 Nameが最小値
	t.Run("minimum_Name", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          "カイ",
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})

	//Test case 16 Nameが最大値
	t.Run("maximum_Name", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          "ああああああああああああああああああああああああああああああああ", //31 characters
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})

	//Test case 17 CreatedAtが欠損している
	t.Run("not_required_CreatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "created at is required"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 18 CreatedAtが最小値よりも小さい
	t.Run("below_minimum_CreatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     time.Unix(TimeValidationMin-1, 0),
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "created at must be after " + time.Unix(TimeValidationMin, 0).String()
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 19 CreatedAtが最大値よりも大きい
	t.Run("above_maximum_CreatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     time.Now().Add(time.Second),
			UpdatedAt:     time.Now().Add(time.Second),
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "created at must be before " + time.Now().Round(time.Second).String()
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 20 CreatedAtが最小値
	t.Run("minimum_CreatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     time.Unix(TimeValidationMin, 0),
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})

	//Test case 21 CreatedAtが最大値
	t.Run("maximum_CreatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     time.Now().Add(-time.Second),
			UpdatedAt:     validUser.UpdatedAt,
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})

	//Test case 22 UpdatedAtが欠損している
	t.Run("not_required_UpdatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     validUser.CreatedAt,
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "updated at is required"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 23 UpdatedAtが最小値よりも小さい
	t.Run("below_minimum_UpdatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     time.Unix(TimeValidationMin-1, 0),
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "updated at must be after " + time.Unix(TimeValidationMin, 0).String()
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 24 UpdatedAtが最大値よりも大きい
	t.Run("above_maximum_UpdatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     validUser.CreatedAt,
			UpdatedAt:     time.Now().Add(time.Second),
		}
		err = uv.UserValidation(user)
		exceptedErrorMessages := "updated at must be before " + time.Now().Round(time.Second).String()
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 25 UpdatedAtが最小値
	t.Run("minimum_UpdatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     time.Unix(TimeValidationMin, 0),
			UpdatedAt:     time.Unix(TimeValidationMin, 0),
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})

	//Test case 26 UpdatedAtが最大値
	t.Run("maximum_UpdatedAt", func(t *testing.T) {
		user := model.User{
			StudentNumber: validUser.StudentNumber,
			Name:          validUser.Name,
			CreatedAt:     time.Now().Add(-time.Second),
			UpdatedAt:     time.Now().Add(-time.Second),
		}
		err = uv.UserValidation(user)
		assert.NoError(t, err)
	})
}
