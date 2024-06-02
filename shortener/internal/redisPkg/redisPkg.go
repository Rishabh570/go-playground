package redisPkg

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	Client *redis.Client
}

func GetClient() *Database {
	// Initialize redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Ensure that the connection is properly closed gracefully
	// Can't do this here otherwise conn will close right after this fn finishes
	// defer redisClient.Close()

	// Perform basic diagnostic to check if the connection is working
	ctx := context.Background()
	status, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("Redis connection was refused")
	}
	fmt.Println("Redis status:", status)

	return &Database{
		Client: redisClient,
	}
}
