package model

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type TimeID string

func NewTimeID(s string) (TimeID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", errors.New("invalid account id")
	}

	return TimeID(id.String()), nil
}

func GenerateTimeID() TimeID {
	return TimeID(uuid.NewString())
}

func (a TimeID) String() string {
	return string(a)
}
