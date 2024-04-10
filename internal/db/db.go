package db

import (
	"os"

	"github.com/Cr4z1k/Avito-test-task/internal/config/dbconf"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbconf.GetConnectionString())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitialQuery(db *sqlx.DB) error {
	pathToDir := "./internal/db/"

	initQueryFile, err := os.ReadFile(pathToDir + "init_sql.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(initQueryFile)); err != nil {
		return err
	}

	insertTagsAndFeaturesFile, err := os.ReadFile(pathToDir + "insert.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(insertTagsAndFeaturesFile)); err != nil {
		return err
	}

	return nil
}
