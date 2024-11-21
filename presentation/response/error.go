package response

import (
	"net/http"
	"pomodoro-rpg-api/pkg/apperr"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperr.ApplicationError); ok {
		switch appErr.Code() {
		case apperr.ErrBadRequest:
			JSON(w, http.StatusBadRequest, errorResponse{
				Code:    appErr.Code().String(),
				Message: appErr.Message(),
			})
			return
		case apperr.ErrNotFound:
			JSON(w, http.StatusNotFound, errorResponse{
				Code:    appErr.Code().String(),
				Message: appErr.Message(),
			})
			return
		default:
			JSON(w, http.StatusInternalServerError, errorResponse{
				Code: "InternalServerError",
			})
			return
		}
	}

	JSON(w, http.StatusInternalServerError, errorResponse{
		Code: "InternalServerError",
	})
}
