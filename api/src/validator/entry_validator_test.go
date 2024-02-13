package validator_test

import (
	"api/model"
	"api/validator"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntryValidator_StudentNumberValidation(t *testing.T) {
	var err error
	StudentNumberMin := 10000000
	StudentNumberMax := 39999999
	sampleStudentNumber := uint(20122027)

	t.Setenv("STUDENT_NUMBER_MIN", strconv.Itoa(StudentNumberMin))
	t.Setenv("STUDENT_NUMBER_MAX", strconv.Itoa(StudentNumberMax))

	ev := validator.IEntryValidator(validator.NewEntryValidator())

	//Test case 1 正しいケース
	t.Run("valid", func(t *testing.T) {
		err = ev.StudentNumberValidation(sampleStudentNumber)
		assert.NoError(t, err)
	})

	//Test case 2 STUDENT_NUMBER_MINが"無効
	t.Run("invalid_STUDENT_NUMBER_MIN", func(t *testing.T) {
		t.Setenv("STUDENT_NUMBER_MIN", "")
		err = ev.StudentNumberValidation(sampleStudentNumber)
		assert.Equal(t, "strconv.ParseUint: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 3 STUDENT_NUMBER_MAXが無効
	t.Run("invalid_STUDENT_NUMBER_MAX", func(t *testing.T) {
		t.Setenv("STUDENT_NUMBER_MAX", "")
		err = ev.StudentNumberValidation(sampleStudentNumber)
		assert.Equal(t, "strconv.ParseUint: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 4 StudentNumberが欠損している
	t.Run("not_required", func(t *testing.T) {
		err = ev.StudentNumberValidation(0)
		exceptedErrorMessages := "student number is required"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 5 StudentNumberが最小値よりも小さい
	t.Run("below_minimum", func(t *testing.T) {
		err = ev.StudentNumberValidation(uint(StudentNumberMin - 1))
		exceptedErrorMessages := "student number must be greater than 10000000"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 6 StudentNumberが最大値よりも大きい
	t.Run("above_maximum", func(t *testing.T) {
		err = ev.StudentNumberValidation(uint(StudentNumberMax + 1))
		exceptedErrorMessages := "student number must be less than 39999999"
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 7 StudentNumberが最小値
	t.Run("minimum", func(t *testing.T) {
		err = ev.StudentNumberValidation(uint(StudentNumberMin))
		assert.NoError(t, err)
	})

	//Test case 8 最大値のStudentNumber
	t.Run("maximum", func(t *testing.T) {
		err = ev.StudentNumberValidation(uint(StudentNumberMax))
		assert.NoError(t, err)
	})
}

func TestEntryValidator_EntryValidation(t *testing.T) {
	var err error
	StudentNumberMin := 10000000
	StudentNumberMax := 39999999
	TimeValidationMin := int64(1704034800)
	sampleStudentNumber := uint(20122027)

	t.Setenv("STUDENT_NUMBER_MIN", strconv.Itoa(StudentNumberMin))
	t.Setenv("STUDENT_NUMBER_MAX", strconv.Itoa(StudentNumberMax))
	t.Setenv("TIME_VALIDATION_MIN", strconv.FormatInt(TimeValidationMin, 10))

	ev := validator.IEntryValidator(validator.NewEntryValidator())

	validEntry := model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      nil,
		StudentNumber: sampleStudentNumber,
	}

	//Test case 1 正しいケース
	t.Run("valid", func(t *testing.T) {
		err = ev.EntryValidation(validEntry)
		assert.NoError(t, err)
	})

	//Test case 2 STUDENT_NUMBER_MINが無効
	t.Run("invalid_STUDENT_NUMBER_MIN", func(t *testing.T) {
		t.Setenv("STUDENT_NUMBER_MIN", "")
		err = ev.EntryValidation(validEntry)
		assert.Equal(t, "strconv.ParseUint: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 3 STUDENT_NUMBER_MAXが無効
	t.Run("invalid_STUDENT_NUMBER_MAX", func(t *testing.T) {
		t.Setenv("STUDENT_NUMBER_MAX", "")
		err = ev.EntryValidation(validEntry)
		assert.Equal(t, "strconv.ParseUint: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 4 TIME_VALIDATION_MINが無効
	t.Run("invalid_TIME_VALIDATION_MIN", func(t *testing.T) {
		t.Setenv("TIME_VALIDATION_MIN", "")
		err = ev.EntryValidation(validEntry)
		assert.Equal(t, "strconv.ParseInt: parsing \"\": invalid syntax", err.Error())
	})

	//Test case 5 EntryTimeが欠損している
	t.Run("not_required_EntryTime", func(t *testing.T) {
		entry := model.Entry{
			ExitTime:      nil,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "entry_time: entry time is required."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 6 EntryTimeが最小値よりも小さい
	t.Run("below_minimum_EntryTime", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Unix(TimeValidationMin-1, 0),
			ExitTime:      nil,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "entry_time: must be after " + time.Unix(TimeValidationMin, 0).String() + "."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 7 EntryTimeが最大値よりも大きい
	t.Run("above_maximum_EntryTime", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Now().Add(time.Second),
			ExitTime:      nil,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "entry_time: must be before " + time.Now().Round(time.Second).String() + "."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 8 EntryTimeが最小値
	t.Run("minimum_EntryTime", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Unix(TimeValidationMin, 0),
			ExitTime:      nil,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		assert.NoError(t, err)
	})

	//Test case 9 EntryTimeが最大値
	t.Run("maximum_EntryTime", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Now(),
			ExitTime:      nil,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		assert.NoError(t, err)
	})

	//Test case 10 StudentNumberが欠損している
	t.Run("not_required_StudentNumber", func(t *testing.T) {
		entry := model.Entry{
			EntryTime: time.Now(),
			ExitTime:  nil,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "student_number: student number is required."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 11 StudentNumberが最小値よりも小さい
	t.Run("below_minimum_StudentNumber", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Now(),
			ExitTime:      nil,
			StudentNumber: 9999999,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "student_number: student number must be greater than 10000000."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 12 StudentNumberが最大値よりも大きい
	t.Run("above_maximum_StudentNumber", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Now(),
			ExitTime:      nil,
			StudentNumber: 40000000,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "student_number: student number must be less than 39999999."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 13 StudentNumberが最小値
	t.Run("minimum_StudentNumber", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Now(),
			ExitTime:      nil,
			StudentNumber: 10000000,
		}
		err = ev.EntryValidation(entry)
		assert.NoError(t, err)
	})

	//Test case 14 StudentNumberが最大値
	t.Run("maximum_StudentNumber", func(t *testing.T) {
		entry := model.Entry{
			EntryTime:     time.Now(),
			ExitTime:      nil,
			StudentNumber: 39999999,
		}
		err = ev.EntryValidation(entry)
		assert.NoError(t, err)
	})

	//Test case 15 ExitTimeが最小値よりも小さい
	t.Run("below_minimum_ExitTime", func(t *testing.T) {
		exitTime := time.Now().Add(-time.Second)
		entry := model.Entry{
			EntryTime:     time.Now(),
			ExitTime:      &exitTime,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "exit_time: must be after " + entry.EntryTime.String() + "."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 16 ExitTimeが最大値よりも大きい
	t.Run("above_maximum_ExitTime", func(t *testing.T) {
		exitTime := time.Now().Add(time.Second)
		entry := model.Entry{
			EntryTime:     time.Now(),
			ExitTime:      &exitTime,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		exceptedErrorMessages := "exit_time: must be before " + time.Now().Round(time.Second).String() + "."
		assert.Equal(t, exceptedErrorMessages, err.Error())
	})

	//Test case 17 ExitTimeが最小値
	t.Run("minimum_ExitTime", func(t *testing.T) {
		exitTime := time.Unix(TimeValidationMin, 0)
		entry := model.Entry{
			EntryTime:     exitTime,
			ExitTime:      &exitTime,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		assert.NoError(t, err)
	})

	//Test case 18 ExitTimeが最大値
	t.Run("maximum_ExitTime", func(t *testing.T) {
		exitTime := time.Now().Round(time.Second)
		entry := model.Entry{
			EntryTime:     time.Unix(TimeValidationMin, 0),
			ExitTime:      &exitTime,
			StudentNumber: sampleStudentNumber,
		}
		err = ev.EntryValidation(entry)
		assert.NoError(t, err)
	})

}
