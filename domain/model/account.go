package model

import "github.com/cockroachdb/errors"

type Account struct {
	ID         AccountID
	CognitoUID string
	Email      string
	Name       string
	Image      string
}

func NewAccount(id AccountID, cognitoUID, email, name, image string) (Account, error) {
	if cognitoUID == "" {
		return Account{}, errors.New("cognitoUID is required")
	}

	if email == "" {
		return Account{}, errors.New("email is required")
	}

	if name == "" {
		return Account{}, errors.New("name is required")
	}

	return Account{
		ID:         id,
		CognitoUID: cognitoUID,
		Email:      email,
		Name:       name,
		Image:      image,
	}, nil
}

func RecreateAccount(id AccountID, cognitoUID, email, name, image string) Account {
	return Account{
		ID:         id,
		CognitoUID: cognitoUID,
		Email:      email,
		Name:       name,
		Image:      image,
	}
}

func (a *Account) UpdateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	if a.Name != name {
		a.Name = name
	}

	return nil
}

func (a *Account) UpdateImage(img string) {
	if a.Image != img {
		a.Image = img
	}
}
