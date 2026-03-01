package repository

import (
	"context"
	"time"
)

type TokenRepository interface {
	SetBlacklist(ctx context.Context, tokenID string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
	SetRefreshToken(ctx context.Context, tokenID string, userID string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, tokenID string) (string, error)
}
