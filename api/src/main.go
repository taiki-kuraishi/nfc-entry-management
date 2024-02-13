package main

import (
	"api/controller"
	"api/db"
	"api/repository"
	"api/router"
	"api/usecase"
	"api/validator"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
)

func main() {
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	mysqlConfig := mysql.Config{
		DriverName:                "mysql",
		DSN:                       url,
		SkipInitializeWithVersion: true,
	}

	db := db.ConnectDB(mysqlConfig)
	entryValidator := validator.NewEntryValidator()
	userValidator := validator.NewUserValidator()
	userRepository := repository.NewUserRepository(db)
	entryRepository := repository.NewEntryRepository(db)
	entryUsecase := usecase.NewEntryUsecase(entryRepository, entryValidator)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	UserAndEntryController := controller.NewUserAndEntryController(userUsecase, entryUsecase)
	e := router.NewRouter(UserAndEntryController)
	e.Logger.Fatal(e.Start(":8080"))
}
