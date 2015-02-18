package main

import (
	"github.com/ryanlower/setting"
)

// Config ...
type Config struct {
	Port string `env:"PORT" default:"3000"`
	Auth struct {
		// optional for HTTP basic auth for link creation
		Password string `env:"AUTH_PASSWORD"`
	}
	Redis struct {
		Address  string `env:"REDIS_ADDRESS" default:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD"` // optional for redis AUTH
	}
}

// Load config from environment
func (c *Config) load() {
	setting.Load(c)
}
