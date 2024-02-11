package controller_test

import (
	"api/controller"
	"api/model"
	"api/repository"
	"api/usecase"
	"api/validator"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
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

func TestApiController_RootController(t *testing.T) {

	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	c := controller.IUserAndEntryController(
		controller.NewUserAndEntryController(
			usecase.IUserUsecase(
				usecase.NewUserUsecase(
					repository.IUserRepository(
						repository.NewUserRepository(gormDB),
					),
					validator.IUserValidator(
						validator.NewUserValidator(),
					),
				),
			),
			usecase.IEntryUsecase(
				usecase.NewEntryUsecase(
					repository.IEntryRepository(
						repository.NewEntryRepository(gormDB),
					),
					validator.IEntryValidator(
						validator.NewEntryValidator(),
					),
				),
			),
		),
	)

	t.Run("valid", func(t *testing.T) {

		request := controller.Request{
			StudentNumber: 20122027,
			Name:          "カイシ　タロウ",
			Timestamp:     float64(time.Now().UnixNano()) / 1e9,
		}

		// convert float64 to time.Time
		seconds := int64(request.Timestamp)
		nanoseconds := int64((request.Timestamp - float64(seconds)) * 1e9)
		requestTimestamp := time.Unix(seconds, nanoseconds)

		//GetUserByStudentNumber
		storedUser := model.User{}
		userRows := sqlmock.NewRows([]string{"name", "created_at", "updated_at", "student_number"}).
			AddRow(storedUser.Name, storedUser.CreatedAt, storedUser.UpdatedAt, storedUser.StudentNumber)
		mock.ExpectQuery("^SELECT \\* FROM `users` WHERE student_number=\\? ORDER BY `users`.`student_number` LIMIT 1$").
			WithArgs(request.StudentNumber).WillReturnRows(userRows)

		//CreateUser
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users` \\(`name`,`created_at`,`updated_at`,`student_number`\\) VALUES \\(\\?,\\?,\\?,\\?\\)").
			WithArgs(request.Name, requestTimestamp, requestTimestamp, request.StudentNumber).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		//GetStudentNumberWithNullExitTime
		initEntry := model.Entry{}
		entryRows := sqlmock.NewRows([]string{"id", "entry_time", "exit_time", "student_number"}).
			AddRow(initEntry.ID, initEntry.EntryTime, initEntry.ExitTime, initEntry.StudentNumber)
		mock.ExpectQuery("SELECT \\* FROM `entries` WHERE student_number=\\? AND exit_time IS NULL ORDER BY `entries`.`id` LIMIT 1").
			WithArgs(request.StudentNumber).
			WillReturnRows(entryRows)

		//CreateEntry
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `entries` \\(`entry_time`,`exit_time`,`student_number`\\) VALUES \\(\\?,\\?,\\?\\)").
			WithArgs(requestTimestamp, nil, request.StudentNumber).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		//request to json
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf(err.Error())
		}

		//create echo context
		echoServer := echo.New()
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader(string(requestJson)),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		echoContext := echoServer.NewContext(req, rec)
		echoContext.SetPath("/")

		if err := c.HandleUserAndEntry(echoContext); err != nil {
			t.Errorf(err.Error())
		}

		if rec.Code != http.StatusOK {
			t.Errorf(err.Error())
		}
	})

	t.Run("empty_body", func(t *testing.T) {
		//create echo context
		echoServer := echo.New()
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader(""),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		echoContext := echoServer.NewContext(req, rec)
		echoContext.SetPath("/")

		err = c.HandleUserAndEntry(echoContext)
		if err != nil {
			t.Errorf(err.Error())
		}

		expectedStatus := http.StatusBadRequest
		assert.Equal(t, expectedStatus, rec.Code)

		expectedBody := "\"code=400, message=Request body can't be empty\"\n"
		assert.Equal(t, expectedBody, rec.Body.String())
	})

	t.Run("invalid_json", func(t *testing.T) {
		//create echo context
		echoServer := echo.New()
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader("{invalid json}"),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		echoContext := echoServer.NewContext(req, rec)
		echoContext.SetPath("/")

		err = c.HandleUserAndEntry(echoContext)
		if err != nil {
			t.Errorf(err.Error())
		}

		expectedStatus := http.StatusBadRequest
		assert.Equal(t, expectedStatus, rec.Code)

		expectedBody := "\"code=400, message=Syntax error: offset=2, error=invalid character 'i' looking for beginning of object key string\"\n"
		assert.Equal(t, expectedBody, rec.Body.String())
	})

	t.Run("not_required_studentNumber", func(t *testing.T) {

		request := controller.Request{
			Name:      "カイシ　タロウ",
			Timestamp: float64(time.Now().UnixNano()) / 1e9,
		}

		//request to json
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf(err.Error())
		}

		//create echo context
		echoServer := echo.New()
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader(string(requestJson)),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		echoContext := echoServer.NewContext(req, rec)
		echoContext.SetPath("/")

		err = c.HandleUserAndEntry(echoContext)
		if err != nil {
			t.Errorf(err.Error())
		}

		expectedStatus := http.StatusBadRequest
		assert.Equal(t, expectedStatus, rec.Code)

		expectedBody := "\"student number is required\"\n"
		assert.Equal(t, expectedBody, rec.Body.String())

	})

	t.Run("not_required_name", func(t *testing.T) {

		request := controller.Request{
			StudentNumber: 20122027,
			Timestamp:     float64(time.Now().UnixNano()) / 1e9,
		}

		//request to json
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf(err.Error())
		}

		//create echo context
		echoServer := echo.New()
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader(string(requestJson)),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		echoContext := echoServer.NewContext(req, rec)
		echoContext.SetPath("/")

		err = c.HandleUserAndEntry(echoContext)
		if err != nil {
			t.Errorf(err.Error())
		}

		expectedStatus := http.StatusBadRequest
		assert.Equal(t, expectedStatus, rec.Code)

		expectedBody := "\"name is required\"\n"
		assert.Equal(t, expectedBody, rec.Body.String())
	})

	t.Run("not_required_timestamp", func(t *testing.T) {

		request := controller.Request{
			StudentNumber: 20122027,
			Name:          "カイシ　タロウ",
		}

		//request to json
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf(err.Error())
		}

		//create echo context
		echoServer := echo.New()
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader(string(requestJson)),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		echoContext := echoServer.NewContext(req, rec)
		echoContext.SetPath("/")

		err = c.HandleUserAndEntry(echoContext)
		if err != nil {
			t.Errorf(err.Error())
		}

		expectedStatus := http.StatusBadRequest
		assert.Equal(t, expectedStatus, rec.Code)

		expectedBody := "\"timestamp is required\"\n"
		assert.Equal(t, expectedBody, rec.Body.String())
	})
}
