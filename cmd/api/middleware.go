package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"maps"
	"net/http"

	"realworldgo.rasc.ch/internal/response"
)

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := app.sessionManager.Get(r.Context(), "userID")
		userID, ok := value.(int64)
		if ok && userID > 0 {
			next.ServeHTTP(w, r)
		} else {
			response.Unauthorized(w)
		}
	})
}

type contextKey string

const (
	transactionKey contextKey = "transaction"
)

type bufferedResponse struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func newBufferedResponse() *bufferedResponse {
	return &bufferedResponse{header: make(http.Header)}
}

func (br *bufferedResponse) Header() http.Header {
	return br.header
}

func (br *bufferedResponse) WriteHeader(statusCode int) {
	if br.status == 0 {
		br.status = statusCode
	}
}

func (br *bufferedResponse) Write(b []byte) (int, error) {
	if br.status == 0 {
		br.status = http.StatusOK
	}
	return br.body.Write(b)
}

func (br *bufferedResponse) writeTo(w http.ResponseWriter) {
	maps.Copy(w.Header(), br.header)
	status := br.status
	if status == 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	if _, err := br.body.WriteTo(w); err != nil {
		slog.Error("writing buffered response failed", "error", err)
	}
}

func (app *application) rwTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := app.database.BeginTx(r.Context(), nil)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		finished := false
		defer func() {
			if !finished {
				_ = tx.Rollback()
			}
		}()

		ctx := context.WithValue(r.Context(), transactionKey, tx)
		buffer := newBufferedResponse()
		next.ServeHTTP(buffer, r.WithContext(ctx))

		if buffer.status >= http.StatusBadRequest {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				slog.Error("rolling back transaction failed", "error", err)
			}
			finished = true
			buffer.writeTo(w)
			return
		}

		if err := tx.Commit(); err != nil {
			response.InternalServerError(w, err)
			return
		}

		finished = true
		buffer.writeTo(w)
	})
}

func (app *application) readonlyTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := app.database.BeginTx(r.Context(), &sql.TxOptions{ReadOnly: true})
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		defer func() {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				slog.Error("rolling back read-only transaction failed", "error", err)
			}
		}()

		ctx := context.WithValue(r.Context(), transactionKey, tx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
