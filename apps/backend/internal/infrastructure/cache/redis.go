// apps/backend/internal/infrastructure/cache/redis.go

package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig شامل تنظیمات اتصال به Redis است.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// ConnectRedis یک کلاینت Redis ایجاد و برمی‌گرداند.
func ConnectRedis(cfg RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// بررسی صحت اتصال
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("FATAL: Failed to connect to Redis: %v", err)
	}

	log.Println("INFO: Redis connection established successfully.")
	return client
}