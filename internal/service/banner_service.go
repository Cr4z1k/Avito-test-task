package service

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
)

type BannerService struct {
	r repository.Banner
}

func NewBannerService(r repository.Banner) *BannerService {
	return &BannerService{r: r}
}

func (s *BannerService) GetBanner(tag_id []int, feature_id int, isLastVer bool) (core.Banner, error) {

	return core.Banner{}, nil
}

func (s *BannerService) GetBannersWithFilter(tag_id []int, feature_id, limit, offset int) ([]core.Banner, error) {

	return nil, nil
}

func (s *BannerService) CreateBanner(tags []int, feature int, bannerCnt core.Banner, isActive bool) error {

	return nil
}

func (s *BannerService) UpdateBanner(bannerID int, tags []int, feature int, bannerCnt core.Banner, isActive bool) error {

	return nil
}

func (s *BannerService) DeleteBanner(bannerID int) error {

	return nil
}
