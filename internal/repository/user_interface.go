package repository

import (
	"context"
	"github.com/xh3sh/go-auth-microservices/internal/models"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id string) error
}
