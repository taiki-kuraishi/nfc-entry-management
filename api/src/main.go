package main

import (
	"api/controller"
	"api/db"
	"api/repository"
	"api/router"
	"api/usecase"
	"api/validator"
)

func main() {
	db := db.ConnectDB()
	userValidator := validator.NewUserValidator()
	userRepository := repository.NewUserRepository(db)
	entryRepository := repository.NewEntryRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	entryUsecase := usecase.NewEntryUsecase(entryRepository)
	apiController := controller.NewApiController(userUsecase, entryUsecase)
	e := router.NewRouter(apiController)
	e.Logger.Fatal(e.Start(":8080"))
}
