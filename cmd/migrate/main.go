// Package main is the entry point for the database migration tool.
package main

import (
	"context"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"realworldgo.rasc.ch/internal/config"
	"realworldgo.rasc.ch/migrations"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	os.Exit(run())
}

func run() int {
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 {
		flags.Usage()
		return 0
	}
	command := args[0]

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("reading config failed", err)
		return 1
	}

	dsn := "postgres://" + cfg.DB.User + ":" + cfg.DB.Password + "@" + cfg.DB.Connection + "/" + cfg.DB.Database
	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Printf("goose: failed to open DB: %v\n", err)
		return 1
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	goose.SetBaseFS(migrations.EmbeddedFiles)

	if err := goose.RunContext(context.Background(), command, db, ".", arguments...); err != nil {
		log.Printf("goose %v: %v", command, err)
		return 1
	}
	return 0
}
