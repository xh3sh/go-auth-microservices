package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
)

// GetUser возвращает профиль пользователя по его ID
func (r *redisRepository) GetUser(ctx context.Context, id string) (models.User, error) {
	var user models.User
	val, err := r.Get(ctx, fmt.Sprintf("%s%s", constants.PrefixUser, id)).Result()
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

// GetAllUsers возвращает список всех пользователей системы
func (r *redisRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users := []models.User{}
	iter := r.Scan(ctx, 0, fmt.Sprintf("%s*", constants.PrefixUser), 0)
	for iter.Next(ctx) {
		val, err := r.client.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}
		var user models.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			users = append(users, user)
		}
	}
	return users, iter.Err()
}

// DeleteUser удаляет профиль пользователя и индекс
func (r *redisRepository) DeleteUser(ctx context.Context, id string) error {
	user, err := r.GetUser(ctx, id)
	if err == nil {
		_ = r.Del(ctx, fmt.Sprintf("%s%s", constants.PrefixUsername, user.Username))
	}
	return r.Del(ctx, fmt.Sprintf("%s%s", constants.PrefixUser, id))
}
