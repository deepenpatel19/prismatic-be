BEGIN;

ALTER TABLE user_connections ADD CONSTRAINT unique_user_and_friend UNIQUE (user_id, friend_id);

COMMIT;