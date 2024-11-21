package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"pomodoro-rpg-api/pkg/apperr"
	"pomodoro-rpg-api/usecase/output"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/cockroachdb/errors"
	"gopkg.in/square/go-jose.v2"
)

var (
	unauthorizedException     *types.NotAuthorizedException
	userNotFoundException     *types.UserNotFoundException
	invalidPasswordException  *types.InvalidPasswordException
	invalidParameterException *types.InvalidParameterException
	codeMismatchException     *types.CodeMismatchException
)

type CognitoService interface {
	SignIn(email, password string) (output.SignIn, error)
	SignUp(email, password string) (string, error)
	SignOut(token string) error
	ConfirmSignUp(email, code string) error
	GetEmail(token string) (string, error)
	ChangePassword(token, previousPass, proposedPass string) error
	ForgotPassword(email string) error
	ConfirmForgotPassword(email, code, password string) error
	GetJSONWebKeys() (*jose.JSONWebKeySet, error)
}

type cognitoService struct {
	Client       *cognitoidentityprovider.Client
	ClientID     string
	ClientSecret string
	UserPoolID   string
}

func (c *cognitoService) ConfirmForgotPassword(email string, code string, password string) error {
	input := cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(c.ClientID),
		ConfirmationCode: aws.String(code),
		Password:         aws.String(password),
		Username:         aws.String(email),
		SecretHash:       aws.String(secretHash(email, c.ClientID, c.ClientSecret)),
	}

	_, err := c.Client.ConfirmForgotPassword(context.TODO(), &input)
	if err != nil {
		if errors.Is(err, invalidParameterException) || errors.Is(err, codeMismatchException) {
			return errors.WithStack(apperr.ErrInvalidParameter)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *cognitoService) ForgotPassword(email string) error {
	input := cognitoidentityprovider.ForgotPasswordInput{
		ClientId:   aws.String(c.ClientID),
		Username:   aws.String(email),
		SecretHash: aws.String(secretHash(email, c.ClientID, c.ClientSecret)),
	}

	_, err := c.Client.ForgotPassword(context.TODO(), &input)
	if err != nil {
		if errors.Is(err, invalidParameterException) {
			return errors.WithStack(apperr.ErrInvalidParameter)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *cognitoService) ChangePassword(token string, previousPass string, proposedPass string) error {
	input := cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(token),
		PreviousPassword: aws.String(previousPass),
		ProposedPassword: aws.String(proposedPass),
	}

	_, err := c.Client.ChangePassword(context.TODO(), &input)
	if err != nil {
		if errors.Is(err, unauthorizedException) {
			return errors.WithStack(apperr.ErrUnautorizedExeption)
		}
		if errors.Is(err, invalidParameterException) {
			return errors.WithStack(apperr.ErrInvalidParameter)
		}

		return errors.WithStack(err)
	}

	return nil
}

func (c *cognitoService) GetEmail(token string) (string, error) {
	input := cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	}

	output, err := c.Client.GetUser(context.TODO(), &input)
	if err != nil {
		if errors.Is(err, userNotFoundException) {
			return "", errors.WithStack(apperr.ErrDataNotFound)
		}
		return "", errors.WithStack(err)
	}

	var email string
	for _, v := range output.UserAttributes {
		if *v.Name == "email" {
			email = *v.Value
		}
	}

	if email == "" {
		return "", errors.New("email not found")
	}

	return email, nil
}

func (c *cognitoService) ConfirmSignUp(email string, code string) error {
	hash := secretHash(email, c.ClientID, c.ClientSecret)
	input := cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         &c.ClientID,
		ConfirmationCode: &code,
		Username:         &email,
		SecretHash:       &hash,
	}

	_, err := c.Client.ConfirmSignUp(context.TODO(), &input)
	if err != nil {
		if errors.Is(err, codeMismatchException) || errors.Is(err, invalidParameterException) {
			return errors.WithStack(apperr.ErrInvalidParameter)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *cognitoService) SignIn(email string, password string) (output.SignIn, error) {
	result, err := c.Client.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(c.ClientID),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash(email, c.ClientID, c.ClientSecret),
		},
	})

	if err != nil {
		switch {
		case errors.As(err, unauthorizedException):
			return output.SignIn{}, errors.WithStack(fmt.Errorf("not authorized: %w", err))
		case errors.As(err, userNotFoundException):
			return output.SignIn{}, errors.WithStack(fmt.Errorf("user not found: %w", err))
		default:
			return output.SignIn{}, errors.WithStack(fmt.Errorf("unknown error: %w", err))
		}
	}

	return output.SignIn{
		AccessToken:  *result.AuthenticationResult.AccessToken,
		IdToken:      *result.AuthenticationResult.IdToken,
		RefreshToken: *result.AuthenticationResult.RefreshToken,
	}, nil
}

func (c *cognitoService) SignUp(email string, password string) (string, error) {
	hash := secretHash(email, c.ClientID, c.ClientSecret)
	result, err := c.Client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(c.ClientID),
		Password:   aws.String(password),
		Username:   aws.String(email),
		SecretHash: aws.String(hash),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	})

	if err != nil {
		if errors.As(err, invalidPasswordException) {
			return "", errors.WithStack(fmt.Errorf("invalid password: %w", err))
		}
		return "", errors.WithStack(fmt.Errorf("unknown error: %w", err))
	}

	return *result.UserSub, nil
}

func (c *cognitoService) GetJSONWebKeys() (*jose.JSONWebKeySet, error) {
	url := fmt.Sprintf("https://cognito-idp.ap-northeast-1.amazonaws.com/%s/.well-known/jwks.json", c.UserPoolID)

	res, err := http.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	var jwks jose.JSONWebKeySet
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return nil, errors.WithStack(err)
	}

	return &jwks, nil
}

func (c *cognitoService) SignOut(token string) error {
	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(token),
	}

	_, err := c.Client.GlobalSignOut(context.TODO(), input)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func secretHash(email string, clientID string, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(email + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewCognitoService(clientID, clientSecret, userPoolID string) (CognitoService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &cognitoService{
		Client:       cognitoidentityprovider.NewFromConfig(cfg),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		UserPoolID:   userPoolID,
	}, nil
}
