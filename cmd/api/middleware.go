package main

import (
	"context"
	"database/sql"
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

func (app *application) rwTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := app.database.BeginTx(r.Context(), nil)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), transactionKey, tx)
		next.ServeHTTP(w, r.WithContext(ctx))

		if err := tx.Commit(); err != nil {
			response.InternalServerError(w, err)
			return
		}
	})
}

func (app *application) readonlyTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := app.database.BeginTx(r.Context(), &sql.TxOptions{ReadOnly: true})
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), transactionKey, tx)
		next.ServeHTTP(w, r.WithContext(ctx))

		if err := tx.Rollback(); err != nil {
			response.InternalServerError(w, err)
			return
		}
	})
}
