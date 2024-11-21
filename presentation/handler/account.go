package handler

import (
	"encoding/json"
	"net/http"
	"pomodoro-rpg-api/pkg/apperr"
	"pomodoro-rpg-api/pkg/contextkey"
	"pomodoro-rpg-api/pkg/logger"
	"pomodoro-rpg-api/presentation/dto"
	"pomodoro-rpg-api/presentation/response"
	"pomodoro-rpg-api/usecase"
	"pomodoro-rpg-api/usecase/input"

	"github.com/cockroachdb/errors"
)

type AccountHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type accountHandler struct {
	au usecase.AccountUsecase
}

func (a *accountHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := ctx.Value(contextkey.Email).(string)

	acc, err := a.au.GetByEmail(ctx, email)
	if err != nil {
		response.Error(w, err)
		return
	}

	res := dto.AccountResponse{
		Email: acc.Email,
		Name:  acc.Name,
		Image: acc.Image,
	}

	response.JSON(w, http.StatusOK, res)
}

func (a *accountHandler) Update(w http.ResponseWriter, r *http.Request) {
	var acc dto.UpdateAccountRequest
	ctx := r.Context()
	email := ctx.Value(contextkey.Email).(string)

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		logger.Event(ctx, logger.INFO, "request decode failed", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "invalid request", err))
		return
	}

	input := input.Account{
		Email: email,
		Name:  acc.Name,
		Image: acc.Image,
	}

	if err := a.au.Update(ctx, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "")
}

func NewAccountHandler(au usecase.AccountUsecase) AccountHandler {
	return &accountHandler{au}
}
