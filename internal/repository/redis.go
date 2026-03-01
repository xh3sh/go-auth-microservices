package repository

import (
	"context"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
	prefix string
}

func NewRedisRepository(client *redis.Client) Repository {
	return &redisRepository{
		client: client,
		prefix: constants.RedisPrefix,
	}
}

func (r *redisRepository) buildKey(key string) string {
	return r.prefix + key
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration *time.Duration) error {
	ttl := constants.DefaultTTL
	if expiration != nil {
		ttl = *expiration
	}
	return r.client.Set(ctx, r.buildKey(key), value, ttl).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, r.buildKey(key))
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.buildKey(key)).Err()
}

func (r *redisRepository) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanIterator {
	return r.client.Scan(ctx, cursor, r.buildKey(match), count).Iterator()
}
