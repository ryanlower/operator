package main

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
)

// Creates (and returns) RedisStore connected to
// database 1 on localhost
// Calls FLUSHDB to delete all keys from database 1
// so we start with a fresh db
func setupRedisStore() *RedisStore {
	conf := new(Config)
	conf.Redis.Address = "localhost:6379"
	conf.Redis.Password = "password"

	store := &RedisStore{config: conf}
	store.connect()
	store.connection.Do("SELECT", "1")
	store.connection.Do("FLUSHDB")

	return store
}

func TestSetNewToken(t *testing.T) {
	store := setupRedisStore()

	// store.Set doesn't error
	err := store.Set("redis", "http://redis.io")
	assert.Nil(t, err)

	// key has been set in redis
	url, _ := redis.String(store.connection.Do("GET", "redis"))
	assert.Equal(t, url, "http://redis.io")
}

func TestSetExistingToken(t *testing.T) {
	store := setupRedisStore()

	store.connection.Do("SET", "redis", "http://redis.io")

	err := store.Set("redis", "http://redisconference.com")
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "Token already used")
	}

	// key has not been overwritten in redis
	url, _ := redis.String(store.connection.Do("GET", "redis"))
	assert.Equal(t, url, "http://redis.io")
}

func TestGetTokenSet(t *testing.T) {
	store := setupRedisStore()

	store.connection.Do("SET", "redis", "http://redis.io")

	// store.Get returns url and doesn't error
	reply, err := store.Get("redis")
	assert.Nil(t, err)
	assert.Equal(t, reply, "http://redis.io")
}

func TestGetTokenMissin(t *testing.T) {
	store := setupRedisStore()

	_, err := store.Get("redis")
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "Token not found")
	}
}
