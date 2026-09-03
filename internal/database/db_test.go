package database

import (
	"net/url"
	"testing"

	"realworldgo.rasc.ch/internal/config"
)

func TestDSNEscapesCredentials(t *testing.T) {
	var cfg config.Config
	cfg.DB.User = "api@example.com"
	cfg.DB.Password = "p@ss:word"
	cfg.DB.Host = "db.internal:5432"
	cfg.DB.Database = "real world"

	parsed, err := url.Parse(DSN(cfg))
	if err != nil {
		t.Fatal(err)
	}
	password, passwordSet := parsed.User.Password()
	if got := parsed.User.Username(); got != cfg.DB.User {
		t.Fatalf("username = %q, want %q", got, cfg.DB.User)
	}
	if !passwordSet || password != cfg.DB.Password {
		t.Fatalf("password = %q, %v, want %q, true", password, passwordSet, cfg.DB.Password)
	}
	if parsed.Host != cfg.DB.Host {
		t.Fatalf("host = %q, want %q", parsed.Host, cfg.DB.Host)
	}
	if parsed.Path != "/"+cfg.DB.Database {
		t.Fatalf("path = %q, want %q", parsed.Path, "/"+cfg.DB.Database)
	}
}
