package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
)

// CreateUserCredentials СЃРѕС…СЂР°РЅСЏРµС‚ СѓС‡РµС‚РЅС‹Рµ РґР°РЅРЅС‹Рµ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ Рё СЃРѕР·РґР°РµС‚ РїРѕРёСЃРєРѕРІС‹Р№ РёРЅРґРµРєСЃ
func (r *redisRepository) CreateUserCredentials(ctx context.Context, creds models.UserCredentials) error {
	jsonData, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	if err := r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixUser, creds.ID), jsonData, nil); err != nil {
		return err
	}

	return r.Set(ctx, fmt.Sprintf("%s%s", constants.PrefixUsername, creds.Username), creds.ID, nil)
}

// GetUserByUsername РЅР°С…РѕРґРёС‚ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ РїРѕ РµРіРѕ РёРјРµРЅРё
func (r *redisRepository) GetUserByUsername(ctx context.Context, username string) (models.UserCredentials, error) {
	userID, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixUsername, username)).Result()
	if err == nil && userID != "" {
		val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixUser, userID)).Result()
		if err == nil {
			var creds models.UserCredentials
			if err := json.Unmarshal([]byte(val), &creds); err == nil {
				return creds, nil
			}
		}
	}

	return models.UserCredentials{}, fmt.Errorf("user not found")
}
