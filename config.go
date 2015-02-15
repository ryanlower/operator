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
		port string
	}
}

// Load config from environment
func (c *Config) load() {
	c.port = envOrDefault("PORT", "3000")
	c.auth.password = os.Getenv("AUTH_PASSWORD")
	c.redis.port = envOrDefault("REDIS_PORT", "6379")
}

func envOrDefault(key string, defaultValue string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return defaultValue
}
