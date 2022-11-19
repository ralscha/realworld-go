package main

import (
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"log"
	"os"
	"realworldgo.rasc.ch/internal/config"
	"realworldgo.rasc.ch/migrations"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 {
		flags.Usage()
		return
	}
	command := args[0]

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("reading config failed", err)
	}

	db, err := goose.OpenDBWithDriver("sqlite3", cfg.DB.Dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	goose.SetBaseFS(migrations.EmbeddedFiles)

	if err := goose.Run(command, db, ".", arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
