package usecase

import (
	"api/model"
	"api/repository"
	"api/validator"
	"time"
)

type IEntryUsecase interface {
	EntryOrExit(studentNumber uint, timestamp time.Time) (string, error)
}

type EntryUsecase struct {
	er repository.IEntryRepository
	ev validator.IEntryValidator
}

func NewEntryUsecase(er repository.IEntryRepository, ev validator.IEntryValidator) IEntryUsecase {
	return &EntryUsecase{er: er, ev: ev}
}

func (eu *EntryUsecase) EntryOrExit(studentNumber uint, timestamp time.Time) (string, error) {
	if err := eu.ev.StudentNumberValidation(studentNumber); err != nil {
		return "", err
	}

	newEntry := model.Entry{}
	if err := eu.er.GetStudentNumberWithNullExitTime(&newEntry, studentNumber); err != nil {
		return "", err
	}

	//entry
	if newEntry.ID == 0 {
		newEntry = model.Entry{
			EntryTime:     timestamp,
			ExitTime:      nil,
			StudentNumber: studentNumber,
		}

		if err := eu.ev.EntryValidation(newEntry); err != nil {
			return "", err
		}

		if err := eu.er.CreateEntry(&newEntry); err != nil {
			return "", err
		}
		return "entry success", nil
	}

	//exit
	newEntry.ExitTime = &timestamp

	if err := eu.ev.EntryValidation(newEntry); err != nil {
		return "", err
	}

	if err := eu.er.UpdateEntry(&newEntry); err != nil {
		return "", err
	}

	return "exit success", nil
}
