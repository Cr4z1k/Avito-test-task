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

CREATE TABLE IF NOT EXISTS banner_feature_tag (
	ID SERIAL PRIMARY KEY NOT NULL,
	banner_id INT REFERENCES banner(ID) ON DELETE CASCADE NOT NULL,
	feature_id INT REFERENCES feature(ID) ON DELETE CASCADE NOT NULL,
	tag_id INT REFERENCES tag(ID) ON DELETE CASCADE NOT NULL,
	CONSTRAINT unique_feature_tag UNIQUE (feature_id, tag_id)
);

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

----------------------------------------------------------------------------------------------------------

UPDATE banner SET is_active = false WHERE id = 3
SELECT * FROM banner

INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id) VALUES (1, 2, 1), (1, 1, 2);
INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id) VALUES (2, 2, 2), (2, 2, 3), (2, 2, 4);
INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id) VALUES (3, 1, 3), (3, 1, 4), (3, 1, 5);
INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id) VALUES (4, 3, 4), (4, 3, 5);

INSERT INTO banner(title, text, url, is_active, created_at, updated_at) VALUES ('4', '4', 'http://4.com', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO feature VALUES (1), (2), (3);
INSERT INTO tag VALUES (1), (2), (3), (4), (5);

SELECT * FROM banner
SELECT * FROM feature
SELECT * FROM tag
SELECT * FROM banner_tags
SELECT * FROM banners_feature

SELECT title, text, url
	FROM banner b
	JOIN banner_feature_tag bft ON b.id = bft.banner_id
	WHERE bft.tag_id = 4 AND bft.feature_id = 1