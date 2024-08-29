BEGIN;

ALTER TABLE users ADD COLUMN type INTEGER DEFAULT 0;
ALTER TABLE users DROP COLUMN created_at;

DROP TABLE post_comments;
DROP TABLE posts;
DROP TABLE user_connections;
DROP TABLE user_login_history;

COMMIT;