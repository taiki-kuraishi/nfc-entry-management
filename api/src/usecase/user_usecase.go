package usecase

import (
	"api/model"
	"api/repository"
	"api/validator"
	"strconv"
)

type IUserUsecase interface {
	CreateOrUpdateUser(user model.User) (string, error)
}

type UserUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &UserUsecase{ur: ur, uv: uv}
}

func (uu *UserUsecase) CreateOrUpdateUser(user model.User) (string, error) {
	var studentNumberAsString string = strconv.Itoa(int(user.StudentNumber))

	//Validate user
	if eer := uu.uv.UserValidation(user); eer != nil {
		return "", eer
	}

	DBUser := model.User{}

	if err := uu.ur.GetUserByStudentNumber(&DBUser, user.StudentNumber); err != nil {
		return "", err
	}

	//Create user
	if DBUser.StudentNumber == 0 {
		if err := uu.ur.CreateUser(&user); err != nil {
			return "", err
		}
		return "Create user : " + studentNumberAsString, nil
	}

	//Update user
	if DBUser.Name != user.Name {
		if err := uu.ur.UpdateUser(&user); err != nil {
			return "", err
		}
		return "User updated : " + studentNumberAsString, nil
	}

	return "User already exists : " + studentNumberAsString, nil
}
