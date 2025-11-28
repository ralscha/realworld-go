// Package main is the entry point for the RealWorld API server.
package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"realworldgo.rasc.ch/internal/config"
	"realworldgo.rasc.ch/internal/database"
	"realworldgo.rasc.ch/internal/scsheader"
	"realworldgo.rasc.ch/internal/version"
)

type application struct {
	config         *config.Config
	database       *sql.DB
	sessionManager *scsheader.HeaderSession
}

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("reading config failed %v\n", err)
		return 1
	}

	var logger *slog.Logger

	switch cfg.Environment {
	case config.Development:
		boil.DebugMode = true
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case config.Production:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)

	db, err := database.New(cfg)
	if err != nil {
		log.Printf("opening database connection failed: %v", err)
		return 1
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	sm := scsheader.HeaderSession{SessionManager: scs.New()}
	sm.Store = postgresstore.NewWithCleanupInterval(db, 30*time.Minute)
	sm.Lifetime = 24 * time.Hour

	err = initAuth(cfg)
	if err != nil {
		log.Printf("init auth failed: %v", err)
		return 1
	}

	app := &application{
		config:         &cfg,
		database:       db,
		sessionManager: &sm,
	}

	slog.Info("starting server", "addr", app.config.HTTP.Port, "version", version.Get())

	err = app.serve()
	if err != nil {
		slog.Error("http serve failed", "error", err)
		return 1
	}

	slog.Info("server stopped")
	return 0
}
