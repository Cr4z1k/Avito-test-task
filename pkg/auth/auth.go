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
	Parse(token, claim string) (interface{}, error)
	NewToken(isAdmin bool) (string, error)
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

func (m *Manager) Parse(token, claim string) (interface{}, error) {
	tkn, err := jwt.Parse(token, func(tkn *jwt.Token) (i interface{}, err error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign-in method: %v", tkn.Header["alg"])
		}

		return []byte(m.signInKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token claims are not of type tokenClaims")
	}

	value, ok := claims[claim]
	if !ok {
		return nil, fmt.Errorf("claim %s not found", claim)
	}

	return value, nil
}

func (m *Manager) NewToken(isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":     jwt.NewNumericDate(time.Now().Add(ttl)),
		"isAdmin": isAdmin,
	})

	return token.SignedString([]byte(m.signInKey))
}
