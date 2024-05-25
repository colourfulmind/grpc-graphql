CREATE TABLE IF NOT EXISTS users
(
    id        SERIAL PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    content     TEXT NOT NULL,
    comments    BOOL NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id     INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments
(
    id          SERIAL PRIMARY KEY,
    content     TEXT NOT NULL,
    user_id     INTEGER NOT NULL,
    post_id     INTEGER NOT NULL,
    comment_id  INTEGER NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);
