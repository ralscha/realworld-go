// Code generated by SQLBoiler 4.13.0 (https://github.com/aarondl/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/queries/qmhelper"
	"github.com/aarondl/strmangle"
	"github.com/friendsofgo/errors"
)

// ArticleTag is an object representing the database table.
type ArticleTag struct {
	ID        int64 `boil:"id" json:"id" toml:"id" yaml:"id"`
	ArticleID int64 `boil:"article_id" json:"article_id" toml:"article_id" yaml:"article_id"`
	TagID     int64 `boil:"tag_id" json:"tag_id" toml:"tag_id" yaml:"tag_id"`

	R *articleTagR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L articleTagL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ArticleTagColumns = struct {
	ID        string
	ArticleID string
	TagID     string
}{
	ID:        "id",
	ArticleID: "article_id",
	TagID:     "tag_id",
}

var ArticleTagTableColumns = struct {
	ID        string
	ArticleID string
	TagID     string
}{
	ID:        "article_tag.id",
	ArticleID: "article_tag.article_id",
	TagID:     "article_tag.tag_id",
}

// Generated where

var ArticleTagWhere = struct {
	ID        whereHelperint64
	ArticleID whereHelperint64
	TagID     whereHelperint64
}{
	ID:        whereHelperint64{field: "\"article_tag\".\"id\""},
	ArticleID: whereHelperint64{field: "\"article_tag\".\"article_id\""},
	TagID:     whereHelperint64{field: "\"article_tag\".\"tag_id\""},
}

// ArticleTagRels is where relationship names are stored.
var ArticleTagRels = struct {
	Article string
	Tag     string
}{
	Article: "Article",
	Tag:     "Tag",
}

// articleTagR is where relationships are stored.
type articleTagR struct {
	Article *Article `boil:"Article" json:"Article" toml:"Article" yaml:"Article"`
	Tag     *Tag     `boil:"Tag" json:"Tag" toml:"Tag" yaml:"Tag"`
}

// NewStruct creates a new relationship struct
func (*articleTagR) NewStruct() *articleTagR {
	return &articleTagR{}
}

func (r *articleTagR) GetArticle() *Article {
	if r == nil {
		return nil
	}
	return r.Article
}

func (r *articleTagR) GetTag() *Tag {
	if r == nil {
		return nil
	}
	return r.Tag
}

// articleTagL is where Load methods for each relationship are stored.
type articleTagL struct{}

var (
	articleTagAllColumns            = []string{"id", "article_id", "tag_id"}
	articleTagColumnsWithoutDefault = []string{"article_id", "tag_id"}
	articleTagColumnsWithDefault    = []string{"id"}
	articleTagPrimaryKeyColumns     = []string{"id"}
	articleTagGeneratedColumns      = []string{}
)

type (
	// ArticleTagSlice is an alias for a slice of pointers to ArticleTag.
	// This should almost always be used instead of []ArticleTag.
	ArticleTagSlice []*ArticleTag

	articleTagQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	articleTagType                 = reflect.TypeOf(&ArticleTag{})
	articleTagMapping              = queries.MakeStructMapping(articleTagType)
	articleTagPrimaryKeyMapping, _ = queries.BindMapping(articleTagType, articleTagMapping, articleTagPrimaryKeyColumns)
	articleTagInsertCacheMut       sync.RWMutex
	articleTagInsertCache          = make(map[string]insertCache)
	articleTagUpdateCacheMut       sync.RWMutex
	articleTagUpdateCache          = make(map[string]updateCache)
	articleTagUpsertCacheMut       sync.RWMutex
	articleTagUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single articleTag record from the query.
func (q articleTagQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ArticleTag, error) {
	o := &ArticleTag{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for article_tag")
	}

	return o, nil
}

// All returns all ArticleTag records from the query.
func (q articleTagQuery) All(ctx context.Context, exec boil.ContextExecutor) (ArticleTagSlice, error) {
	var o []*ArticleTag

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ArticleTag slice")
	}

	return o, nil
}

// Count returns the count of all ArticleTag records in the query.
func (q articleTagQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count article_tag rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q articleTagQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if article_tag exists")
	}

	return count > 0, nil
}

// Article pointed to by the foreign key.
func (o *ArticleTag) Article(mods ...qm.QueryMod) articleQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ArticleID),
	}

	queryMods = append(queryMods, mods...)

	return Articles(queryMods...)
}

// Tag pointed to by the foreign key.
func (o *ArticleTag) Tag(mods ...qm.QueryMod) tagQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.TagID),
	}

	queryMods = append(queryMods, mods...)

	return Tags(queryMods...)
}

// LoadArticle allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (articleTagL) LoadArticle(ctx context.Context, e boil.ContextExecutor, singular bool, maybeArticleTag interface{}, mods queries.Applicator) error {
	var slice []*ArticleTag
	var object *ArticleTag

	if singular {
		var ok bool
		object, ok = maybeArticleTag.(*ArticleTag)
		if !ok {
			object = new(ArticleTag)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeArticleTag)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeArticleTag))
			}
		}
	} else {
		s, ok := maybeArticleTag.(*[]*ArticleTag)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeArticleTag)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeArticleTag))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &articleTagR{}
		}
		args = append(args, object.ArticleID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &articleTagR{}
			}

			for _, a := range args {
				if a == obj.ArticleID {
					continue Outer
				}
			}

			args = append(args, obj.ArticleID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`article`),
		qm.WhereIn(`article.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Article")
	}

	var resultSlice []*Article
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Article")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for article")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for article")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Article = foreign
		if foreign.R == nil {
			foreign.R = &articleR{}
		}
		foreign.R.ArticleTags = append(foreign.R.ArticleTags, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ArticleID == foreign.ID {
				local.R.Article = foreign
				if foreign.R == nil {
					foreign.R = &articleR{}
				}
				foreign.R.ArticleTags = append(foreign.R.ArticleTags, local)
				break
			}
		}
	}

	return nil
}

// LoadTag allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (articleTagL) LoadTag(ctx context.Context, e boil.ContextExecutor, singular bool, maybeArticleTag interface{}, mods queries.Applicator) error {
	var slice []*ArticleTag
	var object *ArticleTag

	if singular {
		var ok bool
		object, ok = maybeArticleTag.(*ArticleTag)
		if !ok {
			object = new(ArticleTag)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeArticleTag)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeArticleTag))
			}
		}
	} else {
		s, ok := maybeArticleTag.(*[]*ArticleTag)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeArticleTag)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeArticleTag))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &articleTagR{}
		}
		args = append(args, object.TagID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &articleTagR{}
			}

			for _, a := range args {
				if a == obj.TagID {
					continue Outer
				}
			}

			args = append(args, obj.TagID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`tag`),
		qm.WhereIn(`tag.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Tag")
	}

	var resultSlice []*Tag
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Tag")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for tag")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for tag")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Tag = foreign
		if foreign.R == nil {
			foreign.R = &tagR{}
		}
		foreign.R.ArticleTags = append(foreign.R.ArticleTags, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.TagID == foreign.ID {
				local.R.Tag = foreign
				if foreign.R == nil {
					foreign.R = &tagR{}
				}
				foreign.R.ArticleTags = append(foreign.R.ArticleTags, local)
				break
			}
		}
	}

	return nil
}

// SetArticle of the articleTag to the related item.
// Sets o.R.Article to related.
// Adds o to related.R.ArticleTags.
func (o *ArticleTag) SetArticle(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Article) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"article_tag\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"article_id"}),
		strmangle.WhereClause("\"", "\"", 2, articleTagPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ArticleID = related.ID
	if o.R == nil {
		o.R = &articleTagR{
			Article: related,
		}
	} else {
		o.R.Article = related
	}

	if related.R == nil {
		related.R = &articleR{
			ArticleTags: ArticleTagSlice{o},
		}
	} else {
		related.R.ArticleTags = append(related.R.ArticleTags, o)
	}

	return nil
}

// SetTag of the articleTag to the related item.
// Sets o.R.Tag to related.
// Adds o to related.R.ArticleTags.
func (o *ArticleTag) SetTag(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Tag) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"article_tag\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"tag_id"}),
		strmangle.WhereClause("\"", "\"", 2, articleTagPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.TagID = related.ID
	if o.R == nil {
		o.R = &articleTagR{
			Tag: related,
		}
	} else {
		o.R.Tag = related
	}

	if related.R == nil {
		related.R = &tagR{
			ArticleTags: ArticleTagSlice{o},
		}
	} else {
		related.R.ArticleTags = append(related.R.ArticleTags, o)
	}

	return nil
}

// ArticleTags retrieves all the records using an executor.
func ArticleTags(mods ...qm.QueryMod) articleTagQuery {
	mods = append(mods, qm.From("\"article_tag\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"article_tag\".*"})
	}

	return articleTagQuery{q}
}

// FindArticleTag retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindArticleTag(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*ArticleTag, error) {
	articleTagObj := &ArticleTag{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"article_tag\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, articleTagObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from article_tag")
	}

	return articleTagObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ArticleTag) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no article_tag provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(articleTagColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	articleTagInsertCacheMut.RLock()
	cache, cached := articleTagInsertCache[key]
	articleTagInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			articleTagAllColumns,
			articleTagColumnsWithDefault,
			articleTagColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(articleTagType, articleTagMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(articleTagType, articleTagMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"article_tag\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"article_tag\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into article_tag")
	}

	if !cached {
		articleTagInsertCacheMut.Lock()
		articleTagInsertCache[key] = cache
		articleTagInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ArticleTag.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ArticleTag) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	articleTagUpdateCacheMut.RLock()
	cache, cached := articleTagUpdateCache[key]
	articleTagUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			articleTagAllColumns,
			articleTagPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("models: unable to update article_tag, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"article_tag\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, articleTagPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(articleTagType, articleTagMapping, append(wl, articleTagPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	_, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update article_tag row")
	}

	if !cached {
		articleTagUpdateCacheMut.Lock()
		articleTagUpdateCache[key] = cache
		articleTagUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q articleTagQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for article_tag")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ArticleTagSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), articleTagPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"article_tag\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, articleTagPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in articleTag slice")
	}

	return nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ArticleTag) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no article_tag provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(articleTagColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	articleTagUpsertCacheMut.RLock()
	cache, cached := articleTagUpsertCache[key]
	articleTagUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			articleTagAllColumns,
			articleTagColumnsWithDefault,
			articleTagColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			articleTagAllColumns,
			articleTagPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert article_tag, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(articleTagPrimaryKeyColumns))
			copy(conflict, articleTagPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"article_tag\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(articleTagType, articleTagMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(articleTagType, articleTagMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert article_tag")
	}

	if !cached {
		articleTagUpsertCacheMut.Lock()
		articleTagUpsertCache[key] = cache
		articleTagUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ArticleTag record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ArticleTag) Delete(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil {
		return errors.New("models: no ArticleTag provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), articleTagPrimaryKeyMapping)
	sql := "DELETE FROM \"article_tag\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from article_tag")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q articleTagQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if q.Query == nil {
		return errors.New("models: no articleTagQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from article_tag")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ArticleTagSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), articleTagPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"article_tag\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, articleTagPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from articleTag slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ArticleTag) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindArticleTag(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ArticleTagSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ArticleTagSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), articleTagPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"article_tag\".* FROM \"article_tag\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, articleTagPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ArticleTagSlice")
	}

	*o = slice

	return nil
}

// ArticleTagExists checks if the ArticleTag row exists.
func ArticleTagExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"article_tag\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if article_tag exists")
	}

	return exists, nil
}
