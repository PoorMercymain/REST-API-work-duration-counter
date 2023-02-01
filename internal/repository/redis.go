package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func RedisConnect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}
func RedisSet(rdb *redis.Client, key string, value string) {
	var ctx = context.Background()
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
func RedisGet(rdb *redis.Client, key string) string {
	var ctx = context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}
