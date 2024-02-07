package usecase

import (
	"api/model"
	"api/repository"
)

type IUserUsecase interface {
	CreateOrUpdateUser(user model.User) (string, error)
}

type UserUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &UserUsecase{ur: ur}
}

func (uu *UserUsecase) CreateOrUpdateUser(user model.User) (string, error) {
	DBUser := model.User{}

	if err := uu.ur.GetUserByStudentNumber(&DBUser, user.StudentNumber); err != nil {
		return "", err
	}

	//Create user
	if DBUser.StudentNumber == 0 {
		if err := uu.ur.CreateUser(&user); err != nil {
			return "", err
		}
		return "User created", nil
	}

	//Update user
	if DBUser.Name != user.Name {
		if err := uu.ur.UpdateUser(&user); err != nil {
			return "", err
		}
		return "User updated", nil
	}

	return "User already exists", nil
}
