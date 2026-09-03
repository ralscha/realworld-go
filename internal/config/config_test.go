package config

import (
	"strings"
	"testing"
)

func TestEnvironmentOverridesDottedConfigKey(t *testing.T) {
	t.Setenv("REALWORLD_DB_HOST", "db.internal:5432")

	v := newViper()
	v.SetConfigType("env")
	if err := v.ReadConfig(strings.NewReader("db.host=localhost:5432")); err != nil {
		t.Fatal(err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.DB.Host != "db.internal:5432" {
		t.Fatalf("DB.Host = %q, want environment override", cfg.DB.Host)
	}
}
