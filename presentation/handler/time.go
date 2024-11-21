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

	"github.com/cockroachdb/errors"
)

type TimeHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type timeHandler struct {
	tu usecase.TimeUsecase
}

func (t *timeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.TimeRequest
	ctx := r.Context()
	email := ctx.Value(contextkey.Email).(string)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Event(ctx, logger.INFO, "invalid request", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "不正なリクエストです", err))
		return
	}

	if err := t.tu.Create(ctx, email, req.FocusTime); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, nil)
}

func NewTimeHandler(tu usecase.TimeUsecase) TimeHandler {
	return &timeHandler{tu}
}
