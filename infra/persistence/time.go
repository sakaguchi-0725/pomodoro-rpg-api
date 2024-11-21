package persistence

import (
	"pomodoro-rpg-api/domain/model"
	"pomodoro-rpg-api/domain/repository"
	"pomodoro-rpg-api/infra/entity"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type timePersistence struct {
	db *gorm.DB
}

func (p *timePersistence) Create(t model.Time) error {
	entity := entity.ToTimeEntity(t)
	if err := p.db.Create(&entity).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (p *timePersistence) GetAll(accID model.AccountID) ([]model.Time, error) {
	var res []model.Time

	err := p.db.Where("account_id = ?", accID).Find(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []model.Time{}, nil
		}
		return []model.Time{}, errors.WithStack(err)
	}

	return res, nil
}

func NewTimePersistence(db *gorm.DB) repository.TimeRepository {
	return &timePersistence{db}
}
