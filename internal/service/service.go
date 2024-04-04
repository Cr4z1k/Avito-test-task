package service

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
)

type Banner interface {
	GetBanner(tag_id []int, feature_id int, isLastVer bool) (core.Banner, error)
	GetBannersWithFilter(tag_id []int, feature_id, limit, offset int) ([]core.Banner, error)
	CreateBanner(tags []int, feature int, bannerCnt core.Banner, isActive bool) error
	UpdateBanner(bannerID int, tags []int, feature int, bannerCnt core.Banner, isActive bool) error
	DeleteBanner(bannerID int) error
}

type Service struct {
	Banner
}

func NewService(r *repository.Repository) *Service {
	return &Service{Banner: NewBannerService(r.Banner)}
}
