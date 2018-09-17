package model

import (
	"github.com/go-redis/redis"
	"fmt"
	"time"
)

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password: 	"",
		DB:			0,
	})

	go pingRedis()
}

func pingRedis() {
	for {
		if _, err := RedisClient.Ping().Result(); err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(time.Second)
	}
}