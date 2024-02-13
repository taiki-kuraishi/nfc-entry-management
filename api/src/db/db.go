package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(config mysql.Config) *gorm.DB {
	db, err := gorm.Open(mysql.New(config), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CloseDB(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}

	err = dbSQL.Close()
	if err != nil {
		log.Println(err)
	}
}
