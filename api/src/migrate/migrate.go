package main

import (
	"api/db"
	"api/model"
	"fmt"
)

func main() {
	dbConn := db.ConnectDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	err := dbConn.AutoMigrate(&model.User{}, &model.Entry{})
	if err != nil {
		fmt.Printf("Failed to auto-migrate: %v", err)
	}
}
