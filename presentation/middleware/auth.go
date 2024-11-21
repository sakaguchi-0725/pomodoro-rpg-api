package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pomodoro-rpg-api/infra/service"
	"pomodoro-rpg-api/pkg/contextkey"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/square/go-jose.v2"
)

type Authenticator struct {
	UserPoolID     string
	ClientID       string
	jwkCache       *jwkCache
	CognitoService service.CognitoService
}

type jwkCache struct {
	jwks      *jose.JSONWebKeySet
	timestamp time.Time
}

func NewAuthenticator(userPoolID, clientID string, cognitoService service.CognitoService) *Authenticator {
	return &Authenticator{
		UserPoolID:     userPoolID,
		ClientID:       clientID,
		jwkCache:       &jwkCache{},
		CognitoService: cognitoService,
	}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := a.getToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: access token not found", http.StatusUnauthorized)
			return
		}

		sub, email, err := a.validateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextkey.UserID, sub)
		ctx = context.WithValue(ctx, contextkey.Email, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) getToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil || cookie.Value == "" {
		return "", errors.New("Unauthorized: access token not found")
	}
	return cookie.Value, nil
}

func (a *Authenticator) validateToken(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return a.lookupKey(token)
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("failed to parse claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("sub claim missing")
	}

	email, err := a.CognitoService.GetEmail(tokenStr)
	if err != nil {
		return "", "", err
	}

	return sub, email, nil
}

func (a *Authenticator) lookupKey(token *jwt.Token) (interface{}, error) {
	jwks, err := a.getJSONWebKeys()
	if err != nil {
		return nil, err
	}

	kid := token.Header["kid"].(string)
	for _, key := range jwks.Keys {
		if key.KeyID == kid {
			return key.Key, nil
		}
	}

	return nil, errors.New("key not found")
}

func (a *Authenticator) getJSONWebKeys() (*jose.JSONWebKeySet, error) {
	if a.jwkCache.jwks != nil && time.Since(a.jwkCache.timestamp).Minutes() < 10.0 {
		return a.jwkCache.jwks, nil
	}

	url := fmt.Sprintf("https://cognito-idp.ap-northeast-1.amazonaws.com/%s/.well-known/jwks.json", a.UserPoolID)

	res, err := http.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	var jwks jose.JSONWebKeySet
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return nil, errors.WithStack(err)
	}

	a.jwkCache.jwks = &jwks
	a.jwkCache.timestamp = time.Now()
	return &jwks, nil
}
