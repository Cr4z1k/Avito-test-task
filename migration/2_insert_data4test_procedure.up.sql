TRUNCATE banner CASCADE;
TRUNCATE tag CASCADE;
TRUNCATE feature CASCADE;

INSERT INTO tag(id) VALUES (1), (2), (3), (4), (5);
INSERT INTO feature(id) VALUES (1), (2), (3);

INSERT INTO banner(id, title, text, url, is_active, created_at, updated_at) VALUES
(1, 'test_1', 'test_1', 'http://test_1.com', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'test_2', 'test_2', 'http://test_2.com', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'test_3', 'test_3', 'http://test_3.com', false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'test_4', 'test_4', 'http://test_4.com', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id) VALUES
(1, 1, 1), (1, 1, 2),
(2, 2, 2), (2, 2, 3), (2, 2, 4),
(3, 1, 3), (3, 1, 4), (3, 1, 5),
(4, 3, 4), (4, 3, 5);