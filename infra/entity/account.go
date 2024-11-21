package entity

import (
	"pomodoro-rpg-api/domain/model"
	"time"
)

type Account struct {
	ID         string `gorm:"primaryKey"`
	CognitoUID string
	Email      string
	Name       string
	Image      string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func ToAccountEntity(acc model.Account) Account {
	return Account{
		ID:         acc.ID.String(),
		CognitoUID: acc.CognitoUID,
		Email:      acc.Email,
		Name:       acc.Name,
		Image:      acc.Image,
	}
}
