-- +goose Up
CREATE TABLE user
(
    id       INTEGER      NOT NULL PRIMARY KEY AUTOINCREMENT,
    username varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    email    varchar(255) NOT NULL UNIQUE,
    bio      TEXT,
    image    varchar(511)
);

CREATE TABLE article
(
    id          INTEGER      NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER      NOT NULL,
    slug        varchar(255) NOT NULL UNIQUE,
    title       varchar(255) NOT NULL,
    description TEXT         NOT NULL,
    body        TEXT         NOT NULL,
    created_at  INTEGER      NOT NULL,
    updated_at  INTEGER      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    constraint fk_article_user foreign key (user_id) references user (id) ON DELETE CASCADE
);

CREATE TABLE article_favorite
(
    id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    constraint fk_article_favorite_article foreign key (article_id) references article (id) ON DELETE CASCADE,
    constraint fk_article_favorite_user foreign key (user_id) references user (id) ON DELETE CASCADE
);

CREATE TABLE follow
(
    id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id   INTEGER NOT NULL,
    follow_id INTEGER NOT NULL,
    constraint fk_follow_user foreign key (user_id) references user (id) ON DELETE CASCADE,
    constraint fk_follow_follow_user foreign key (follow_id) references user (id) ON DELETE CASCADE
);

CREATE TABLE tag
(
    id   INTEGER      NOT NULL PRIMARY KEY AUTOINCREMENT,
    name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE article_tag
(
    id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER NOT NULL,
    tag_id     INTEGER NOT NULL,
    constraint fk_article_tag_article foreign key (article_id) references article (id) ON DELETE CASCADE,
    constraint fk_article_tag_tag foreign key (tag_id) references tag (id) ON DELETE CASCADE
);

CREATE TABLE comment
(
    id         INTEGER   NOT NULL PRIMARY KEY AUTOINCREMENT,
    body       TEXT      NOT NULL,
    article_id INTEGER   NOT NULL,
    user_id    INTEGER   NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP,
    constraint fk_comment_user foreign key (user_id) references user (id) ON DELETE CASCADE,
    constraint fk_comment_article foreign key (article_id) references article (id) ON DELETE CASCADE
);

CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BLOB NOT NULL,
    expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);


-- +goose Down
DROP TABLE comment;
DROP TABLE article_tag;
DROP TABLE tag;
DROP TABLE follow;
DROP TABLE article_favorite;
DROP TABLE article;
DROP TABLE user;
DROP TABLE sessions;
