package repository

import (
	"database/sql"

	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/jmoiron/sqlx"
)

type BannerRepo struct {
	db *sqlx.DB
}

func NewBannerRepo(db *sqlx.DB) *BannerRepo {
	return &BannerRepo{db: db}
}

// TODO: use IN CU funcs
// func resetBannerUpdTime(db *sqlx.DB, banner_id uint64) error {
// 	query := `
// 	UPDATE banner
// 	SET updated_at = $2
// 	WHERE id = $1;
// 	`

// 	_, err := db.Exec(query, banner_id, time.Now())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *BannerRepo) GetBannerLastRevision(tag_id, feature_id uint64, isAdmin bool) (core.Banner, error) {
	var banner core.Banner

	query := `
	SELECT b.id, title, text, url, is_active, created_at, updated_at
	FROM banner b
	JOIN banner_feature_tag bft ON b.id = bft.banner_id
	WHERE bft.tag_id = $1 AND bft.feature_id = $2
	`

	if !isAdmin {
		query += ` AND is_active;`
	}

	err := r.db.QueryRowx(query, tag_id, feature_id).StructScan(&banner)
	if err == sql.ErrNoRows {
		return core.Banner{}, nil
	} else if err != nil {
		return core.Banner{}, err
	}

	return banner, nil
}

func (r *BannerRepo) GetBannersWithFilter(tag_ids []uint64, feature_id uint64, limit, offset int) ([]core.BannerWithFilters, error) {

	return nil, nil
}

func (r *BannerRepo) CreateBanner(tag_ids []uint64, feature_id uint64, bannerCnt core.Banner, isActive bool) error {

	return nil
}

func (r *BannerRepo) UpdateBanner(bannerID, feature_id uint64, tag_ids []uint64, NewBanner core.Banner, isActive bool) error {

	return nil
}

func (r *BannerRepo) DeleteBanner(bannerID uint64) error {

	return nil
}
