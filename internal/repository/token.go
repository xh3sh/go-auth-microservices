package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/redis/go-redis/v9"
)

// SetBlacklist РґРѕР±Р°РІР»СЏРµС‚ С‚РѕРєРµРЅ РІ С‡РµСЂРЅС‹Р№ СЃРїРёСЃРѕРє РІ Redis
func (r *redisRepository) SetBlacklist(ctx context.Context, tokenID string, ttl time.Duration) error {
	return r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixBlacklist, tokenID), constants.ValueRevoked, &ttl)
}

// IsBlacklisted РїСЂРѕРІРµСЂСЏРµС‚, РЅР°С…РѕРґРёС‚СЃСЏ Р»Рё С‚РѕРєРµРЅ РІ С‡РµСЂРЅРѕРј СЃРїРёСЃРєРµ
func (r *redisRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixBlacklist, tokenID)).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == constants.ValueRevoked, err
}

// SetRefreshToken СЃРѕС…СЂР°РЅСЏРµС‚ refresh С‚РѕРєРµРЅ РІ Redis
func (r *redisRepository) SetRefreshToken(ctx context.Context, tokenID string, userID string, ttl time.Duration) error {
	return r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixRefreshToken, tokenID), userID, &ttl)
}

// GetRefreshToken РёР·РІР»РµРєР°РµС‚ userID, СЃРІСЏР·Р°РЅРЅС‹Р№ СЃ refresh С‚РѕРєРµРЅРѕРј
func (r *redisRepository) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixRefreshToken, tokenID)).Result()
	return val, err
}
