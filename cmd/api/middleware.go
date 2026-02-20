package main

import (
	"context"
	"database/sql"
	"errors"
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
			response.Forbidden(w)
		}
	})
}

type contextKey string

const (
	transactionKey contextKey = "transaction"
)

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (sr *statusRecorder) WriteHeader(statusCode int) {
	sr.status = statusCode
	sr.wroteHeader = true
	sr.ResponseWriter.WriteHeader(statusCode)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if !sr.wroteHeader {
		sr.WriteHeader(http.StatusOK)
	}

	return sr.ResponseWriter.Write(b)
}

func (app *application) rwTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := app.database.BeginTx(r.Context(), nil)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		committed := false
		defer func() {
			if !committed {
				_ = tx.Rollback()
			}
		}()

		ctx := context.WithValue(r.Context(), transactionKey, tx)
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r.WithContext(ctx))

		if recorder.status >= http.StatusBadRequest {
			return
		}

		if err := tx.Commit(); err != nil {
			if !recorder.wroteHeader {
				response.InternalServerError(w, err)
			}
			return
		}

		committed = true
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
				response.InternalServerError(w, err)
			}
		}()

		ctx := context.WithValue(r.Context(), transactionKey, tx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
