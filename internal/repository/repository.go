package repository

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/jmoiron/sqlx"
)

type Banner interface {
	GetBannerLastRevision(tagID, featureID uint64, isAdmin bool) (core.Banner, error)
	GetBannersWithFilter(tagIDs []uint64, featureID uint64, limit, offset int) ([]core.BannerWithFilters, error)
	CreateBanner(tagIDs []int, featureID uint64, bannerCnt core.BannerContent, isActive bool) (int, error)
	UpdateBanner(bannerID, featureID uint64, tagIDs []uint64, NewBanner core.Banner, isActive bool) error
	DeleteBanner(bannerID uint64) (uint64, uint64, error)
}

type Repository struct {
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Banner: NewBannerRepo(db)}
}
