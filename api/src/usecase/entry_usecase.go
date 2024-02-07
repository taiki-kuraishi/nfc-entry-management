package usecase

import (
	"api/model"
	"api/repository"
	"time"
)

type IEntryUsecase interface {
	EntryOrExit(studentNumber uint, timestamp time.Time) (string, error)
}

type EntryUsecase struct {
	er repository.IEntryRepository
}

func NewEntryUsecase(er repository.IEntryRepository) IEntryUsecase {
	return &EntryUsecase{er: er}
}

func (eu *EntryUsecase) EntryOrExit(studentNumber uint, timestamp time.Time) (string, error) {
	newEntry := model.Entry{}
	if err := eu.er.FindUnexitedEntry(&newEntry, studentNumber); err != nil {
		return "", err
	}

	//entry
	if newEntry.ID == 0 {
		if err := eu.er.Entry(studentNumber, timestamp); err != nil {
			return "", err
		}
		return "entry success", nil
	}

	//exit
	if err := eu.er.Exit(newEntry.ID, timestamp); err != nil {
		return "", err
	}

	return "exit success", nil
}
