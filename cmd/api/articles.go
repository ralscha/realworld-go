package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
)

func (app *application) articlesFeed(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	articles, err := models.Articles().Feed(r.Context(), app.db, user.ID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, articles)
}

func (app *application) articlesGet(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, article)
}

func (app *application) articleGet(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, article)
}

func (app *application) articlesCreate(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	text := slug.Make(...)
	var input models.ArticleCreateInput
	if err := app.readJSON(w, r, &input); err != nil {
		return
	}

	article, err := input.Create(r.Context(), app.db, user.ID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, article)
}

func (app *application) articlesUpdate(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")

	var input models.ArticleUpdateInput
	if err := app.readJSON(w, r, &input); err != nil {
		return
	}

	article, err := input.Update(r.Context(), app.db, user.ID, slug)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, article)
}

func (app *application) articlesDelete(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	if article.AuthorID != user.ID {
		response.NotFound(w, r)
		return
	}

	if _, err := article.Delete(r.Context(), app.db); err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}


func (app *application) getArticleById(articleId, userId int64) (dto.Article, error) {
	article, err := models.Articles(
		qm.Select(
			models.ArticleColumns.ID,
			models.ArticleColumns.UserID,
			models.ArticleColumns.Title,
			models.ArticleColumns.Description,
			models.ArticleColumns.Body,
			models.ArticleColumns.Slug,
			models.ArticleColumns.CreatedAt,
			models.ArticleColumns.UpdatedAt,
		),
		models.ArticleWhere.ID.EQ(articleId)).One(app.ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	return app.getArticle(article, userId)
}

func (app *application) getArticleBySlug(slug string, userId int64) (dto.Article, error) {
	article, err := models.Articles(
		qm.Select(
			models.ArticleColumns.ID,
			models.ArticleColumns.UserID,
			models.ArticleColumns.Title,
			models.ArticleColumns.Description,
			models.ArticleColumns.Body,
			models.ArticleColumns.Slug,
			models.ArticleColumns.CreatedAt,
			models.ArticleColumns.UpdatedAt,
		), models.ArticleWhere.Slug.EQ(slug)).One(app.ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	return app.getArticle(article, userId)
}

func (app *application) getArticle(article *models.Article, userId int64) (dto.Article, error) {
	author, err := models.Users(qm.Select(models.UserColumns.ID, models.UserColumns.Username,
		models.UserColumns.Bio, models.UserColumns.Image),
		models.UserWhere.ID.EQ(article.UserID.Int64)).One(app.ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	following, err := models.Follows(models.FollowWhere.UserID.EQ(userId), models.FollowWhere.FollowID.EQ(author.ID)).
		Exists(app.ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	authorProfile := dto.Profile{
		Username:  author.Username,
		Bio:       author.Bio.String,
		Image:     author.Image.String,
		Following: following,
	}

	favorited, err := models.ArticleFavorites(models.ArticleFavoriteWhere.UserID.EQ(userId),
		models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).
		Exists(app.ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	tags, err := models.Tags(qm.Select(models.TagColumns.Name),
		qm.InnerJoin(models.TableNames.ArticleTag+" ON "+models.TableNames.ArticleTag+"."+models.ArticleTagColumns.TagID+" = "+models.TableNames.Tag+"."+models.TagColumns.ID),
		models.ArticleTagWhere.ArticleID.EQ(article.ID)).
		All(app.ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	tagList := make([]string, len(tags))
	for _, tag := range tags {
		tagList = append(tagList, tag.Name)
	}

	favoritesCount, err := models.ArticleFavorites(models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).
		Count(app.ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	articleDto := dto.Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description.String,
		Body:           article.Body.String,
		TagList:        tagList,
		//CreatedAt:      article.CreatedAt,
		//UpdatedAt:      article.UpdatedAt,
		Favorited:      favorited,
		FavoritesCount: int(favoritesCount),
	}
	/*
			return new Article(record.get(ARTICLE.SLUG), record.get(ARTICLE.TITLE),
					record.get(ARTICLE.DESCRIPTION), record.get(ARTICLE.BODY),
					dsl.select(TAG.NAME).from(TAG).innerJoin(ARTICLE_TAG).onKey()
							.where(ARTICLE_TAG.ARTICLE_ID.eq(articleId))
							.fetchSet(TAG.NAME),
					record.get(ARTICLE.CREATED_AT), record.get(ARTICLE.UPDATED_AT),
					favorited, favoritesCount, author);
	 */


	return articleDto, nil
}

/*
	public static Article getArticle(DSLContext dsl, ArticleRecord record, long userId) {

		if (record != null) {
			long articleUserId = record.get(ARTICLE.USER_ID);
			long articleId = record.getId();

			AppUserRecord authorRecord = dsl.selectFrom(APP_USER)
					.where(APP_USER.ID.eq(articleUserId)).fetchOne();

			Profile author = null;
			if (authorRecord != null) {
				boolean following = dsl.selectCount().from(FOLLOW)
						.where(FOLLOW.USER_ID.eq(userId)
								.and(FOLLOW.FOLLOW_ID.eq(authorRecord.getId())))
						.fetchOne(0, int.class) == 1;

				author = new Profile(authorRecord.get(APP_USER.USERNAME),
						authorRecord.get(APP_USER.BIO), authorRecord.get(APP_USER.IMAGE),
						following);
			}

			boolean favorited = dsl.selectCount().from(ARTICLE_FAVORITE)
					.where(ARTICLE_FAVORITE.USER_ID.eq(userId)
							.and(ARTICLE_FAVORITE.ARTICLE_ID.eq(articleId)))
					.fetchOne(0, int.class) == 1;
			int favoritesCount = dsl.selectCount().from(ARTICLE_FAVORITE)
					.where(ARTICLE_FAVORITE.ARTICLE_ID.eq(articleId))
					.fetchOne(0, int.class);

			return new Article(record.get(ARTICLE.SLUG), record.get(ARTICLE.TITLE),
					record.get(ARTICLE.DESCRIPTION), record.get(ARTICLE.BODY),
					dsl.select(TAG.NAME).from(TAG).innerJoin(ARTICLE_TAG).onKey()
							.where(ARTICLE_TAG.ARTICLE_ID.eq(articleId))
							.fetchSet(TAG.NAME),
					record.get(ARTICLE.CREATED_AT), record.get(ARTICLE.UPDATED_AT),
					favorited, favoritesCount, author);

		}
		return null;
	}
*/
