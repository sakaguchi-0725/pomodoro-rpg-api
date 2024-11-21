package repository

import "pomodoro-rpg-api/domain/model"

type TimeRepository interface {
	GetAll(accID model.AccountID) ([]model.Time, error)
	Create(t model.Time) error
}
