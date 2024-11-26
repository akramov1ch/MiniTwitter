package redis

import (
	"comment-service/config"
	"log"

	"github.com/go-redis/redis/v8"			
)

func NewRedisClient() *redis.Client {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: conf.REDIS_HOST + ":" + conf.REDIS_PORT,
	})

	log.Println("Connected to Redis")
	return client
}
