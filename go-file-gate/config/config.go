package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config structure for application configuration
type Config struct {
	App      AppConfig      `yaml:"app"`
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
	Log      LogConfig      `yaml:"log"`
	Redis    RedisConfig    `yaml:"redis"`
}

// AppConfig holds general application configurations
type AppConfig struct {
	Name    string `yaml:"name" env:"APP_NAME" env-default:"GoChatApp"`
	Version string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
}

// HTTPConfig holds the configuration for the HTTP server
type HTTPConfig struct {
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

// PostgresConfig holds the configuration for the PostgreSQL database
type PostgresConfig struct {
	PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX" env-default:"10"`
	URL     string `env:"DATABASE_URL" env-required:"true"`
}

// Log holds the configuration for the logger
type LogConfig struct {
	Level string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
}

// NewConfig reads application configuration and returns it
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - ReadConfig: %w", err)
	}
	return cfg, nil
}
