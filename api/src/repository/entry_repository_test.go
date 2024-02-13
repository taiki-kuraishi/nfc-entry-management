package repository_test

import (
	"api/model"
	"api/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestEntryRepository_CreateEntry(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	ter := repository.IEntryRepository(repository.NewEntryRepository(gormDB))

	sampleTime := time.Unix(0, 0)
	entry := model.Entry{
		EntryTime:     sampleTime,
		ExitTime:      &sampleTime,
		StudentNumber: uint(20122027),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `entries` \\(`entry_time`,`exit_time`,`student_number`\\) VALUES \\(\\?,\\?,\\?\\)").
		WithArgs(entry.EntryTime, entry.ExitTime, entry.StudentNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err = ter.CreateEntry(&entry); err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestEntryRepository_UpdateEntry(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	ter := repository.IEntryRepository(repository.NewEntryRepository(gormDB))

	sampleTime := time.Unix(0, 0)
	entry := model.Entry{
		ID:            uint(1),
		EntryTime:     sampleTime,
		ExitTime:      &sampleTime,
		StudentNumber: uint(20122027),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `entries` SET `exit_time`=\\? WHERE id=\\?").
		WithArgs(entry.EntryTime, entry.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err = ter.UpdateEntry(&entry); err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestEntryRepository_GetStudentNumberWithNullExitTime(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Error(err)
	}

	ter := repository.IEntryRepository(repository.NewEntryRepository(gormDB))

	sampleTime := time.Unix(0, 0)
	entry := model.Entry{
		EntryTime:     sampleTime,
		ExitTime:      &sampleTime,
		StudentNumber: uint(20122027),
	}

	rows := sqlmock.NewRows([]string{"id", "entry_time", "exit_time", "student_number"}).
		AddRow(1, sampleTime, nil, entry.StudentNumber)

	mock.ExpectQuery("SELECT \\* FROM `entries` WHERE student_number=\\? AND exit_time IS NULL ORDER BY `entries`.`id` LIMIT 1").
		WithArgs(entry.StudentNumber).
		WillReturnRows(rows)

	if err = ter.GetStudentNumberWithNullExitTime(&entry, entry.StudentNumber); err != nil {
		t.Error(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
