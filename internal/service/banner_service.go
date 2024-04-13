package service

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
)

type BannerService struct {
	r     repository.Banner
	Cache map[string]core.Banner
}

func NewBannerService(r repository.Banner) *BannerService {
	return &BannerService{
		r:     r,
		Cache: make(map[string]core.Banner),
	}
}

func (s *BannerService) GetBanner(tagID, featureID uint64, isLastVer, isAdmin bool) (core.BannerContent, error) {
	var banner core.Banner
	var err error

	cacheKey := strconv.FormatUint(tagID, 10) + strconv.FormatUint(featureID, 10)

	if !isLastVer {
		banner, ok := s.Cache[cacheKey]
		if ok && time.Now().Add(time.Minute*5).Before(banner.UpdatedAt) {
			if banner.IsActive || isAdmin {
				return core.BannerContent{
					Title: banner.Title,
					Text:  banner.Text,
					Url:   banner.Url,
				}, nil
			}
		}
	}

	banner, err = s.r.GetBannerLastRevision(tagID, featureID, isAdmin)
	if err != nil {
		return core.BannerContent{}, err
	}

	s.Cache[cacheKey] = banner

	bannerContent := core.BannerContent{
		Title: banner.Title,
		Text:  banner.Text,
		Url:   banner.Url,
	}

	return bannerContent, nil
}

func (s *BannerService) GetBannersWithFilter(tagID, featureID *int, limit, offset int) ([]core.BannerWithFilters, error) {
	searchResult, err := s.r.GetBannersWithFilter(tagID, featureID, limit, offset)
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

func (s *BannerService) CreateBanner(tagIDs []int, featureID uint64, bannerCnt core.BannerContent, isActive bool) (int, error) {
	linkRegex := regexp.MustCompile(`^(http[s]?:\/\/)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}([\/?#].*)?$`)

	if !linkRegex.MatchString(bannerCnt.Url) {
		return -1, errors.New("wrong format for URL string")
	}

	id, err := s.r.CreateBanner(tagIDs, featureID, bannerCnt, isActive)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *BannerService) UpdateBanner(bannerID, featureID uint64, tagIDs []int, bannerCnt core.BannerContent, isActive bool) error {
	newBanner := core.Banner{
		Title:    bannerCnt.Title,
		Text:     bannerCnt.Text,
		Url:      bannerCnt.Url,
		IsActive: isActive,
	}

	if err := s.r.UpdateBanner(bannerID, featureID, tagIDs, newBanner); err != nil {
		return err
	}

	return nil
}

func (s *BannerService) DeleteBanner(bannerID int) error {
	if bannerID < 0 {
		return errors.New("ID cannot be less than 0")
	}

	tagID, featureID, err := s.r.DeleteBanner(uint64(bannerID))
	if err != nil {
		return err
	}

	key := strconv.FormatUint(tagID, 10) + strconv.FormatUint(featureID, 10)

	delete(s.Cache, key)

	return nil
}
