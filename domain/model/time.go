package model

import (
	"time"

	"github.com/pkg/errors"
)

type Time struct {
	ID            TimeID
	FocusTime     float64
	AccountID     AccountID
	ExecutionDate time.Time
}

func NewTime(id TimeID, focusTime float64, accID AccountID) (Time, error) {
	if focusTime <= 0 {
		return Time{}, errors.New("focus time is 0 or more")
	}

	return Time{
		ID:            id,
		FocusTime:     focusTime,
		AccountID:     accID,
		ExecutionDate: time.Now(),
	}, nil
}

func RecreateTime(id TimeID, focusTime float64, accID AccountID, dateStr string) Time {
	layout := "2006-01-02T15:04"
	executionDate, _ := time.Parse(layout, dateStr)

	return Time{
		ID:            id,
		FocusTime:     focusTime,
		AccountID:     accID,
		ExecutionDate: executionDate,
	}
}
