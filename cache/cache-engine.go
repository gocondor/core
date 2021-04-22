// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

//redisCTX context for redis
var redisCTX context.Context

// redisTTL expiry time for redis records
var redisTTL time.Duration

// CacheEngine handles the caching operations
type CacheEngine struct {
	redisDB *redis.Client
}

var cEngine *CacheEngine

// New initiates a new caching engine
func New() *CacheEngine {
	redisCTX = context.Background()
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbName, _ := strconv.ParseInt(os.Getenv("REDIS_DB_NAME"), 10, 32)

	cEngine = &CacheEngine{
		redisDB: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       int(dbName),
		}),
	}
	return cEngine
}

// Resolve resolves initiated caching engine
func Resolve() *CacheEngine {
	return cEngine
}

// Set set a key, val pair in the cache
func (c *CacheEngine) Set(key string, val string) (bool, error) {
	status := c.redisDB.Set(redisCTX, key, val, 0)
	if status.Err() != nil {
		return false, status.Err()
	}

	return true, nil
}

// Get retrieves a val from cache by a given key
func (c *CacheEngine) Get(key string) (interface{}, error) {
	val, err := c.redisDB.Get(redisCTX, key).Result()
	if err != nil {
		return false, err
	}
	return val, nil
}

// Delete removes a record from cache by a given key
func (c *CacheEngine) Delete(key string) error {
	status := c.redisDB.Del(redisCTX, key)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}
