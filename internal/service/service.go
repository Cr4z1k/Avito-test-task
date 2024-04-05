package service

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
	"github.com/Cr4z1k/Avito-test-task/pkg/auth"
)

type Banner interface {
	GetBanner(tag_id []int, feature_id int, isLastVer bool) (core.Banner, error)
	GetBannersWithFilter(tag_id []int, feature_id, limit, offset int) ([]core.Banner, error)
	CreateBanner(tags []int, feature int, bannerCnt core.Banner, isActive bool) error
	UpdateBanner(bannerID int, tags []int, feature int, bannerCnt core.Banner, isActive bool) error
	DeleteBanner(bannerID int) error
}

type Auth interface {
	ParseToken(token, claim string) (interface{}, error)
	GetToken(isAdmin bool) (string, error)
}

type Service struct {
	Banner
	Auth
}

func NewService(r *repository.Repository, t auth.TokenManager) *Service {
	return &Service{
		Banner: NewBannerService(r.Banner),
		Auth:   NewAuthService(t),
	}
}
