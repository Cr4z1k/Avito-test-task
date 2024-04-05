package service

import (
	"github.com/Cr4z1k/Avito-test-task/pkg/auth"
)

type AuthService struct {
	tokenManager auth.TokenManager
}

func NewAuthService(t auth.TokenManager) *AuthService {
	return &AuthService{tokenManager: t}
}

func (s *AuthService) ParseToken(token, claim string) (interface{}, error) {
	claimValue, err := s.tokenManager.Parse(token, claim)
	if err != nil {
		return nil, err
	}

	return claimValue, nil
}

func (s *AuthService) GetToken(isAdmin bool) (string, error) {
	return s.tokenManager.NewToken(isAdmin)
}
