package repository

import (
	"api/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	GetUserByStudentNumber(user *model.User, studentNumber uint) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) UpdateUser(user *model.User) error {
	if err := ur.db.Save(user).Error; err != nil {
		return err
	}
	return nil

}

func (ur *UserRepository) GetUserByStudentNumber(user *model.User, studentNumber uint) error {
	if err := ur.db.Where("student_number=?", studentNumber).FirstOrInit(user).Error; err != nil {
		return err
	}
	return nil
}
