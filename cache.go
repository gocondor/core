package core

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx context.Context

type Cache struct {
	redis *redis.Client
}

func NewCache(cacheConfig CacheConfig) *Cache {
	ctx = context.Background()
	dbStr := os.Getenv("REDIS_DB")
	db64, err := strconv.ParseInt(dbStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing redis db env var: %v", err))
	}
	db := int(db64)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       db,                          // use default DB
	})

	_, err = rdb.Ping(ctx).Result()
	if cacheConfig.EnableCache && err != nil {
		panic(fmt.Sprintf("problem connecting to redis cache, (if it's not needed you can disable it in config/cache.go): %v", err))
	}

	return &Cache{
		redis: rdb,
	}
}

func (c *Cache) Set(key string, value string) error {
	err := c.redis.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) SetWithExpiration(key string, value string, expiration time.Duration) error {
	err := c.redis.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	result, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Cache) Delete(key string) error {
	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
