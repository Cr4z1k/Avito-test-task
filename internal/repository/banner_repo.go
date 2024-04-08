package repository

import (
	"database/sql"
	"time"

	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type BannerRepo struct {
	db *sqlx.DB
}

func NewBannerRepo(db *sqlx.DB) *BannerRepo {
	return &BannerRepo{db: db}
}

// TODO: use IN U func
func resetBannerUpdTime(db *sqlx.DB, banner_id uint64) error {
	query := `
	UPDATE banner
	SET updated_at = $2
	WHERE id = $1;
	`

	_, err := db.Exec(query, banner_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) GetBannerLastRevision(tagID, featureID uint64, isAdmin bool) (core.Banner, error) {
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

	err := r.db.QueryRowx(query, tagID, featureID).StructScan(&banner)
	if err == sql.ErrNoRows {
		return core.Banner{}, nil
	} else if err != nil {
		return core.Banner{}, err
	}

	if err := resetBannerUpdTime(r.db, banner.ID); err != nil {
		return core.Banner{}, err
	}

	return banner, nil
}

func (r *BannerRepo) GetBannersWithFilter(tagIDs []uint64, featureID uint64, limit, offset int) ([]core.BannerWithFilters, error) {

	return nil, nil
}

func (r *BannerRepo) CreateBanner(tagIDs []int, featureID uint64, bannerCnt core.BannerContent, isActive bool) (int, error) {
	var id int

	tagIDsArray := pq.Array(tagIDs)

	query := `
	SELECT create_banner($1, $2, $3, $4, $5, $6);
	`

	if err := r.db.QueryRowx(query, tagIDsArray, featureID, bannerCnt.Title, bannerCnt.Text, bannerCnt.Url, isActive).Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *BannerRepo) UpdateBanner(bannerID, featureID uint64, tagIDs []uint64, NewBanner core.Banner, isActive bool) error {

	return nil
}

func (r *BannerRepo) DeleteBanner(bannerID uint64) (uint64, uint64, error) {
	query := `
		SELECT tag_id, featue_id FROM banner_feature_tag WHERE banner_id = $1 LIMIT 1;
	`

	var featureTagIds struct {
		tag_id     uint64
		feature_id uint64
	}

	if err := r.db.QueryRowx(query, bannerID).StructScan(featureTagIds); err != nil {
		return 0, 0, err
	}

	query = `
		DELETE FROM banner WHERE id = $1;
	`

	if _, err := r.db.Exec(query, bannerID); err != nil {
		return 0, 0, err
	}

	return featureTagIds.tag_id, featureTagIds.feature_id, nil
}
