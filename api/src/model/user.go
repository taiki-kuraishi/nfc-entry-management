package model

import (
	"time"
)

type User struct {
	StudentNumber uint      `json:"student_number" gorm:"primary_key"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
