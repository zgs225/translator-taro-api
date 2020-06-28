package services

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var _redisClient *redis.Client

// MustGetRedisClient 获取 Redis 客户端
func MustGetRedisClient() *redis.Client {
	if _redisClient == nil {
		_redisClient = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
			DB:           viper.GetInt("redis.db"),
			Password:     viper.GetString("redis.password"),
			PoolSize:     50,
			MinIdleConns: 25,
		})

		if err := _redisClient.Ping().Err(); err != nil {
			panic(err)
		}
	}
	return _redisClient
}
