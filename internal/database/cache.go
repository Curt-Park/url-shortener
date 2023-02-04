package database

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	urlCache *cache.Cache
}

func NewCache(redisServer string, refreshTokenExpire int, size int) *Cache {
	urlCache := cache.New(&cache.Options{
		Redis:      redis.NewClient(&redis.Options{Addr: redisServer}),
		LocalCache: cache.NewTinyLFU(size, (time.Duration)(refreshTokenExpire)*time.Minute),
	})
	return &Cache{urlCache: urlCache}
}

func (c *Cache) Set(key string, value interface{}) {
	item := cache.Item{
		Ctx:   context.TODO(),
		Key:   key,
		Value: value,
	}
	if err := c.urlCache.Set(&item); err != nil {
		log.Panicf("failed to set the item in Cache server: %v", err)
	}
}

func (c *Cache) Get(key string) (string, bool) {
	ctx := context.TODO()
	var value string
	if err := c.urlCache.Get(ctx, key, &value); err != nil {
		log.Printf("failed to get the item from Cache server: %v", err)
		return value, false
	}
	return value, true
}

func (c *Cache) Delete(key string) {
	if err := c.urlCache.Delete(context.TODO(), key); err != nil {
		log.Panicf("failed to delete the item from Cache server: %v", err)
	}
}
