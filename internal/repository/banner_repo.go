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

func (r *BannerRepo) GetBannersWithFilter(tagID, featureID *int, limit, offset int) ([]core.BannerWithFilters, error) {
	query := `
		SELECT b.id, ARRAY(SELECT tag_id FROM banner_feature_tag bft WHERE banner_id = b.id) AS tag_ids, feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at
		FROM banner_feature_tag bft
		JOIN banner b ON b.id = bft.banner_id
	`

	args := []interface{}{}

	if tagID != nil && featureID != nil {
		query += ` WHERE bft.tag_id = $1 AND bft.feature_id = $2 GROUP BY b.id, feature_id LIMIT $3 OFFSET $4`

		args = append(args, *tagID, *featureID, limit, offset)
	} else if tagID == nil && featureID != nil {
		query += ` WHERE bft.feature_id = $1 GROUP BY b.id, feature_id LIMIT $2 OFFSET $3`

		args = append(args, *featureID, limit, offset)
	} else if tagID != nil && featureID == nil {
		query += ` WHERE bft.tag_id = $1 GROUP BY b.id, feature_id LIMIT $2 OFFSET $3`

		args = append(args, *tagID, limit, offset)
	} else {
		query += ` GROUP BY b.id, feature_id LIMIT $1 OFFSET $2`

		args = append(args, limit, offset)
	}

	resultRows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var banners []core.BannerWithFilters

	for resultRows.Next() {
		var banner core.BannerWithFilters

		err := resultRows.Scan(&banner.BannerID, pq.Array(&banner.TagIDs), &banner.FeatureID,
			&banner.Content.Title, &banner.Content.Text, &banner.Content.Url, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			return nil, err
		}

		banners = append(banners, banner)
	}

	if err := resultRows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
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

func (r *BannerRepo) UpdateBanner(bannerID, featureID uint64, tagIDs []int, newBanner core.Banner) error {
	tagIDsArray := pq.Array(tagIDs)

	query := `
		CALL update_banner($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := r.db.Exec(query, bannerID, newBanner.Title, newBanner.Text, newBanner.Url, newBanner.IsActive, featureID, tagIDsArray); err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) DeleteBanner(bannerID uint64) (uint64, uint64, error) {
	query := `
		SELECT tag_id, feature_id FROM banner_feature_tag WHERE banner_id = $1 LIMIT 1;
	`

	var featureTagIds struct {
		TadID     uint64 `db:"tag_id"`
		FeatureID uint64 `db:"feature_id"`
	}

	if err := r.db.QueryRowx(query, bannerID).StructScan(&featureTagIds); err != nil {
		return 0, 0, err
	}

	query = `
		DELETE FROM banner WHERE id = $1;
	`

	if _, err := r.db.Exec(query, bannerID); err != nil {
		return 0, 0, err
	}

	return featureTagIds.TadID, featureTagIds.FeatureID, nil
}
