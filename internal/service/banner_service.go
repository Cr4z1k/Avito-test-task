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
	cache map[string]core.Banner
}

func NewBannerService(r repository.Banner) *BannerService {
	return &BannerService{
		r:     r,
		cache: make(map[string]core.Banner),
	}
}

func (s *BannerService) GetBanner(tagID, featureID uint64, isLastVer, isAdmin bool) (core.BannerContent, error) {
	var banner core.Banner
	var err error

	cacheKey := strconv.FormatUint(tagID, 10) + strconv.FormatUint(featureID, 10)

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

	banner, err = s.r.GetBannerLastRevision(tagID, featureID, isAdmin)
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

func (s *BannerService) UpdateBanner(bannerID int, tags []int, feature int, bannerCnt core.Banner, isActive bool) error {

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

	delete(s.cache, key)

	return nil
}
