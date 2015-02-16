package main

import (
	"os"
)

// Config ...
type Config struct {
	port string
	auth struct {
		password string
	}
	redis struct {
		address string
	}
}

// Load config from environment
func (c *Config) load() {
	c.port = envOrDefault("PORT", "3000")
	c.auth.password = os.Getenv("AUTH_PASSWORD")
	c.redis.address = envOrDefault("REDIS_ADDRESS", "localhost:6379")
}

func envOrDefault(key string, defaultValue string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return defaultValue
}
