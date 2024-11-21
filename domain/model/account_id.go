package model

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type AccountID string

func NewAccountID(s string) (AccountID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", errors.New("invalid account id")
	}

	return AccountID(id.String()), nil
}

func GenerateAccountID() AccountID {
	return AccountID(uuid.NewString())
}

func (a AccountID) String() string {
	return string(a)
}
