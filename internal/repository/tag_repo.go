package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TagRepo struct {
	db *sqlx.DB
}

func NewTagRepo(db *sqlx.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) CreateTags(tagIDs []int) error {
	query := `INSERT INTO tag VALUES `

	for _, tagID := range tagIDs {
		query += fmt.Sprintf("(%d), ", tagID)
	}

	query = query[:len(query)-2] + ";"

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}
