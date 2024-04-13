package service

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
	"github.com/Cr4z1k/Avito-test-task/pkg/auth"
)

type Banner interface {
	GetBanner(tagID, featureID uint64, isLastVer, isAdmin bool) (core.BannerContent, error)
	GetBannersWithFilter(tagID, featureID *int, limit, offset int) ([]core.BannerWithFilters, error)
	CreateBanner(tagIDs []int, featureID uint64, bannerCnt core.BannerContent, isActive bool) (int, error)
	UpdateBanner(bannerID, featureID uint64, tagIDs []int, bannerCnt core.BannerContent, isActive bool) error
	DeleteBanner(bannerID int) error
}

type Auth interface {
	ParseToken(token, claim string) (interface{}, error)
	GetToken(isAdmin bool) (string, error)
}

type Feature interface {
	CreateFeatures(featureIDs []int) error
}

type Tag interface {
	CreateTags(tagIDs []int) error
}

type Service struct {
	Banner
	Auth
	Feature
	Tag
}

func NewService(r *repository.Repository, t auth.TokenManager) *Service {
	return &Service{
		Banner:  NewBannerService(r.Banner),
		Auth:    NewAuthService(t),
		Feature: NewFeatureService(r.Feature),
		Tag:     NewTagService(r.Tag),
	}
}
