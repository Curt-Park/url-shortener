package internal

import (
	"context"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type DatabaseOperations interface {
	Set(key string, value interface{})
	Get(key string) (string, bool)
	Delete(key string)
}

type Database struct {
	urlDB    *redis.Client
	urlCache *cache.Cache
}

func NewDatabase(redisServer string, password string, db int) *Database {
	urlDB := redis.NewClient(&redis.Options{
		Addr:     redisServer,
		Password: password,
		DB:       db,
	})
	urlCache := cache.New(time.Hour, 2*time.Hour)
	return &Database{urlDB: urlDB, urlCache: urlCache}
}

func (db *Database) Set(key string, value interface{}) {
	if err := db.urlDB.Set(ctx, key, value, 0).Err(); err != nil {
		log.Panicf("failed to set the item in the database: %v", err)
	}
	db.urlCache.Set(key, value, cache.DefaultExpiration)
}

func (db *Database) Get(key string) (string, bool) {
	// Get the value from the cache.
	value, exist := db.urlCache.Get(key)
	if exist {
		return value.(string), true
	}

	// Find the value from the db.
	var err error
	value, err = db.urlDB.Get(ctx, key).Result()
	// Could not find.
	if err != nil {
		return value.(string), false
	}

	// Set cache.
	db.urlCache.Set(key, value, cache.DefaultExpiration)
	return value.(string), true
}

func (db *Database) Delete(key string) {
	db.urlCache.Delete(key)
	if err := db.urlDB.Del(ctx, key).Err(); err != nil {
		log.Panicf("failed to delete the item from the database: %v", err)
	}
}
