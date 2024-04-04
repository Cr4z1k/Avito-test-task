package db

import (
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
