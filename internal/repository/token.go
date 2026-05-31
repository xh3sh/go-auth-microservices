package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/redis/go-redis/v9"
)

// SetBlacklist добавляет токен в черный список в Redis
func (r *redisRepository) SetBlacklist(ctx context.Context, tokenID string, ttl time.Duration) error {
	return r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixBlacklist, tokenID), constants.ValueRevoked, &ttl)
}

// IsBlacklisted проверяет, находится ли токен в черном списке
func (r *redisRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixBlacklist, tokenID)).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == constants.ValueRevoked, err
}

// SetRefreshToken сохраняет refresh токен в Redis
func (r *redisRepository) SetRefreshToken(ctx context.Context, tokenID string, userID string, ttl time.Duration) error {
	return r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixRefreshToken, tokenID), userID, &ttl)
}

// GetRefreshToken извлекает userID, связанный с refresh токеном
func (r *redisRepository) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixRefreshToken, tokenID)).Result()
	return val, err
}
