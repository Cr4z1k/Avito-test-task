package repository

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/jmoiron/sqlx"
)

type Banner interface {
	GetBannerLastRevision(tagID, featureID uint64, isAdmin bool) (core.Banner, error)
	GetBannersWithFilter(tagID, featureID *int, limit, offset int) ([]core.BannerWithFilters, error)
	CreateBanner(tagIDs []int, featureID uint64, bannerCnt core.BannerContent, isActive bool) (int, error)
	UpdateBanner(bannerID, featureID uint64, tagIDs []int, newBanner core.Banner) error
	DeleteBanner(bannerID uint64) (uint64, uint64, error)
}

type Feature interface {
	CreateFeatures(featureIDs []int) error
}

type Tag interface {
	CreateTags(tagIDs []int) error
}

type Repository struct {
	Banner
	Feature
	Tag
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Banner:  NewBannerRepo(db),
		Feature: NewFeatureRepo(db),
		Tag:     NewTagRepo(db),
	}
}
