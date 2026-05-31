package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService предоставляет методы для работы с JWT токенами
type JWTService struct {
	secret            string
	expiration        int64
	refreshExpiration int64
	repo              repository.TokenRepository
}

func NewJWTService(secret string, expiration int64, refreshExpiration int64, repo repository.TokenRepository) *JWTService {
	return &JWTService{
		secret:            secret,
		expiration:        expiration,
		refreshExpiration: refreshExpiration,
		repo:              repo,
	}
}

// GenerateTokenPair генерирует пару токенов: access и refresh
func (j *JWTService) GenerateTokenPair(userID string, username string) (*models.JWTTokenPair, error) {
	accessToken, err := j.generateAccessToken(userID, username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateRefreshToken(userID, username)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ttl := time.Duration(j.refreshExpiration) * time.Second
	err = j.repo.SetRefreshToken(ctx, refreshToken, userID, ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &models.JWTTokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    j.expiration,
		TokenType:    "Bearer",
	}, nil
}

func (j *JWTService) generateAccessToken(userID string, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.expiration) * time.Second)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"type":     "access",
		"exp":      expiresAt.Unix(),
		"iat":      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTService) generateRefreshToken(userID string, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.refreshExpiration) * time.Second)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"type":     "refresh",
		"exp":      expiresAt.Unix(),
		"iat":      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ValidateToken проверяет валидность токена и наличие его в черном списке
func (j *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	ctx := context.Background()
	blacklisted, err := j.repo.IsBlacklisted(ctx, tokenString)
	if err == nil && blacklisted {
		return nil, errors.New("token is revoked")
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// ExtractClaims извлекает данные из JWT токена
func (j *JWTService) ExtractClaims(tokenString string) (*models.JWTClaims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userID, _ := claims["user_id"].(string)
	username, _ := claims["username"].(string)
	exp, _ := claims["exp"].(float64)
	iat, _ := claims["iat"].(float64)

	return &models.JWTClaims{
		UserID:   userID,
		Username: username,
		Exp:      int64(exp),
		Iat:      int64(iat),
	}, nil
}

// RefreshAccessToken обновляет access токен на основе валидного refresh токена
func (j *JWTService) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if t, ok := claims["type"].(string); !ok || t != "refresh" {
		return "", errors.New("token is not a refresh token")
	}

	userID, _ := claims["user_id"].(string)
	username, _ := claims["username"].(string)

	ctx := context.Background()
	storedUserID, err := j.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil || storedUserID != userID {
		return "", errors.New("refresh token not found or expired")
	}

	return j.generateAccessToken(userID, username)
}

// RevokeToken добавляет токен в черный список
func (j *JWTService) RevokeToken(tokenString string, ttl time.Duration) error {
	ctx := context.Background()
	return j.repo.SetBlacklist(ctx, tokenString, ttl)
}
