package RedisCache

import (
	logger "Server/Logger"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func New(cacheAddr string) *Cache {
	cache := &Cache{
		client: redis.NewClient(&redis.Options{
			Addr:     cacheAddr,
			Password: "",
			DB:       0, //default db
		}),
		ctx: context.Background(),
	}

	_, err := cache.client.Ping(cache.ctx).Result()
	if err != nil {
		logger.FailOnError(logger.SERVER, logger.ESSENTIAL, "Unable to connect to caching layer with error %v", err)
	} else {
		logger.LogInfo(logger.SERVER, logger.ESSENTIAL, "Successfully connected to caching layer")
	}

	return cache
}

func (cache *Cache) Add(key string, value string, ttl time.Duration) error {
	return cache.client.Set(cache.ctx, key, value, ttl).Err()
}

func (cache *Cache) Get(key string) (string, error) {
	return cache.client.Get(cache.ctx, key).Result()
}

func (cache *Cache) GetBytes(key string) ([]byte, error) {
	return cache.client.Get(cache.ctx, key).Bytes()
}
