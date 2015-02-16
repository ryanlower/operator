package main

import (
	"os"
)

// Config ...
type Config struct {
	port string
	auth struct {
		password string // optional for HTTP basic auth for link creation
	}
	redis struct {
		address  string
		password string // optional for redis AUTH
	}
}

// Load config from environment
func (c *Config) load() {
	c.port = envOrDefault("PORT", "3000")
	c.auth.password = os.Getenv("AUTH_PASSWORD")
	c.redis.address = envOrDefault("REDIS_ADDRESS", "localhost:6379")
	c.redis.password = os.Getenv("REDIS_PASSWORD")
}

func envOrDefault(key string, defaultValue string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return defaultValue
}
