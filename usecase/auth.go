package usecase

import (
	"context"
	"pomodoro-rpg-api/domain/model"
	"pomodoro-rpg-api/domain/repository"
	"pomodoro-rpg-api/infra/service"
	"pomodoro-rpg-api/pkg/apperr"
	"pomodoro-rpg-api/pkg/logger"
	"pomodoro-rpg-api/usecase/input"
	"pomodoro-rpg-api/usecase/output"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase interface {
	SignIn(ctx context.Context, email, password string) (output.SignIn, error)
	SignUp(ctx context.Context, input input.SignUp) error
	SignOut(ctx context.Context, token string) error
	ConfirmSignUp(ctx context.Context, email, code string) error
	ChangePassword(ctx context.Context, token string, previousPass string, proposedPass string) error
	ForgotPassword(ctx context.Context, email string) error
	ConfirmForgotPassword(ctx context.Context, email, code, password string) error
	VerifyToken(ctx context.Context, tokenStr string) (bool, error)
}

type authUsecase struct {
	cs service.CognitoService
	ar repository.AccountRepository
}

func (a *authUsecase) ConfirmForgotPassword(ctx context.Context, email, code, password string) error {
	if err := a.cs.ConfirmForgotPassword(email, code, password); err != nil {
		if errors.Is(err, apperr.ErrInvalidParameter) {
			logger.Event(ctx, logger.INFO, "invalid input", err)
			return apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err)
		}
		logger.Event(ctx, logger.ERROR, "confirm forgot password failed", err)
		return err
	}

	return nil
}

func (a *authUsecase) ForgotPassword(ctx context.Context, email string) error {
	if err := a.cs.ForgotPassword(email); err != nil {
		if errors.Is(err, apperr.ErrInvalidParameter) {
			logger.Event(ctx, logger.INFO, "invalid email", err)
			return apperr.NewApplicationError(apperr.ErrBadRequest, "ユーザーが見つかりませんでした", err)
		}
		logger.Event(ctx, logger.ERROR, "forgot password failed", err)
		return err
	}

	return nil
}

func (a *authUsecase) ChangePassword(ctx context.Context, token string, previousPass string, proposedPass string) error {
	if err := a.cs.ChangePassword(token, previousPass, proposedPass); err != nil {
		if errors.Is(err, apperr.ErrUnautorizedExeption) {
			logger.Event(ctx, logger.INFO, "unauthorized", err)
			return apperr.NewApplicationError(apperr.ErrUnautorized, "アクセストークンが不正です", err)
		}
		if errors.Is(err, apperr.ErrInvalidParameter) {
			logger.Event(ctx, logger.INFO, "invalid parameter", err)
			return apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が不正です", err)
		}
		logger.Event(ctx, logger.ERROR, "change password failed", err)
		return err
	}

	return nil
}

func (a *authUsecase) ConfirmSignUp(ctx context.Context, email string, code string) error {
	if err := a.cs.ConfirmSignUp(email, code); err != nil {
		if errors.Is(err, apperr.ErrInvalidParameter) {
			logger.Event(ctx, logger.INFO, err.Error(), err)
			return apperr.NewApplicationError(apperr.ErrBadRequest, "入力が間違っています", err)
		}
		logger.Event(ctx, logger.ERROR, "confirm signup failed", err)
		return err
	}

	return nil
}

func (a *authUsecase) SignIn(ctx context.Context, email string, password string) (output.SignIn, error) {
	res, err := a.cs.SignIn(email, password)
	if err != nil {
		logger.Event(ctx, logger.INFO, "signin failed", err)
		return output.SignIn{}, apperr.NewApplicationError(apperr.ErrUnautorized, "signin failed", err)
	}

	return res, nil
}

func (a *authUsecase) SignUp(ctx context.Context, input input.SignUp) error {
	cognitoUID, err := a.cs.SignUp(input.Email, input.Password)
	if err != nil {
		logger.Event(ctx, logger.INFO, "signup failed", err)
		return apperr.NewApplicationError(apperr.ErrBadRequest, "sign up failed", err)
	}

	accID := model.GenerateAccountID()
	acc, err := model.NewAccount(accID, cognitoUID, input.Email, input.Name, "")
	if err != nil {
		logger.Event(ctx, logger.INFO, "invalid input", err)
		return apperr.NewApplicationError(apperr.ErrBadRequest, "invalid input", err)
	}

	if err := a.ar.Create(acc); err != nil {
		logger.Event(ctx, logger.INFO, "account create failed", err)
		return apperr.NewApplicationError(apperr.ErrBadRequest, "sign up failed", err)
	}

	return nil
}

func (a *authUsecase) SignOut(ctx context.Context, token string) error {
	if err := a.cs.SignOut(token); err != nil {
		logger.Event(ctx, logger.ERROR, "SignOut failed", err)
		return err
	}

	return nil
}

func (a *authUsecase) VerifyToken(ctx context.Context, tokenStr string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return a.lookupKey(token)
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, nil
	}

	if time.Now().Unix() > int64(exp) {
		return false, nil
	}

	return true, nil
}

func (a *authUsecase) lookupKey(token *jwt.Token) (interface{}, error) {
	jwks, err := a.cs.GetJSONWebKeys()
	if err != nil {
		return nil, err
	}

	kid := token.Header["kid"].(string)
	for _, key := range jwks.Keys {
		if key.KeyID == kid {
			return key.Key, nil
		}
	}

	return nil, errors.New("key not found")
}

func NewAuthUsecase(cs service.CognitoService, ar repository.AccountRepository) AuthUsecase {
	return &authUsecase{cs, ar}
}
