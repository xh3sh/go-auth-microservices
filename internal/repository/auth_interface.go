package repository

import (
	"context"
	"github.com/xh3sh/go-auth-microservices/internal/models"
)

type AuthRepository interface {
	CreateUserCredentials(ctx context.Context, creds models.UserCredentials) error
	GetUserByUsername(ctx context.Context, username string) (models.UserCredentials, error)
}
