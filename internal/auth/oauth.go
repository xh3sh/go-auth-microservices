package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
)

// OAuthService СЂРµР°Р»РёР·СѓРµС‚ Р»РѕРіРёРєСѓ OAuth 2.0 СЃРµСЂРІРµСЂР°
type OAuthService struct {
	repo repository.TokenRepository
}

func NewOAuthService(repo repository.TokenRepository) *OAuthService {
	return &OAuthService{
		repo: repo,
	}
}

// ValidateClient РїСЂРѕРІРµСЂСЏРµС‚ СѓС‡РµС‚РЅС‹Рµ РґР°РЅРЅС‹Рµ OAuth РєР»РёРµРЅС‚Р°
func (o *OAuthService) ValidateClient(clientID, clientSecret string) (bool, error) {
	if clientID == "test_client" && clientSecret == "test_secret" {
		return true, nil
	}
	return len(clientID) > 0 && len(clientSecret) > 0, nil
}

// GenerateOAuthToken СЃРѕР·РґР°РµС‚ РЅРѕРІС‹Р№ OAuth С‚РѕРєРµРЅ
func (o *OAuthService) GenerateOAuthToken(clientID string, scope string) (*models.OAuthToken, error) {
	if len(clientID) == 0 {
		return nil, errors.New("invalid client ID")
	}

	now := time.Now()
	expiresIn := int64(3600)

	accessToken, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	token := &models.OAuthToken{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Scope:       scope,
		IssuedAt:    now.Unix(),
	}

	return token, nil
}

// ValidateOAuthToken РїСЂРѕРІРµСЂСЏРµС‚ РІР°Р»РёРґРЅРѕСЃС‚СЊ OAuth С‚РѕРєРµРЅР°
func (o *OAuthService) ValidateOAuthToken(token string) (bool, error) {
	if len(token) == 0 {
		return false, nil
	}
	
	ctx := context.Background()
	revoked, err := o.repo.IsBlacklisted(ctx, token)
	if err == nil && revoked {
		return false, nil
	}

	return true, nil
}

// RevokeOAuthToken РѕС‚Р·С‹РІР°РµС‚ OAuth С‚РѕРєРµРЅ
func (o *OAuthService) RevokeOAuthToken(token string) error {
	if len(token) == 0 {
		return errors.New("invalid token")
	}
	
	ctx := context.Background()
	return o.repo.SetBlacklist(ctx, token, time.Hour)
}

// ExchangeAuthorizationCode РѕР±РјРµРЅРёРІР°РµС‚ РєРѕРґ Р°РІС‚РѕСЂРёР·Р°С†РёРё РЅР° С‚РѕРєРµРЅ
func (o *OAuthService) ExchangeAuthorizationCode(clientID, clientSecret, code, redirectURI string) (*models.OAuthToken, error) {
	if len(clientID) == 0 || len(code) == 0 || len(redirectURI) == 0 {
		return nil, errors.New("invalid parameters")
	}

	valid, err := o.ValidateClient(clientID, clientSecret)
	if err != nil || !valid {
		return nil, errors.New("invalid client credentials")
	}

	return o.GenerateOAuthToken(clientID, "")
}

// ValidateScope РїСЂРѕРІРµСЂСЏРµС‚ РїСЂР°РІР° РґРѕСЃС‚СѓРїР° (scope) РєР»РёРµРЅС‚Р°
func (o *OAuthService) ValidateScope(clientID, requestedScope string) (bool, error) {
	if len(clientID) == 0 || len(requestedScope) == 0 {
		return false, nil
	}

	allowedScopes := []string{"read", "write", "admin"}
	requestedScopes := strings.Fields(requestedScope)

	for _, rs := range requestedScopes {
		found := false
		for _, as := range allowedScopes {
			if rs == as {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	return true, nil
}
