package repository

import (
	"context"
	"time"
)

type APIKeyRepository interface {
	SetAPIKey(ctx context.Context, apiKeyID string, data interface{}, ttl time.Duration) error
	GetAPIKey(ctx context.Context, apiKeyID string, dest interface{}) error
	DeleteAPIKey(ctx context.Context, apiKeyID string) error
}
