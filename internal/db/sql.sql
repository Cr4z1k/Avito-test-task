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

CREATE OR REPLACE FUNCTION create_banner(
    IN tag_ids INT[],
    IN feature_id INT,
    IN title VARCHAR(255),
    IN text VARCHAR(255),
    IN url VARCHAR(255),
    IN is_active BOOLEAN
)
RETURNS INT AS $$
DECLARE
    banner_id INT;
    tag_id INT;
BEGIN
    BEGIN
        BEGIN
            INSERT INTO banner(title, text, url, is_active, created_at, updated_at)
            VALUES (title, text, url, is_active, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
            RETURNING id INTO banner_id;
            
            FOREACH tag_id IN ARRAY tag_ids
            LOOP
                INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id)
                VALUES (banner_id, feature_id, tag_id);
            END LOOP;
			
			RETURN banner_id;
           
        EXCEPTION WHEN others THEN
			RAISE EXCEPTION 'Ошибка при добавлении данных: %', SQLERRM;
            ROLLBACK;
        END;
    END;
END;
$$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------------

UPDATE banner SET is_active = false WHERE id = 3
SELECT * FROM banner_feature_tag

INSERT INTO banner(title, text, url, is_active, created_at, updated_at)
VALUES ('5', '5', 'http://5.com', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id)
VALUES (5, 2, 1), (5, 1, 2);

INSERT INTO banner(title, text, url, is_active, created_at, updated_at)
VALUES ('5', '5', 'http://5.com', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id)
VALUES (5, 2, 1), (5, 1, 2);

SELECT 'DROP FUNCTION IF EXISTS ' || proname || '(' || oidvectortypes(proargtypes) || ');'
FROM pg_proc
WHERE proname = 'create_banner';

SELECT create_banner(
    ARRAY[1, 2],
    3,
    '5',
    '5',
    'http://5.com',
    true
);

DELETE FROM banner_feature_tag WHERE ID > 43;
DELETE FROM BANNER WHERE ID > 4;

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