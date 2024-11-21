package config

import "os"

type AWS struct {
	UserPoolID   string
	ClientID     string
	ClientSecret string
}

func newAWSConfig() *AWS {
	return &AWS{
		UserPoolID:   os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:     os.Getenv("COGNITO_CLIENT_ID"),
		ClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
	}
}
