package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
)

// APIKeyService РїСЂРµРґРѕСЃС‚Р°РІР»СЏРµС‚ С„СѓРЅРєС†РёРѕРЅР°Р» РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ API РєР»СЋС‡Р°РјРё
type APIKeyService struct {
	repo repository.APIKeyRepository
}

func NewAPIKeyService(repo repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		repo: repo,
	}
}

// GenerateAPIKey СЃРѕР·РґР°РµС‚ РЅРѕРІС‹Р№ API РєР»СЋС‡ РґР»СЏ СѓРєР°Р·Р°РЅРЅРѕРіРѕ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
func (a *APIKeyService) GenerateAPIKey(userID string, name string) (*models.APIKeyResponse, error) {
	keyID, err := utils.GenerateRandomString(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key ID: %w", err)
	}

	keyValue, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	now := time.Now()
	expiresAt := now.Add(constants.APIKeyTTL)

	apiKey := &models.APIKeyResponse{
		ID:        keyID,
		Key:       fmt.Sprintf("%s.%s", keyID, keyValue),
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	storageData := *apiKey
	storageData.Key = keyValue

	ctx := context.Background()
	if err := a.repo.SetAPIKey(ctx, keyID, &storageData, constants.APIKeyTTL); err != nil {
		return nil, fmt.Errorf("failed to save API key to repository: %w", err)
	}

	return apiKey, nil
}

// ValidateAPIKey РїСЂРѕРІРµСЂСЏРµС‚ API РєР»СЋС‡ (С„РѕСЂРјР°С‚ ID.Secret)
func (a *APIKeyService) ValidateAPIKey(fullKey string) (*models.APIKeyResponse, error) {
	if len(fullKey) == 0 {
		return nil, fmt.Errorf("missing api key")
	}

	parts := strings.Split(fullKey, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid api key format")
	}

	keyID := parts[0]
	keyValue := parts[1]

	var apiKey models.APIKeyResponse
	ctx := context.Background()
	err := a.repo.GetAPIKey(ctx, keyID, &apiKey)
	if err != nil {
		return nil, fmt.Errorf("key not found")
	}

	if apiKey.Key != keyValue {
		return nil, fmt.Errorf("invalid key value")
	}

	if time.Now().After(apiKey.ExpiresAt) {
		return nil, fmt.Errorf("key expired")
	}

	return &apiKey, nil
}

// RevokeAPIKey РѕС‚Р·С‹РІР°РµС‚ API РєР»СЋС‡ РїРѕ РµРіРѕ ID
func (a *APIKeyService) RevokeAPIKey(keyID string) error {
	if len(keyID) == 0 {
		return fmt.Errorf("invalid key ID")
	}

	ctx := context.Background()
	return a.repo.DeleteAPIKey(ctx, keyID)
}
