package repository

import (
	"api/model"
	"time"

	"gorm.io/gorm"
)

type IEntryRepository interface {
	Entry(studentNumber uint, entryTime time.Time) error
	Exit(id uint, exitTime time.Time) error
	FindUnexitedEntry(entry *model.Entry, studentNumber uint) error
}

type EntryRepository struct {
	db *gorm.DB
}

func NewEntryRepository(db *gorm.DB) *EntryRepository {
	return &EntryRepository{db}
}

func (er *EntryRepository) Entry(studentNumber uint, entryTime time.Time) error {
	newEntry := model.Entry{
		EntryTime:     entryTime,
		ExitTime:      nil,
		StudentNumber: studentNumber,
	}

	if err := er.db.Create(&newEntry).Error; err != nil {
		return err
	}

	return nil
}

func (er *EntryRepository) Exit(id uint, exitTime time.Time) error {
	if err := er.db.Model(&model.Entry{}).Where("id=?", id).Update("exit_time", exitTime).Error; err != nil {
		return err
	}
	return nil
}

func (er *EntryRepository) FindUnexitedEntry(entry *model.Entry, studentNumber uint) error {
	if err := er.db.Where("student_number=? AND exit_time IS NULL", studentNumber).FirstOrInit(&entry).Error; err != nil {
		return err
	}
	return nil
}
