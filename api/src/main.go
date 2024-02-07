package main

import (
	"api/controller"
	"api/db"
	"api/repository"
	"api/router"
	"api/usecase"
	"fmt"
	"os"
	"time"
)

func main() {
	location, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		fmt.Println(err.Error())
	}

	db := db.ConnectDB()
	userRepository := repository.NewUserRepository(db)
	entryRepository := repository.NewEntryRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	entryUsecase := usecase.NewEntryUsecase(entryRepository)
	apiController := controller.NewApiController(userUsecase, entryUsecase, location)
	e := router.NewRouter(apiController)
	e.Logger.Fatal(e.Start(":8080"))
}
