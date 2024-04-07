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

func CreateTablesAndTriggers(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS banner(
		ID SERIAL PRIMARY KEY NOT NULL,
		title TEXT,
		text TEXT,
		url TEXT,
		is_active BOOL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS feature (
		ID SERIAL PRIMARY KEY NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS tag (
		ID SERIAL PRIMARY KEY NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS banner_tags (
		ID SERIAL PRIMARY KEY NOT NULL,
		banner_id INT REFERENCES banner(ID) ON DELETE CASCADE NOT NULL,
		tag_id INT REFERENCES tag(ID) ON DELETE CASCADE NOT NULL,
		CONSTRAINT unique_banner_tag UNIQUE (banner_id, tag_id)
	);
	
	CREATE TABLE IF NOT EXISTS banners_feature (
		ID SERIAL PRIMARY KEY NOT NULL,
		banner_id INT REFERENCES banner(ID) ON DELETE CASCADE UNIQUE NOT NULL,
		feature_id INT REFERENCES feature(ID) ON DELETE CASCADE NOT NULL,
		CONSTRAINT unique_banner_feature UNIQUE (banner_id, feature_id)
	);

	DROP TRIGGER IF EXISTS trigger_insert_banner_feature ON banner_feature_tag;

	CREATE OR REPLACE FUNCTION insert_banner_feature()
	RETURNS TRIGGER AS $$
	BEGIN
		IF (SELECT COUNT(DISTINCT feature_id) FROM banner_feature_tag WHERE banner_id = NEW.banner_id) > 1
		THEN
			RAISE EXCEPTION 'This banner already has a feature';
		END IF;
		
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER trigger_insert_banner_feature
	AFTER INSERT ON banner_feature_tag
	FOR EACH ROW
	EXECUTE FUNCTION insert_banner_feature();
	`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
