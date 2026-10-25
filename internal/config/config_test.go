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

func TestTrustedProxyCIDRsFromEnvironment(t *testing.T) {
	t.Setenv("REALWORLD_HTTP_TRUSTEDPROXYCIDRS", "10.0.0.2/32,::1/128")
	var cfg Config
	if err := newViper().Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if len(cfg.HTTP.TrustedProxyCIDRs) != 2 || cfg.HTTP.TrustedProxyCIDRs[0] != "10.0.0.2/32" {
		t.Fatalf("unexpected proxies: %v", cfg.HTTP.TrustedProxyCIDRs)
	}
	if err := cfg.validateTrustedProxies(); err != nil {
		t.Fatal(err)
	}
	cfg.HTTP.TrustedProxyCIDRs = []string{"not-a-cidr"}
	if err := cfg.validateTrustedProxies(); err == nil {
		t.Fatal("invalid CIDR accepted")
	}
}
