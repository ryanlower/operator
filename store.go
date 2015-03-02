package main

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

// Store is an interface ...
type Store interface {
	Get(token string) (string, error)
	Set(token, url string) error
}

// RedisStore is a redis backed Store
type RedisStore struct {
	Store
	config     *Config
	connection redis.Conn
}

// Set ...
func (s *RedisStore) Set(token, url string) error {
	if s.connection == nil {
		s.connect()
	}

	reply, _ := redis.Int(s.connection.Do("SETNX", token, url))
	if reply != 1 {
		return errors.New("Token already used")
	}
	return nil
}

// Get ...
func (s *RedisStore) Get(token string) (string, error) {
	if s.connection == nil {
		s.connect()
	}

	reply, _ := redis.String(s.connection.Do("GET", token))
	if reply == "" {
		return reply, errors.New("Token not found")
	}
	return reply, nil
}

func (s *RedisStore) connect() {
	conn, err := redis.Dial("tcp", s.config.Redis.Address)
	if err != nil {
		panic(err) // Can't do much without a redis connection
	}

	// AUTH if config specifies redis passwoed
	if s.config.Redis.Password != "" {
		conn.Do("AUTH", s.config.Redis.Password)
	}

	s.connection = conn
}
