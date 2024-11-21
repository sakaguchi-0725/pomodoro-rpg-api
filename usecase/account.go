package usecase

import (
	"context"
	"pomodoro-rpg-api/domain/repository"
	"pomodoro-rpg-api/pkg/apperr"
	"pomodoro-rpg-api/pkg/logger"
	"pomodoro-rpg-api/usecase/input"
	"pomodoro-rpg-api/usecase/output"

	"github.com/cockroachdb/errors"
)

type AccountUsecase interface {
	GetByEmail(ctx context.Context, email string) (output.Account, error)
	Update(ctx context.Context, input input.Account) error
}

type accountUsecase struct {
	ar repository.AccountRepository
}

func (a *accountUsecase) GetByEmail(ctx context.Context, email string) (output.Account, error) {
	acc, err := a.ar.FindByEmail(email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			logger.Event(ctx, logger.INFO, "account not found", err)
			return output.Account{}, apperr.NewApplicationError(apperr.ErrBadRequest, "アカウント情報を取得できませんでした", err)
		}
		logger.Event(ctx, logger.ERROR, "account find failed", err)
		return output.Account{}, err
	}

	output := output.Account{
		Email: acc.Email,
		Name:  acc.Name,
		Image: acc.Image,
	}

	return output, nil
}

func (a *accountUsecase) Update(ctx context.Context, input input.Account) error {
	acc, err := a.ar.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			logger.Event(ctx, logger.INFO, "account not found", err)
			return apperr.NewApplicationError(apperr.ErrBadRequest, "アカウント情報を取得できませんでした", err)
		}
		logger.Event(ctx, logger.ERROR, "account find failed", err)
		return err
	}

	acc.UpdateImage(input.Image)
	if err := acc.UpdateName(input.Name); err != nil {
		logger.Event(ctx, logger.INFO, "update name failed", err)
		return apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が不正です", err)
	}

	if err := a.ar.Update(acc); err != nil {
		logger.Event(ctx, logger.ERROR, "account update failed", err)
		return err
	}

	return nil
}

func NewAccountUsecase(ar repository.AccountRepository) AccountUsecase {
	return &accountUsecase{ar}
}
