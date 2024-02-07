package usecase

import (
	"api/model"
	"api/repository"
	"errors"

	"gorm.io/gorm"
)

type IUserUsecase interface {
	CreateOrUpdateUser(user model.User) error
}

type UserUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &UserUsecase{ur: ur}
}

func (uu *UserUsecase) CreateOrUpdateUser(user model.User) error {
	DBUser := model.User{}

	if err := uu.ur.GetUserByStudentNumber(&DBUser, user.StudentNumber); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := uu.ur.CreateUser(&user); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	//Update user
	if DBUser.Name != user.Name {
		if err := uu.ur.UpdateUser(&user); err != nil {
			return err
		}
	}

	return nil
}
