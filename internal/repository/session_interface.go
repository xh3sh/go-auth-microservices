package repository

import (
	"context"
	"time"
)

type SessionRepository interface {
	SetSession(ctx context.Context, sessionID string, data interface{}, ttl time.Duration) error
	GetSession(ctx context.Context, sessionID string, dest interface{}) error
	DeleteSession(ctx context.Context, sessionID string) error
}
