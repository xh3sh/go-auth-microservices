package repository

type Repository interface {
	AuthRepository
	SessionRepository
	APIKeyRepository
	TokenRepository
	UserRepository
	LogRepository
}
