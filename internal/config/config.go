package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Config struct {
	Environment Environment
	DB          struct {
		User         string
		Password     string
		Host         string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
		MaxLifetime  string
	}
	HTTP struct {
		Port                  string
		ReadTimeoutInSeconds  int64
		WriteTimeoutInSeconds int64
		IdleTimeoutInSeconds  int64
	}
	Argon2 struct {
		Memory      uint32
		Iterations  uint32
		Parallelism uint8
		SaltLength  uint32
		KeyLength   uint32
	}
}

func newViper() *viper.Viper {
	v := viper.New()
	v.SetDefault("environment", Production)
	v.SetDefault("http.readTimeoutInSeconds", 30)
	v.SetDefault("http.writeTimeoutInSeconds", 30)
	v.SetDefault("http.idleTimeoutInSeconds", 120)
	v.SetDefault("db.maxOpenConns", 4)
	v.SetDefault("db.maxIdleConns", 2)
	v.SetDefault("db.maxIdleTime", "15m")
	v.SetDefault("db.maxLifetime", "2h")
	v.SetDefault("argon2.memory", 1<<17)
	v.SetDefault("argon2.iterations", 20)
	v.SetDefault("argon2.parallelism", 8)
	v.SetDefault("argon2.saltLength", 16)
	v.SetDefault("argon2.keyLength", 32)
	v.SetEnvPrefix("REALWORLD")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	return v
}

func LoadConfig() (Config, error) {
	var cfg Config

	v := newViper()
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}
	if cfg.Environment != Development && cfg.Environment != Production {
		return cfg, fmt.Errorf("unsupported environment %q", cfg.Environment)
	}

	return cfg, nil
}
