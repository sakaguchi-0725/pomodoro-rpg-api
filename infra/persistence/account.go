package persistence

import (
	"pomodoro-rpg-api/domain/model"
	"pomodoro-rpg-api/domain/repository"
	"pomodoro-rpg-api/infra/entity"
	"pomodoro-rpg-api/pkg/apperr"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type accountPersistence struct {
	db *gorm.DB
}

func (p *accountPersistence) FindByEmail(email string) (model.Account, error) {
	var entity entity.Account
	if err := p.db.Where("email = ?", email).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Account{}, errors.WithStack(apperr.ErrDataNotFound)
		}
		return model.Account{}, errors.WithStack(err)
	}

	return model.RecreateAccount(
		model.AccountID(entity.ID),
		entity.CognitoUID,
		entity.Email,
		entity.Name,
		entity.Image,
	), nil
}

func (p *accountPersistence) Create(acc model.Account) error {
	entity := entity.ToAccountEntity(acc)

	if err := p.db.Create(&entity).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (p *accountPersistence) FindByID(id model.AccountID) (model.Account, error) {
	var acc entity.Account

	err := p.db.Where("id = ?", id).First(&acc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Account{}, errors.WithStack(apperr.ErrDataNotFound)
		}
		return model.Account{}, errors.WithStack(err)
	}

	return model.RecreateAccount(
		model.AccountID(acc.ID),
		acc.CognitoUID,
		acc.Email,
		acc.Name,
		acc.Image,
	), nil
}

func (p *accountPersistence) Update(acc model.Account) error {
	entity := entity.ToAccountEntity(acc)

	if err := p.db.Save(&entity).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func NewaccountPersistence(db *gorm.DB) repository.AccountRepository {
	return &accountPersistence{db}
}
