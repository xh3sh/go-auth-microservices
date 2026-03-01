package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
)

// SetSession СЃРѕС…СЂР°РЅСЏРµС‚ РґР°РЅРЅС‹Рµ СЃРµСЃСЃРёРё РІ Redis
func (r *redisRepository) SetSession(ctx context.Context, sessionID string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixSession, sessionID), jsonData, &ttl)
}

// GetSession РёР·РІР»РµРєР°РµС‚ РґР°РЅРЅС‹Рµ СЃРµСЃСЃРёРё РёР· Redis
func (r *redisRepository) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixSession, sessionID)).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// DeleteSession СѓРґР°Р»СЏРµС‚ СЃРµСЃСЃРёСЋ РёР· Redis
func (r *redisRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.Del(ctx, fmt.Sprintf("%s%s", constants.PrefixSession, sessionID))
}
