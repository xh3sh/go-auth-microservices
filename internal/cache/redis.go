package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient РїСЂРµРґРѕСЃС‚Р°РІР»СЏРµС‚ РѕР±РµСЂС‚РєСѓ РЅР°Рґ РєР»РёРµРЅС‚РѕРј Redis
type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(host, port, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
