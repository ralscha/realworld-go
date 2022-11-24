package main

import (
	"database/sql"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/exp/slog"
	"log"
	"math/rand"
	"os"
	"realworldgo.rasc.ch/internal/config"
	"realworldgo.rasc.ch/internal/database"
	"realworldgo.rasc.ch/internal/scsheader"
	"realworldgo.rasc.ch/internal/version"
	"time"
)

type application struct {
	config         *config.Config
	database       *sql.DB
	sessionManager *scsheader.HeaderSession
}

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("reading config failed %v\n", err)
	}

	var logger *slog.Logger

	switch cfg.Environment {
	case config.Development:
		boil.DebugMode = true
		logger = slog.New(slog.NewTextHandler(os.Stdout))
	case config.Production:
		logger = slog.New(slog.NewJSONHandler(os.Stdout))
	}

	slog.SetDefault(logger)

	db, err := database.New(cfg)
	if err != nil {
		logger.Error("opening database connection failed", err)
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	sm := scsheader.HeaderSession{SessionManager: scs.New()}
	sm.Store = postgresstore.NewWithCleanupInterval(db, 30*time.Minute)
	sm.Lifetime = 24 * time.Hour

	err = initAuth(cfg)
	if err != nil {
		logger.Error("init auth failed", err)
		os.Exit(1)
	}

	app := &application{
		config:         &cfg,
		database:       db,
		sessionManager: &sm,
	}

	logger.Info("starting server", "addr", app.config.HTTP.Port, "version", version.Get())

	err = app.serve()
	if err != nil {
		logger.Error("http serve failed", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
