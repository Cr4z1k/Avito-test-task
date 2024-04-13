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

CREATE OR REPLACE PROCEDURE update_banner(
    banner_id_p INT, 
    title_p TEXT, 
    text_p TEXT, 
    url_p TEXT, 
    is_active_p BOOL, 
    feature_id_p INT, 
    tag_ids INT[]
)
LANGUAGE plpgsql
AS $$
DECLARE
    bft_ids INT[];
    i INT;
    bft_ids_len INT;
    tag_ids_len INT;
BEGIN
    UPDATE banner 
    SET title = title_p, text = text_p, url = url_p, is_active = is_active_p, updated_at = CURRENT_TIMESTAMP
    WHERE id = banner_id_p;

    SELECT array_agg(id) FROM banner_feature_tag WHERE banner_id = banner_id_p
    INTO bft_ids;
	
	IF bft_ids IS NULL THEN
        RAISE EXCEPTION 'No banner with such ID';
        RETURN;
    END IF;

    bft_ids_len := array_length(bft_ids, 1);
    tag_ids_len := array_length(tag_ids, 1);

    CASE
        WHEN bft_ids_len = tag_ids_len THEN
            FOR i IN 1..bft_ids_len LOOP
                UPDATE banner_feature_tag
                SET tag_id = tag_ids[i], feature_id = feature_id_p
                WHERE id = bft_ids[i];
            END LOOP;

        WHEN bft_ids_len > tag_ids_len THEN
            FOR i IN 1..tag_ids_len LOOP
                UPDATE banner_feature_tag
                SET tag_id = tag_ids[i], feature_id = feature_id_p
                WHERE id = bft_ids[i];
            END LOOP;

            FOR i IN tag_ids_len + 1..bft_ids_len LOOP
                DELETE FROM banner_feature_tag
                WHERE id = bft_ids[i];
            END LOOP;

        ELSE
            FOR i IN 1..bft_ids_len LOOP
                UPDATE banner_feature_tag
                SET tag_id = tag_ids[i], feature_id = feature_id_p
                WHERE id = bft_ids[i];
            END LOOP;

            FOR i IN bft_ids_len + 1..tag_ids_len LOOP
                INSERT INTO banner_feature_tag (banner_id, feature_id, tag_id)
                VALUES (banner_id_p, feature_id_p, tag_ids[i]);
            END LOOP;
    END CASE;

    EXCEPTION WHEN others THEN
        RAISE EXCEPTION 'Ошибка при изменении данных: %', SQLERRM;
        ROLLBACK;
END;
$$;