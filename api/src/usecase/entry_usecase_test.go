package usecase_test

import (
	"api/model"
	"api/repository"
	"api/usecase"
	"api/validator"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true}}, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestEntryUsecase_EntryOrExit(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	eu := usecase.IEntryUsecase(
		usecase.NewEntryUsecase(
			repository.IEntryRepository(repository.NewEntryRepository(gormDB)),
			validator.IEntryValidator(validator.NewEntryValidator()),
		),
	)

	sampleStudentNumber := uint(20122027)
	sampleTimestamp := time.Now().Add(-time.Second)
	TimeValidationMin := int64(1704034800)

	t.Run("Entry", func(t *testing.T) {
		//GetStudentNumberWithNullExitTime
		initEntry := model.Entry{}
		rows := sqlmock.NewRows([]string{"id", "entry_time", "exit_time", "student_number"}).
			AddRow(initEntry.ID, initEntry.EntryTime, initEntry.ExitTime, initEntry.StudentNumber)
		mock.ExpectQuery("SELECT \\* FROM `entries` WHERE student_number=\\? AND exit_time IS NULL ORDER BY `entries`.`id` LIMIT 1").
			WithArgs(sampleStudentNumber).
			WillReturnRows(rows)

		//CreateEntry
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `entries` \\(`entry_time`,`exit_time`,`student_number`\\) VALUES \\(\\?,\\?,\\?\\)").
			WithArgs(sampleTimestamp, nil, sampleStudentNumber).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result, err := eu.EntryOrExit(sampleStudentNumber, sampleTimestamp)
		assert.NoError(t, err)
		assert.Equal(t, "entry success", result)
	})

	t.Run("Exit", func(t *testing.T) {
		storedEntry := model.Entry{
			ID:            1,
			EntryTime:     time.Unix(TimeValidationMin, 0),
			ExitTime:      nil,
			StudentNumber: sampleStudentNumber,
		}

		rows := sqlmock.NewRows([]string{"id", "entry_time", "exit_time", "student_number"}).
			AddRow(storedEntry.ID, storedEntry.EntryTime, storedEntry.ExitTime, storedEntry.StudentNumber)
		mock.ExpectQuery("SELECT \\* FROM `entries` WHERE student_number=\\? AND exit_time IS NULL ORDER BY `entries`.`id` LIMIT 1").
			WithArgs(sampleStudentNumber).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `entries` SET `exit_time`=\\? WHERE id=\\?").
			WithArgs(sampleTimestamp, storedEntry.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result, err := eu.EntryOrExit(sampleStudentNumber, sampleTimestamp)
		assert.NoError(t, err)
		assert.Equal(t, "exit success", result)
	})
}
