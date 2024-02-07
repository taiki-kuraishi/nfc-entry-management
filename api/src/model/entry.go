package model

import "time"

type Entry struct {
	ID            uint       `json:"id" gorm:"primary_key"`
	EntryTime     time.Time  `json:"entry_time" gorm:"not null"`
	ExitTime      *time.Time `json:"exit_time"`
	StudentNumber uint       `json:"student_number" gorm:"not null;foreignKey"`
}

type EntryRequest struct {
	StudentNumber uint    `json:"student_number"`
	Name          string  `json:"name"`
	Timestamp     float64 `json:"timestamp"`
}
