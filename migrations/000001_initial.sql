-- +goose Up
CREATE TABLE app_user
(
    id       BIGSERIAL PRIMARY KEY,
    username varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    email    varchar(255) NOT NULL UNIQUE,
    bio      TEXT,
    image    varchar(511)
);

CREATE TABLE article
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT   NOT NULL,
    slug        varchar(255) NOT NULL UNIQUE,
    title       varchar(255) NOT NULL,
    description TEXT         NOT NULL,
    body        TEXT         NOT NULL,
    created_at  timestamptz  NOT NULL,
    updated_at  timestamptz  NOT NULL,
    constraint fk_article_user foreign key (user_id) references app_user (id) ON DELETE CASCADE
);

CREATE TABLE article_favorite
(
    id         BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL,
    user_id    BIGINT NOT NULL,
    constraint fk_article_favorite_article foreign key (article_id) references article (id) ON DELETE CASCADE,
    constraint fk_article_favorite_user foreign key (user_id) references app_user (id) ON DELETE CASCADE
);

CREATE TABLE follow
(
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL,
    follow_id BIGINT NOT NULL,
    constraint fk_follow_user foreign key (user_id) references app_user (id) ON DELETE CASCADE,
    constraint fk_follow_follow_user foreign key (follow_id) references app_user (id) ON DELETE CASCADE
);

CREATE TABLE tag
(
    id   BIGSERIAL PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE article_tag
(
    id         BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL,
    tag_id     BIGINT NOT NULL,
    constraint fk_article_tag_article foreign key (article_id) references article (id) ON DELETE CASCADE,
    constraint fk_article_tag_tag foreign key (tag_id) references tag (id) ON DELETE CASCADE
);

CREATE TABLE comment
(
    id         BIGSERIAL PRIMARY KEY,
    body       TEXT        NOT NULL,
    article_id BIGINT  NOT NULL,
    user_id    BIGINT  NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    constraint fk_comment_user foreign key (user_id) references app_user (id) ON DELETE CASCADE,
    constraint fk_comment_article foreign key (article_id) references article (id) ON DELETE CASCADE
);

CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BYTEA       NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);


-- +goose Down
DROP TABLE sessions;
DROP TABLE comment;
DROP TABLE article_tag;
DROP TABLE tag;
DROP TABLE follow;
DROP TABLE article_favorite;
DROP TABLE article;
DROP TABLE app_user;
