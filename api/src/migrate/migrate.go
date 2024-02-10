package main

import (
	"api/db"
	"api/model"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
)

func main() {
	mysqlConfig := mysql.Config{
		DriverName: "mysql",
		DSN:        fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE")),
	}

	dbConn := db.ConnectDB(mysqlConfig)
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	err := dbConn.AutoMigrate(&model.User{}, &model.Entry{})
	if err != nil {
		fmt.Printf("Failed to auto-migrate: %v", err)
	}
}
