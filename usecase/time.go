package usecase

import (
	"context"
	"pomodoro-rpg-api/domain/model"
	"pomodoro-rpg-api/domain/repository"
	"pomodoro-rpg-api/pkg/apperr"
	"pomodoro-rpg-api/pkg/logger"

	"github.com/cockroachdb/errors"
)

type TimeUsecase interface {
	Create(ctx context.Context, email string, focusTime float64) error
}

type timeUsecase struct {
	ar repository.AccountRepository
	tr repository.TimeRepository
}

func (t *timeUsecase) Create(ctx context.Context, email string, focusTime float64) error {
	acc, err := t.ar.FindByEmail(email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			logger.Event(ctx, logger.INFO, "account not found", err)
			return apperr.NewApplicationError(apperr.ErrBadRequest, "アカウント情報が取得できませんでした", err)
		}
		logger.Event(ctx, logger.ERROR, "find account failed", err)
		return err
	}

	id := model.GenerateTimeID()
	model, err := model.NewTime(id, focusTime, acc.ID)
	if err != nil {
		logger.Event(ctx, logger.INFO, err.Error(), err)
		return apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err)
	}

	if err := t.tr.Create(model); err != nil {
		logger.Event(ctx, logger.ERROR, "create failed", err)
		return err
	}

	return nil
}

func NewTimeUsecase(ar repository.AccountRepository, tr repository.TimeRepository) TimeUsecase {
	return &timeUsecase{ar, tr}
}
