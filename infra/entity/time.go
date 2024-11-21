package entity

import (
	"pomodoro-rpg-api/domain/model"
	"time"
)

type Time struct {
	ID            string    `gorm:"primaryKey"`
	FocusTime     float64   `gorm:"not null"`
	ExecutionDate time.Time `gorm:"not null"`
	AccountID     string
	Account       Account   `gorm:"foreignKey:AccountID"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func ToTimeEntity(t model.Time) Time {
	return Time{
		ID:            t.ID.String(),
		FocusTime:     t.FocusTime,
		ExecutionDate: t.ExecutionDate,
		AccountID:     t.AccountID.String(),
	}
}
