package db_test

import (
	"api/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(
		mysql.Dialector{
			Config: &mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			},
		},
		&gorm.Config{},
	)
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestConnectDB(t *testing.T) {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf(err.Error())
	}

	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).
			AddRow("8.0.36"))

	mysqlConfig := mysql.Config{
		DriverName: "mysql",
		Conn:       mockdb,
	}

	gormDB := db.ConnectDB(mysqlConfig)

	if gormDB == nil {
        t.Errorf("failed to connect to database")
    }

	if _, err := gormDB.DB(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestCloseDB(t *testing.T) {
	gormDB, mock, err := NewDBMock()
	if err != nil {
		t.Errorf(err.Error())
	}

	mock.ExpectClose()
	db.CloseDB(gormDB)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(err.Error())
	}
}
