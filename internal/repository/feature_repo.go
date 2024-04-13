package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FeatureRepo struct {
	db *sqlx.DB
}

func NewFeatureRepo(db *sqlx.DB) *FeatureRepo {
	return &FeatureRepo{db: db}
}

func (r *FeatureRepo) CreateFeatures(fetureIDs []int) error {
	query := `INSERT INTO feature VALUES `

	for _, featureID := range fetureIDs {
		query += fmt.Sprintf("(%d), ", featureID)
	}

	query = query[:len(query)-2] + ";"

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}
