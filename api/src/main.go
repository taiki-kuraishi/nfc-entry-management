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
	userUsecase := usecase.NewUserUsecase(userRepository)
	apiController := controller.NewApiController(userUsecase, location)
	e := router.NewRouter(apiController)
	e.Logger.Fatal(e.Start(":8080"))
}
