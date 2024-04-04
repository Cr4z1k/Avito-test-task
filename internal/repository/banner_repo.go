package repository

import (
	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/jmoiron/sqlx"
)

type BannerRepo struct {
	db *sqlx.DB
}

func NewBannerRepo(db *sqlx.DB) *BannerRepo {
	return &BannerRepo{db: db}
}

func (r *BannerRepo) GetBanner(tag_id []int, feature_id int, isLastVer bool) (core.Banner, error) {

	return core.Banner{}, nil
}

func (r *BannerRepo) GetBannersWithFilter(tag_id []int, feature_id, limit, offset int) ([]core.Banner, error) {

	return nil, nil
}

func (r *BannerRepo) CreateBanner(tags []int, feature int, bannerCnt core.Banner, isActive bool) error {

	return nil
}

func (r *BannerRepo) UpdateBanner(bannerID int, tags []int, feature int, bannerCnt core.Banner, isActive bool) error {

	return nil
}

func (r *BannerRepo) DeleteBanner(bannerID int) error {

	return nil
}
