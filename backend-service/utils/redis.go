package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectToRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "secretpass", // password set
		DB:       0,            // use default DB
	})

	res := rdb.Ping(context.Background())
	if res.Err() != nil {
		log.Println(res.Err())
		return
	}
	log.Println("Connected to redis cache successfully")
	Redis = rdb
}
