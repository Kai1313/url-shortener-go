package database

import (
    "context"
    "fmt"
    "log"

    "github.com/Kai1313/url-shortener-fullstack/backend/internal/config"
    "github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
        Password: cfg.RedisPassword,
        DB:       0,
    })

    // Test connection
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        log.Printf("Warning: Redis connection failed: %v", err)
        return client // Return anyway, we'll handle errors in handlers
    }

    log.Println("Redis connection established")
    return client
}