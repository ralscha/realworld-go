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

// Comment is an object representing the database table.
type Comment struct {
	ID        int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	Body      string    `boil:"body" json:"body" toml:"body" yaml:"body"`
	ArticleID int64     `boil:"article_id" json:"article_id" toml:"article_id" yaml:"article_id"`
	UserID    int64     `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *commentR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L commentL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CommentColumns = struct {
	ID        string
	Body      string
	ArticleID string
	UserID    string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Body:      "body",
	ArticleID: "article_id",
	UserID:    "user_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var CommentTableColumns = struct {
	ID        string
	Body      string
	ArticleID string
	UserID    string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "comment.id",
	Body:      "comment.body",
	ArticleID: "comment.article_id",
	UserID:    "comment.user_id",
	CreatedAt: "comment.created_at",
	UpdatedAt: "comment.updated_at",
}

// Generated where

var CommentWhere = struct {
	ID        whereHelperint64
	Body      whereHelperstring
	ArticleID whereHelperint64
	UserID    whereHelperint64
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperint64{field: "\"comment\".\"id\""},
	Body:      whereHelperstring{field: "\"comment\".\"body\""},
	ArticleID: whereHelperint64{field: "\"comment\".\"article_id\""},
	UserID:    whereHelperint64{field: "\"comment\".\"user_id\""},
	CreatedAt: whereHelpertime_Time{field: "\"comment\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"comment\".\"updated_at\""},
}

// CommentRels is where relationship names are stored.
var CommentRels = struct {
	Article string
	User    string
}{
	Article: "Article",
	User:    "User",
}

// commentR is where relationships are stored.
type commentR struct {
	Article *Article `boil:"Article" json:"Article" toml:"Article" yaml:"Article"`
	User    *AppUser `boil:"User" json:"User" toml:"User" yaml:"User"`
}

// NewStruct creates a new relationship struct
func (*commentR) NewStruct() *commentR {
	return &commentR{}
}

func (r *commentR) GetArticle() *Article {
	if r == nil {
		return nil
	}
	return r.Article
}

func (r *commentR) GetUser() *AppUser {
	if r == nil {
		return nil
	}
	return r.User
}

// commentL is where Load methods for each relationship are stored.
type commentL struct{}

var (
	commentAllColumns            = []string{"id", "body", "article_id", "user_id", "created_at", "updated_at"}
	commentColumnsWithoutDefault = []string{"body", "article_id", "user_id", "created_at", "updated_at"}
	commentColumnsWithDefault    = []string{"id"}
	commentPrimaryKeyColumns     = []string{"id"}
	commentGeneratedColumns      = []string{}
)

type (
	// CommentSlice is an alias for a slice of pointers to Comment.
	// This should almost always be used instead of []Comment.
	CommentSlice []*Comment

	commentQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	commentType                 = reflect.TypeOf(&Comment{})
	commentMapping              = queries.MakeStructMapping(commentType)
	commentPrimaryKeyMapping, _ = queries.BindMapping(commentType, commentMapping, commentPrimaryKeyColumns)
	commentInsertCacheMut       sync.RWMutex
	commentInsertCache          = make(map[string]insertCache)
	commentUpdateCacheMut       sync.RWMutex
	commentUpdateCache          = make(map[string]updateCache)
	commentUpsertCacheMut       sync.RWMutex
	commentUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single comment record from the query.
func (q commentQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Comment, error) {
	o := &Comment{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for comment")
	}

	return o, nil
}

// All returns all Comment records from the query.
func (q commentQuery) All(ctx context.Context, exec boil.ContextExecutor) (CommentSlice, error) {
	var o []*Comment

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Comment slice")
	}

	return o, nil
}

// Count returns the count of all Comment records in the query.
func (q commentQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count comment rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q commentQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if comment exists")
	}

	return count > 0, nil
}

// Article pointed to by the foreign key.
func (o *Comment) Article(mods ...qm.QueryMod) articleQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ArticleID),
	}

	queryMods = append(queryMods, mods...)

	return Articles(queryMods...)
}

// User pointed to by the foreign key.
func (o *Comment) User(mods ...qm.QueryMod) appUserQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	return AppUsers(queryMods...)
}

// LoadArticle allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (commentL) LoadArticle(ctx context.Context, e boil.ContextExecutor, singular bool, maybeComment interface{}, mods queries.Applicator) error {
	var slice []*Comment
	var object *Comment

	if singular {
		var ok bool
		object, ok = maybeComment.(*Comment)
		if !ok {
			object = new(Comment)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeComment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeComment))
			}
		}
	} else {
		s, ok := maybeComment.(*[]*Comment)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeComment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeComment))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &commentR{}
		}
		args = append(args, object.ArticleID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &commentR{}
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
		foreign.R.Comments = append(foreign.R.Comments, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ArticleID == foreign.ID {
				local.R.Article = foreign
				if foreign.R == nil {
					foreign.R = &articleR{}
				}
				foreign.R.Comments = append(foreign.R.Comments, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (commentL) LoadUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybeComment interface{}, mods queries.Applicator) error {
	var slice []*Comment
	var object *Comment

	if singular {
		var ok bool
		object, ok = maybeComment.(*Comment)
		if !ok {
			object = new(Comment)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeComment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeComment))
			}
		}
	} else {
		s, ok := maybeComment.(*[]*Comment)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeComment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeComment))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &commentR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &commentR{}
			}

			for _, a := range args {
				if a == obj.UserID {
					continue Outer
				}
			}

			args = append(args, obj.UserID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`app_user`),
		qm.WhereIn(`app_user.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load AppUser")
	}

	var resultSlice []*AppUser
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice AppUser")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for app_user")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for app_user")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &appUserR{}
		}
		foreign.R.UserComments = append(foreign.R.UserComments, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &appUserR{}
				}
				foreign.R.UserComments = append(foreign.R.UserComments, local)
				break
			}
		}
	}

	return nil
}

// SetArticle of the comment to the related item.
// Sets o.R.Article to related.
// Adds o to related.R.Comments.
func (o *Comment) SetArticle(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Article) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"comment\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"article_id"}),
		strmangle.WhereClause("\"", "\"", 2, commentPrimaryKeyColumns),
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
		o.R = &commentR{
			Article: related,
		}
	} else {
		o.R.Article = related
	}

	if related.R == nil {
		related.R = &articleR{
			Comments: CommentSlice{o},
		}
	} else {
		related.R.Comments = append(related.R.Comments, o)
	}

	return nil
}

// SetUser of the comment to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserComments.
func (o *Comment) SetUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *AppUser) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"comment\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, commentPrimaryKeyColumns),
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

	o.UserID = related.ID
	if o.R == nil {
		o.R = &commentR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &appUserR{
			UserComments: CommentSlice{o},
		}
	} else {
		related.R.UserComments = append(related.R.UserComments, o)
	}

	return nil
}

// Comments retrieves all the records using an executor.
func Comments(mods ...qm.QueryMod) commentQuery {
	mods = append(mods, qm.From("\"comment\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"comment\".*"})
	}

	return commentQuery{q}
}

// FindComment retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindComment(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Comment, error) {
	commentObj := &Comment{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"comment\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, commentObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from comment")
	}

	return commentObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Comment) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no comment provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(commentColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	commentInsertCacheMut.RLock()
	cache, cached := commentInsertCache[key]
	commentInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			commentAllColumns,
			commentColumnsWithDefault,
			commentColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(commentType, commentMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(commentType, commentMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"comment\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"comment\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into comment")
	}

	if !cached {
		commentInsertCacheMut.Lock()
		commentInsertCache[key] = cache
		commentInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Comment.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Comment) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	commentUpdateCacheMut.RLock()
	cache, cached := commentUpdateCache[key]
	commentUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			commentAllColumns,
			commentPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("models: unable to update comment, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"comment\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, commentPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(commentType, commentMapping, append(wl, commentPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update comment row")
	}

	if !cached {
		commentUpdateCacheMut.Lock()
		commentUpdateCache[key] = cache
		commentUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q commentQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for comment")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CommentSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), commentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"comment\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, commentPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in comment slice")
	}

	return nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Comment) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no comment provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(commentColumnsWithDefault, o)

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

	commentUpsertCacheMut.RLock()
	cache, cached := commentUpsertCache[key]
	commentUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			commentAllColumns,
			commentColumnsWithDefault,
			commentColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			commentAllColumns,
			commentPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert comment, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(commentPrimaryKeyColumns))
			copy(conflict, commentPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"comment\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(commentType, commentMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(commentType, commentMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert comment")
	}

	if !cached {
		commentUpsertCacheMut.Lock()
		commentUpsertCache[key] = cache
		commentUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Comment record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Comment) Delete(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil {
		return errors.New("models: no Comment provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), commentPrimaryKeyMapping)
	sql := "DELETE FROM \"comment\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from comment")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q commentQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if q.Query == nil {
		return errors.New("models: no commentQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from comment")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CommentSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), commentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"comment\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, commentPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from comment slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Comment) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindComment(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CommentSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CommentSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), commentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"comment\".* FROM \"comment\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, commentPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in CommentSlice")
	}

	*o = slice

	return nil
}

// CommentExists checks if the Comment row exists.
func CommentExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"comment\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if comment exists")
	}

	return exists, nil
}
