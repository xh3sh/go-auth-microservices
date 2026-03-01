package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
)

// SetAPIKey СЃРѕС…СЂР°РЅСЏРµС‚ API РєР»СЋС‡ РІ Redis
func (r *redisRepository) SetAPIKey(ctx context.Context, apiKeyID string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixAPIKey, apiKeyID), jsonData, &ttl)
}

// GetAPIKey РёР·РІР»РµРєР°РµС‚ API РєР»СЋС‡ РёР· Redis
func (r *redisRepository) GetAPIKey(ctx context.Context, apiKeyID string, dest interface{}) error {
	val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixAPIKey, apiKeyID)).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// DeleteAPIKey СѓРґР°Р»СЏРµС‚ API РєР»СЋС‡ РёР· Redis
func (r *redisRepository) DeleteAPIKey(ctx context.Context, apiKeyID string) error {
	return r.Del(ctx, fmt.Sprintf("%s%s", constants.PrefixAPIKey, apiKeyID))
}
