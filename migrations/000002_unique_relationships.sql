-- +goose Up
DELETE FROM article_favorite duplicate
USING article_favorite original
WHERE duplicate.id > original.id
  AND duplicate.article_id = original.article_id
  AND duplicate.user_id = original.user_id;

DELETE FROM follow duplicate
USING follow original
WHERE duplicate.id > original.id
  AND duplicate.user_id = original.user_id
  AND duplicate.follow_id = original.follow_id;

DELETE FROM article_tag duplicate
USING article_tag original
WHERE duplicate.id > original.id
  AND duplicate.article_id = original.article_id
  AND duplicate.tag_id = original.tag_id;

ALTER TABLE article_favorite
    ADD CONSTRAINT article_favorite_article_user_unique UNIQUE (article_id, user_id);

ALTER TABLE follow
    ADD CONSTRAINT follow_user_follow_unique UNIQUE (user_id, follow_id);

ALTER TABLE article_tag
    ADD CONSTRAINT article_tag_article_tag_unique UNIQUE (article_id, tag_id);

-- +goose Down
ALTER TABLE article_tag DROP CONSTRAINT article_tag_article_tag_unique;
ALTER TABLE follow DROP CONSTRAINT follow_user_follow_unique;
ALTER TABLE article_favorite DROP CONSTRAINT article_favorite_article_user_unique;
