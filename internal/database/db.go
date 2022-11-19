package database

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"realworldgo.rasc.ch/internal/config"
	"time"
)

func New(cfg config.Config) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", cfg.DB.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	connMaxIdleTime, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(connMaxIdleTime)

	connMaxLifetime, err := time.ParseDuration(cfg.DB.MaxLifetime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
