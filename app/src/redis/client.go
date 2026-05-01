package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379"
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("[Redis] Invalid URL: %v", err)
	}

	Client = redis.NewClient(opts)

	if err := Client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("[Redis] Could not connect: %v", err)
	}

	log.Println("[Redis] Connected")
}
