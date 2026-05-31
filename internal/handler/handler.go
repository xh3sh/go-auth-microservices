package handler

import (
	"github.com/xh3sh/go-auth-microservices/internal/auth"
	"github.com/xh3sh/go-auth-microservices/internal/config"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
)

// Handler является центральным узлом для всех HTTP хендлеров приложения
type Handler struct {
	jwtService     *auth.JWTService
	apiKeyService  *auth.APIKeyService
	oauthService   *auth.OAuthService
	sessionService *auth.SessionService
	repo           repository.Repository
	eventRepo      repository.EventRepository
	config         *config.Config
}

func NewHandler(
	jwtService *auth.JWTService,
	apiKeyService *auth.APIKeyService,
	oauthService *auth.OAuthService,
	sessionService *auth.SessionService,
	repo repository.Repository,
	eventRepo repository.EventRepository,
	cfg *config.Config,
) *Handler {
	return &Handler{
		jwtService:     jwtService,
		apiKeyService:  apiKeyService,
		oauthService:   oauthService,
		sessionService: sessionService,
		repo:           repo,
		eventRepo:      eventRepo,
		config:         cfg,
	}
}

func (h *Handler) GetEventRepo() repository.EventRepository {
	return h.eventRepo
}
