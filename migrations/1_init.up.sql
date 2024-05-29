CREATE TABLE IF NOT EXISTS users
(
    id        SERIAL PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id                  SERIAL PRIMARY KEY,
    title               TEXT NOT NULL,
    content             TEXT NOT NULL,
    comments_allowed    BOOL NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id             INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments
(
    id          SERIAL PRIMARY KEY,
    content     TEXT NOT NULL,
    user_id     INTEGER NOT NULL,
    post_id     INTEGER NOT NULL,
    parent_id   INTEGER,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (parent_id) REFERENCES comments(id)
);

INSERT INTO users (email, password) VALUES ('test@test.com', '12345678');
INSERT INTO users (email, password) VALUES ('test2@test.com', '12345678');
INSERT INTO users (email, password) VALUES ('test3@test.com', '12345678');

INSERT INTO posts (title, content, comments, user_id) VALUES ('Test1', 'hello, world', true);
INSERT INTO posts (title, content, comments, user_id) VALUES ('Test2', 'hello, world', true);
INSERT INTO posts (title, content, comments, user_id) VALUES ('Test3', 'hello, world', false);

INSERT INTO comments (content, user_id, post_id, parent_id) VALUES ('TestComment1', 1, 1, 0);
INSERT INTO comments (content, user_id, post_id, parent_id) VALUES ('TestComment2', 1, 1, 1);
INSERT INTO comments (content, user_id, post_id, parent_id) VALUES ('TestComment3', 1, 1, 1);