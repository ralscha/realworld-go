package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type failingCommitDriver struct{}

func (failingCommitDriver) Open(string) (driver.Conn, error) {
	return failingCommitConn{}, nil
}

type failingCommitConn struct{}

func (failingCommitConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (failingCommitConn) Close() error { return nil }

func (failingCommitConn) Begin() (driver.Tx, error) { return failingCommitTx{}, nil }

func (failingCommitConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return failingCommitTx{}, nil
}

type failingCommitTx struct{}

func (failingCommitTx) Commit() error   { return errors.New("commit failed") }
func (failingCommitTx) Rollback() error { return nil }

func TestRWTransactionDoesNotPublishResponseBeforeCommit(t *testing.T) {
	const driverName = "failing-commit"
	sql.Register(driverName, failingCommitDriver{})
	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	app := application{database: db}
	handler := app.rwTransaction(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("success"))
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/", nil))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if strings.Contains(recorder.Body.String(), "success") {
		t.Fatalf("uncommitted success response was published: %q", recorder.Body.String())
	}
}
