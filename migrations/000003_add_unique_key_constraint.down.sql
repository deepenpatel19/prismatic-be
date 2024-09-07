BEGIN;

ALTER TABLE user_connections DROP CONSTRAINT unique_user_and_friend;

COMMIT;