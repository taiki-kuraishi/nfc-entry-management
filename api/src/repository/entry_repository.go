package repository

import (
	"api/model"

	"gorm.io/gorm"
)

type IEntryRepository interface {
	CreateEntry(entry *model.Entry) error
	UpdateEntry(entry *model.Entry) error
	GetStudentNumberWithNullExitTime(entry *model.Entry, studentNumber uint) error
}

type EntryRepository struct {
	db *gorm.DB
}

func NewEntryRepository(db *gorm.DB) *EntryRepository {
	return &EntryRepository{db}
}

func (er *EntryRepository) CreateEntry(entry *model.Entry) error {
	if err := er.db.Create(&entry).Error; err != nil {
		return err
	}

	return nil
}

func (er *EntryRepository) UpdateEntry(entry *model.Entry) error {
	if err := er.db.Model(&model.Entry{}).Where("id=?", entry.ID).Update("exit_time", entry.ExitTime).Error; err != nil {
		return err
	}
	return nil
}

func (er *EntryRepository) GetStudentNumberWithNullExitTime(entry *model.Entry, studentNumber uint) error {
	if err := er.db.Where("student_number=? AND exit_time IS NULL", studentNumber).FirstOrInit(&entry).Error; err != nil {
		return err
	}
	return nil
}
