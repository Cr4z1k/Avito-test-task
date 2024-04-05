package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ttl = time.Hour * 24
)

type TokenManager interface {
	Parse(token string) (string, error)
	NewUserToken() (string, error)
	NewAdmToken() (string, error)
}

type Manager struct {
	signInKey string
}

func NewTokenManager(signInKey string) (*Manager, error) {
	if signInKey == "" {
		return nil, errors.New("empty sign-in key")
	}

	return &Manager{signInKey: signInKey}, nil
}

func (m *Manager) Parse(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(tkn *jwt.Token) (i interface{}, err error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign-in method: %v", tkn.Header["alg"])
		}

		return []byte(m.signInKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token claims are not of type tokenClaims")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewUserToken() (string, error) {
	token := 
}
