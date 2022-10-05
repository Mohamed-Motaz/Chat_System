package RedisCache

import (
	logger "Server/Logger"
	utils "Server/Utils"
	"context"
	"fmt"
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

	if !utils.IN_DOCKER {
		//just for testing
		go cache.debug()
	}
	return cache
}

func (cache *Cache) Incr(key string) (int64, error) {
	return cache.client.Incr(cache.ctx, key).Result()
}

func (cache *Cache) Decr(key string) (int64, error) {
	return cache.client.Decr(cache.ctx, key).Result()
}

func (cache *Cache) Get(key string) (string, error) {
	return cache.client.Get(cache.ctx, key).Result()
}

func (cache *Cache) Set(key string, value int, ttl time.Duration) error {
	return cache.client.Set(cache.ctx, key, value, ttl).Err()
}

//for setting the key for the number of a specific chat, format is {APP_TOKEN}.CHAT
func (cache *Cache) MakeChatCtrForAppCacheKey(appToken string) string {
	return fmt.Sprintf("%s.CHAT", appToken)
}

//for setting the key for the number of a specific message, format is {APP_TOKEN}.CHAT.{CHAT_NUM}
func (cache *Cache) MakeMessageCtrForChatCacheKey(appToken string, chatNum int) string {
	return fmt.Sprintf("%s.%d", cache.MakeChatCtrForAppCacheKey(appToken), chatNum)
}

//for setting the key for the number of a specific message, format is {APP_TOKEN}.CHAT.{CHAT_NUM}.CHAT_ID
func (cache *Cache) MakeChatIdForChatNumCacheKey(appToken string, chatNum int) string {
	return fmt.Sprintf("%s.%d.CHAT_ID", cache.MakeChatCtrForAppCacheKey(appToken), chatNum)
}

func (cache *Cache) debug() {
	for {
		fmt.Println("About to print all redis keys and vals\n")
		iter := cache.client.Scan(cache.ctx, 0, "*", 0).Iterator()
		for iter.Next(cache.ctx) {
			key := iter.Val()
			val, _ := cache.Get(key)
			fmt.Printf("%v -- %v\n", key, val)
		}
		if err := iter.Err(); err != nil {
			panic(err)
		}
		fmt.Println("\nDone printing all redis keys and vals\n")

		time.Sleep(10 * time.Second)
	}

}
