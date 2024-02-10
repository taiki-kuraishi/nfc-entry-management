package usecase_test

import (
	"api/model"
	"api/repository"
	"api/usecase"
	"api/validator"
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

func TestUserUsecase_CreateOrUpdateUser(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	uu := usecase.IUserUsecase(
		usecase.NewUserUsecase(
			repository.IUserRepository(repository.NewUserRepository(gormDB)),
			validator.IUserValidator(validator.NewUserValidator()),
		),
	)

	t.Run("CreateUser", func(t *testing.T) {
		user := model.User{
			StudentNumber: 20122027,
			Name:          "カイシ　タロウ",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		//GetUserByStudentNumber
		storedUser := model.User{}
		rows := sqlmock.NewRows([]string{"name", "created_at", "updated_at", "student_number"}).
			AddRow(storedUser.Name, storedUser.CreatedAt, storedUser.UpdatedAt, storedUser.StudentNumber)
		mock.ExpectQuery("^SELECT \\* FROM `users` WHERE student_number=\\? ORDER BY `users`.`student_number` LIMIT 1$").
			WithArgs(user.StudentNumber).WillReturnRows(rows)

		//CreateUser
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users` \\(`name`,`created_at`,`updated_at`,`student_number`\\) VALUES \\(\\?,\\?,\\?,\\?\\)").
			WithArgs(user.Name, user.CreatedAt, user.UpdatedAt, user.StudentNumber).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result, err := uu.CreateOrUpdateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, "User created", result)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user := model.User{
			StudentNumber: 20122027,
			Name:          "カイシ　タロウ",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		//GetUserByStudentNumber
		storedUser := model.User{
			StudentNumber: 20122027,
			Name:          "ヨネヤマ　タロウ",
			CreatedAt:     time.Now().Add(-time.Second),
			UpdatedAt:     time.Now().Add(-time.Second),
		}
		rows := sqlmock.NewRows([]string{"name", "created_at", "updated_at", "student_number"}).
			AddRow(storedUser.Name, storedUser.CreatedAt, storedUser.UpdatedAt, storedUser.StudentNumber)
		mock.ExpectQuery("^SELECT \\* FROM `users` WHERE student_number=\\? ORDER BY `users`.`student_number` LIMIT 1$").
			WithArgs(user.StudentNumber).WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users` SET `name`=\\?,`created_at`=\\?,`updated_at`=\\? WHERE `student_number` = \\?").
			WithArgs(user.Name, user.CreatedAt, AnyTime{}, user.StudentNumber).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result, err := uu.CreateOrUpdateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, "User updated", result)
	})

	t.Run("AlreadyExistsUser", func(t *testing.T) {
		user := model.User{
			StudentNumber: 20122027,
			Name:          "カイシ　タロウ",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		//GetUserByStudentNumber
		rows := sqlmock.NewRows([]string{"name", "created_at", "updated_at", "student_number"}).
			AddRow(user.Name, user.CreatedAt, user.UpdatedAt, user.StudentNumber)
		mock.ExpectQuery("^SELECT \\* FROM `users` WHERE student_number=\\? ORDER BY `users`.`student_number` LIMIT 1$").
			WithArgs(user.StudentNumber).WillReturnRows(rows)

		result, err := uu.CreateOrUpdateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, "User already exists", result)
	})
}
