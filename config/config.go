package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"log"`
		PG      `yaml:"postgres"`
		Service `yaml:"service"`
		LinkGen `yaml:"generator"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port         string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		WriteTimeout time.Duration `env-required:"true" yaml:"write_timeout" env:"WRITE_TIMEOUT"`
		ReadTimeout  time.Duration `env-required:"true" yaml:"read_timeout" env:"READ_TIMEOUT"`
	}

	Log struct {
		Level string `yaml:"log_level"`
	}

	PG struct {
		Name     string `env:"DB_NAME"`
		User     string `env:"DB_USER"`
		Port     int    `env:"DB_PORT"`
		Password string `env:"DB_PASSWORD"`
		Host     string `env:"DB_HOST"`
		PoolMax  int    `yaml:"pool_max"`
	}

	Service struct {
		Host                  string `yaml:"host"`
		Port                  int    `yaml:"port"`
		RecalculationInterval int    `yaml:"interval"`
	}

	LinkGen struct {
		Alphabet string `yaml:"alphabet"`
		Length   int    `yaml:"length"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yaml", cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}

	return cfg, nil
}
