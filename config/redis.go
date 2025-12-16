package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password by default
		DB:       0,  // default DB
	})

	// Test connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Printf("⚠️  Redis connection failed: %v (caching will be disabled)", err)
		RedisClient = nil
		return
	}

	log.Println("✅ Redis connected successfully")
}

// Helper functions for caching
func SetCache(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("redis not available")
	}
	return RedisClient.Set(Ctx, key, value, expiration).Err()
}

func GetCache(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("redis not available")
	}
	return RedisClient.Get(Ctx, key).Result()
}

func DeleteCache(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("redis not available")
	}
	return RedisClient.Del(Ctx, key).Err()
}

func DeleteCachePattern(pattern string) error {
	if RedisClient == nil {
		return fmt.Errorf("redis not available")
	}

	iter := RedisClient.Scan(Ctx, 0, pattern, 0).Iterator()
	for iter.Next(Ctx) {
		err := RedisClient.Del(Ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
