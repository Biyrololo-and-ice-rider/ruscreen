package config

import "github.com/kelseyhightower/envconfig"

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port      string `envconfig:"PORT"       default:"8080"`
	RedisURL  string `envconfig:"REDIS_URL"  default:"redis://localhost:6379"`
	JWTSecret string `envconfig:"JWT_SECRET" default:"change-me-in-development"`
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
