package repository

import "pomodoro-rpg-api/domain/model"

type AccountRepository interface {
	FindByID(id model.AccountID) (model.Account, error)
	FindByEmail(email string) (model.Account, error)
	Create(acc model.Account) error
	Update(acc model.Account) error
}
