BEGIN;

ALTER TABLE users DROP COLUMN type;
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT NOW();

CREATE TABLE user_login_history(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGSERIAL,
    login_at TIMESTAMP DEFAULT NOW(),
    ip_address TEXT DEFAULT '',

    CONSTRAINT user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
                ON DELETE CASCADE
);

CREATE TABLE user_connections(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGSERIAL,
    friend_id BIGSERIAL,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
                ON DELETE CASCADE,

    CONSTRAINT friend_id
        FOREIGN KEY(friend_id)
            REFERENCES users(id)
                ON DELETE CASCADE
);

CREATE TABLE posts(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGSERIAL,
    title VARCHAR(100) DEFAULT '' NOT NULL,
    description TEXT DEFAULT '',
    file TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
                ON DELETE CASCADE
);

CREATE TABLE post_comments(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGSERIAL,
    post_id BIGSERIAL,
    comment TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
                ON DELETE CASCADE,

    CONSTRAINT post_id
        FOREIGN KEY(post_id)
            REFERENCES posts(id)
                ON DELETE CASCADE
);

COMMIT;