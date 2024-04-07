package repository

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/jmoiron/sqlx"
)

type Banner interface {
	GetBannerLastRevision(tag_id, feature_id uint64, isAdmin bool) (core.Banner, error)
	GetBannersWithFilter(tag_ids []uint64, feature_id uint64, limit, offset int) ([]core.BannerWithFilters, error)
	CreateBanner(tag_ids []uint64, feature_id uint64, bannerCnt core.Banner, isActive bool) error
	UpdateBanner(bannerID, feature_id uint64, tag_ids []uint64, NewBanner core.Banner, isActive bool) error
	DeleteBanner(bannerID uint64) error
}

type Repository struct {
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Banner: NewBannerRepo(db)}
}
