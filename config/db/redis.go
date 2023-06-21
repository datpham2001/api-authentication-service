package db

// import (
// 	"context"
// 	"fmt"
// 	"realworld-authentication/config"

// 	"github.com/redis/go-redis/v9"
// )

// var (
// 	RedisClient *redis.Client
// )

// func ConnectRedis()
// 	RedisClient = redis.NewClient(&redis.Options{
// 		Addr: config.AppConfig.RedisUrl,
// 	})

// 	if _, err := RedisClient.Ping(context.TODO()).Result(); err != nil {
// 		panic(err)
// 	}

// 	err := RedisClient.Set(context.TODO(), "test", "This is response from redis server", 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Redis client connected successfully")
// }
//
