-- +goose Up
create table user (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  username varchar(255) UNIQUE,
  password varchar(255),
  email varchar(255) UNIQUE,
  bio text,
  image varchar(511)
);

create table article (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  user_id int,
  slug varchar(255) UNIQUE,
  title varchar(255),
  description text,
  body text,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  constraint fk_article_user foreign key (user_id) references user(id) ON DELETE CASCADE
);

create table article_favorite (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  article_id int not null,
  user_id int not null,
  constraint fk_article_favorite_article foreign key (article_id) references article(id) ON DELETE CASCADE,
  constraint fk_article_favorite_user foreign key (user_id) references user(id) ON DELETE CASCADE
);

create table follow (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  user_id int not null,
  follow_id int not null,
  constraint fk_follow_user foreign key (user_id) references user(id) ON DELETE CASCADE,
  constraint fk_follow_follow_user foreign key (follow_id) references user(id) ON DELETE CASCADE
);

create table tag (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name varchar(255)
);

create table article_tag (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  article_id int not null,
  tag_id int not null,
  constraint fk_article_tag_article foreign key (article_id) references article(id) ON DELETE CASCADE,
  constraint fk_article_tag_tag foreign key (tag_id) references tag(id) ON DELETE CASCADE
);

create table comment (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  body text,
  article_id int not null,
  user_id int not null,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  constraint fk_comment_user foreign key (user_id) references user(id) ON DELETE CASCADE,
  constraint fk_comment_article foreign key (article_id) references article(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE comment;
DROP TABLE article_tag;
DROP TABLE tag;
DROP TABLE follow;
DROP TABLE article_favorite;
DROP TABLE article;
DROP TABLE user;
