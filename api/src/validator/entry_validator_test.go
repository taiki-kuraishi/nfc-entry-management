package validator

import (
	"api/model"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntryValidator_StudentNumberValidation(t *testing.T) {
	StudentNumberMin, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MIN"), 10, 64)
	assert.NoError(t, err)

	StudentNumberMax, err := strconv.ParseUint(os.Getenv("STUDENT_NUMBER_MAX"), 10, 64)
	assert.NoError(t, err)

	sampleStudentNumber := uint(20122027)

	ev := &EntryValidator{}

	// Test case 1 valid student number
	err = ev.StudentNumberValidation(sampleStudentNumber)
	assert.NoError(t, err)

	// Test case 2 invalid student number (not required)
	err = ev.StudentNumberValidation(0)
	exceptedErrorMessages := "student number is required"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	// Test case 3 invalid student number (just below minimum value)
	err = ev.StudentNumberValidation(uint(StudentNumberMin - 1))
	exceptedErrorMessages = "student number must be greater than 10000000"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	// Test case 4 invalid student number (just above maximum value)
	err = ev.StudentNumberValidation(uint(StudentNumberMax + 1))
	exceptedErrorMessages = "student number must be less than 39999999"
	assert.Equal(t, exceptedErrorMessages, err.Error())

	// Test case 5 valid student number (minimum allowed value)
	err = ev.StudentNumberValidation(uint(StudentNumberMin))
	assert.NoError(t, err)

	// Test case 6 valid student number (maximum allowed value)
	err = ev.StudentNumberValidation(uint(StudentNumberMax))
	assert.NoError(t, err)

}

func TestEntryValidator_EntryValidation(t *testing.T) {
	TimeValidationMin, err := strconv.ParseInt(os.Getenv("TIME_VALIDATION_MIN"), 10, 64)
	assert.NoError(t, err)

	sampleStudentNumber := uint(20122027)

	ev := &EntryValidator{}

	// Test case 1 valid entry
	entry := model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      nil,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	assert.NoError(t, err)

	//Test case 2 invalid EntryTime	(not required ExitTime)
	entry = model.Entry{
		ExitTime:      nil,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages := "entry_time: entry time is required."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 3 invalid EntryTime	(EntryTime with just below minimum value)
	entry = model.Entry{
		EntryTime:     time.Unix(TimeValidationMin-1, 0),
		ExitTime:      nil,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages = "entry_time: must be after " + time.Unix(TimeValidationMin, 0).String() + "."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 4 invalid EntryTime	(EntryTime with just above maximum value)
	entry = model.Entry{
		EntryTime:     time.Now().Add(time.Second),
		ExitTime:      nil,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages = "entry_time: must be before " + time.Now().Round(time.Second).String() + "."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 5 invalid EntryTime	(Entry with minimum allowed value)
	entry = model.Entry{
		EntryTime:     time.Unix(TimeValidationMin, 0),
		ExitTime:      nil,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	assert.NoError(t, err)

	//Test case 6 valid EntryTime (EntryTime with maximum allowed value)
	entry = model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      nil,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	assert.NoError(t, err)

	//Test case 7 invalid StudentNumber	(StudentNumber not required)
	entry = model.Entry{
		EntryTime: time.Now(),
		ExitTime:  nil,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages = "student_number: student number is required."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 8 invalid StudentNumber	(StudentNumber with just below minimum value)
	entry = model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      nil,
		StudentNumber: 9999999,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages = "student_number: student number must be greater than 10000000."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 9 invalid StudentNumber	(StudentNumber with just above maximum value)
	entry = model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      nil,
		StudentNumber: 40000000,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages = "student_number: student number must be less than 39999999."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 10 valid StudentNumber (Student with minimum allowed value)
	entry = model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      nil,
		StudentNumber: 10000000,
	}
	err = ev.EntryValidation(entry)
	assert.NoError(t, err)

	//Test case 11 valid StudentNumber (Student with maximum allowed value)
	entry = model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      nil,
		StudentNumber: 39999999,
	}
	err = ev.EntryValidation(entry)
	assert.NoError(t, err)

	//Test case 12 invalid ExitTime	(ExitTime with just below minimum value)
	exitTime := time.Now().Add(-time.Second)
	entry = model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      &exitTime,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages = "exit_time: must be after " + entry.EntryTime.String() + "."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 13 valid ExitTime	(ExitTime with just above maximum value)
	exitTime = time.Now().Add(time.Second)
	entry = model.Entry{
		EntryTime:     time.Now(),
		ExitTime:      &exitTime,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	exceptedErrorMessages = "exit_time: must be before " + time.Now().Round(time.Second).String() + "."
	assert.Equal(t, exceptedErrorMessages, err.Error())

	//Test case 14 valid ExitTime	(ExitTime with minimum allowed value)
	exitTime = time.Unix(TimeValidationMin, 0)
	entry = model.Entry{
		EntryTime:     exitTime,
		ExitTime:      &exitTime,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	assert.NoError(t, err)

	//Test case 15 valid ExitTime	(ExitTime with maximum allowed value)
	exitTime = time.Now().Round(time.Second)
	entry = model.Entry{
		EntryTime:     time.Unix(TimeValidationMin, 0),
		ExitTime:      &exitTime,
		StudentNumber: sampleStudentNumber,
	}
	err = ev.EntryValidation(entry)
	assert.NoError(t, err)

}
