package handler

import (
	"encoding/json"
	"net/http"
	"pomodoro-rpg-api/pkg/apperr"
	"pomodoro-rpg-api/pkg/formatter"
	"pomodoro-rpg-api/pkg/logger"
	"pomodoro-rpg-api/presentation/dto"
	"pomodoro-rpg-api/presentation/response"
	"pomodoro-rpg-api/usecase"
	"pomodoro-rpg-api/usecase/input"
	"pomodoro-rpg-api/usecase/output"
	"reflect"
	"time"

	"github.com/cockroachdb/errors"
)

type AuthHandler interface {
	IsAuth(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	ConfirmSignUp(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	ForgotPassword(w http.ResponseWriter, r *http.Request)
	ConfirmForgotPassword(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	au usecase.AuthUsecase
}

func (a *authHandler) IsAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("access_token")
	if err != nil || cookie.Value == "" {
		response.JSON(w, http.StatusOK, map[string]bool{"isAuthenticated": false})
		return
	}

	isAuth, err := a.au.VerifyToken(ctx, cookie.Value)
	if err != nil {
		logger.Event(ctx, logger.ERROR, "verifyToken failed", err)
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"isAuthenticated": isAuth})
}

func (a *authHandler) ConfirmForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmForgotPasswordRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Event(ctx, logger.INFO, "invalid input", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err))
		return
	}

	if err := a.au.ConfirmForgotPassword(ctx, req.Email, req.Code, req.Password); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (a *authHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Event(ctx, logger.INFO, "invalid input", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err))
		return
	}

	if err := a.au.ForgotPassword(ctx, req.Email); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (a *authHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ChangePasswordRequest
	ctx := r.Context()

	cookie, err := r.Cookie("access_token")
	if err != nil {
		logger.Event(ctx, logger.ERROR, "access token not found", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrUnautorized, "アクセストークンが不正です", err))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Event(ctx, logger.INFO, "decode error", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err))
		return
	}

	err = a.au.ChangePassword(ctx, cookie.Value, req.PreviousPass, req.ProposedPass)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (a *authHandler) ConfirmSignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmSignUpRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Event(ctx, logger.INFO, err.Error(), errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err))
		return
	}

	if err := a.au.ConfirmSignUp(ctx, req.Email, req.Code); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (a *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Event(ctx, logger.INFO, "request decode failed", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err))
		return
	}

	output, err := a.au.SignIn(ctx, req.Email, req.Passowrd)
	if err != nil {
		response.Error(w, err)
		return
	}

	setTokenCookies(w, output)
	response.JSON(w, http.StatusOK, nil)
}

func (a *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Event(ctx, logger.INFO, "request decode failed", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrBadRequest, "入力値が正しくありません", err))
		return
	}

	input := input.SignUp{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := a.au.SignUp(ctx, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (a *authHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("access_token")
	if err != nil {
		logger.Event(ctx, logger.INFO, "cookie not found", errors.WithStack(err))
		response.Error(w, apperr.NewApplicationError(apperr.ErrUnautorized, "access tokenが見つかりません", err))
		return
	}

	if err := a.au.SignOut(ctx, cookie.Value); err != nil {
		response.Error(w, err)
		return
	}

	deleteAllCookies(w, r)
	response.JSON(w, http.StatusOK, "logout success")
}

func setTokenCookies(w http.ResponseWriter, tokens output.SignIn) {
	v := reflect.ValueOf(tokens)
	t := reflect.TypeOf(tokens)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).String()

		cookie := &http.Cookie{
			Name:     formatter.ToSnakeCase(fieldName),
			Value:    fieldValue,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		}

		http.SetCookie(w, cookie)
	}
}

func deleteAllCookies(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		http.SetCookie(w, &http.Cookie{
			Name:     cookie.Name,
			Value:    "",
			Path:     cookie.Path,
			Expires:  time.Now(),
			MaxAge:   -1,
			HttpOnly: cookie.HttpOnly,
			// Secure:   cookie.Secure,
			SameSite: http.SameSiteNoneMode,
		})
	}
}

func NewAuthHandler(au usecase.AuthUsecase) AuthHandler {
	return &authHandler{au}
}
