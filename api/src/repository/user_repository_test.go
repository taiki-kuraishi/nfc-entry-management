package repository_test

import (
	"api/model"
	"api/repository"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateUser(t *testing.T) {
	sampleStudentNumber := uint(20122027)
	sampleName := "カイシ　タロウ"
	sampleTime := time.Unix(0, 0)

	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	tr := repository.IUserRepository(repository.NewUserRepository(gormDB))

	sampleUser := model.User{
		StudentNumber: sampleStudentNumber,
		Name:          sampleName,
		CreatedAt:     sampleTime,
		UpdatedAt:     sampleTime,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` \\(`name`,`created_at`,`updated_at`,`student_number`\\) VALUES \\(\\?,\\?,\\?,\\?\\)").
		WithArgs(sampleName, sampleTime, sampleTime, sampleStudentNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err = tr.CreateUser(&sampleUser); err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateUser(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	tr := repository.IUserRepository(repository.NewUserRepository(gormDB))

	sampleUser := model.User{
		StudentNumber: uint(20122027),
		Name:          "カイシ　タロウ",
		CreatedAt:     time.Unix(0, 0),
		UpdatedAt:     time.Unix(0, 0),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `name`=\\?,`created_at`=\\?,`updated_at`=\\? WHERE `student_number` = \\?").
		WithArgs(sampleUser.Name, sampleUser.CreatedAt, AnyTime{}, sampleUser.StudentNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err = tr.UpdateUser(&sampleUser); err != nil {
		t.Errorf(err.Error())
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetUserByStudentNumber(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	tr := repository.IUserRepository(repository.NewUserRepository(gormDB))

	// テストデータの作成
	sampleUser := &model.User{
		Name:          "カイシ　タロウ",
		CreatedAt:     time.Now().Round(time.Second),
		UpdatedAt:     time.Now().Round(time.Second),
		StudentNumber: uint(20122027),
	}

	// モックの期待値の設定
	rows := sqlmock.NewRows([]string{"name", "created_at", "updated_at", "student_number"}).
		AddRow(sampleUser.Name, sampleUser.CreatedAt, sampleUser.UpdatedAt, sampleUser.StudentNumber)
	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE student_number=\\? ORDER BY `users`.`student_number` LIMIT 1$").
		WithArgs(sampleUser.StudentNumber).WillReturnRows(rows)

	responseUser := &model.User{}
	err = tr.GetUserByStudentNumber(responseUser, sampleUser.StudentNumber)

	assert.NoError(t, err)
	assert.Equal(t, sampleUser, responseUser)
}
