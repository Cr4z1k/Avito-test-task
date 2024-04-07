package service

import (
	"strconv"
	"time"

	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
)

type BannerService struct {
	r     repository.Banner
	cache map[string]core.Banner
}

func NewBannerService(r repository.Banner) *BannerService {
	return &BannerService{
		r:     r,
		cache: make(map[string]core.Banner),
	}
}

func (s *BannerService) GetBanner(tag_id, feature_id uint64, isLastVer, isAdmin bool) (core.BannerContent, error) {
	var banner core.Banner
	var err error

	cacheKey := strconv.FormatUint(tag_id, 10) + strconv.FormatUint(feature_id, 10)

	if !isLastVer {
		banner, ok := s.cache[cacheKey]
		if ok {
			if time.Now().Add(time.Minute * 5).Before(banner.UpdatedAt) {
				return core.BannerContent{
					Title: banner.Title,
					Text:  banner.Text,
					Url:   banner.Url,
				}, nil
			}
		}
	}

	banner, err = s.r.GetBannerLastRevision(tag_id, feature_id, isAdmin)
	if err != nil {
		return core.BannerContent{}, err
	}

	s.cache[cacheKey] = banner

	bannerContent := core.BannerContent{
		Title: banner.Title,
		Text:  banner.Text,
		Url:   banner.Url,
	}

	return bannerContent, nil
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
