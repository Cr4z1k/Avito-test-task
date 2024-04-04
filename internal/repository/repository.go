package repository

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/jmoiron/sqlx"
)

type Banner interface {
	GetBanner(tag_id []int, feature_id int, isLastVer bool) (core.Banner, error)
	GetBannersWithFilter(tag_id []int, feature_id, limit, offset int) ([]core.Banner, error)
	CreateBanner(tags []int, feature int, bannerCnt core.Banner, isActive bool) error
	UpdateBanner(bannerID int, tags []int, feature int, bannerCnt core.Banner, isActive bool) error
	DeleteBanner(bannerID int) error
}

type Repository struct {
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Banner: NewBannerRepo(db)}
}
