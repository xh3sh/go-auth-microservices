package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/auth"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"github.com/xh3sh/go-auth-microservices/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware предоставляет методы для аутентификации запросов
type AuthMiddleware struct {
	jwtService     *auth.JWTService
	apiKeyService  *auth.APIKeyService
	sessionService *auth.SessionService
	repo           repository.Repository
	eventRepo      repository.EventRepository
}

func NewAuthMiddleware(jwt *auth.JWTService, api *auth.APIKeyService, session *auth.SessionService, repo repository.Repository, eventRepo repository.EventRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:     jwt,
		apiKeyService:  api,
		sessionService: session,
		repo:           repo,
		eventRepo:      eventRepo,
	}
}

// UniversalAuth проверяет все доступные способы авторизации
func (m *AuthMiddleware) UniversalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := ""
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			if t, err := c.Cookie("access_token"); err == nil {
				token = t
			}
		}

		if token != "" {
			claims, err := m.jwtService.ExtractClaims(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("auth_method", "jwt")
				c.Set("jwt_claims", claims)
				c.Request.Header.Set("X-User-ID", claims.UserID)
				c.Request.Header.Set("X-Auth-Method", "jwt")

				m.publishValidation(claims.UserID, "jwt", true, "", utils.GetTraceID(c.Request.Context()))

				c.Next()
				return
			}
		}

		oauthToken, err := c.Cookie("oauth_access_token")
		if err == nil && oauthToken != "" {
			claims, err := m.jwtService.ExtractClaims(oauthToken)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("auth_method", "oauth")
				c.Set("jwt_claims", claims)
				c.Request.Header.Set("X-User-ID", claims.UserID)
				c.Request.Header.Set("X-Auth-Method", "oauth")

				m.publishValidation(claims.UserID, "oauth", true, "", utils.GetTraceID(c.Request.Context()))

				c.Next()
				return
			}
		}

		keyVal := c.GetHeader("X-API-Key")
		if keyVal == "" {
			if kval, err := c.Cookie("api_key"); err == nil {
				keyVal = kval
			}
		}

		if keyVal != "" {
			apiKey, err := m.apiKeyService.ValidateAPIKey(keyVal)
			if err == nil {
				c.Set("user_id", apiKey.UserID)
				c.Set("auth_method", "api_key")
				c.Set("api_key_data", apiKey)
				c.Request.Header.Set("X-User-ID", apiKey.UserID)
				c.Request.Header.Set("X-Auth-Method", "api_key")

				m.publishValidation(apiKey.UserID, "api_key", true, "", utils.GetTraceID(c.Request.Context()))

				c.Next()
				return
			}
		}

		sessionID, err := c.Cookie("session_id")
		if err != nil {
			sessionID = c.GetHeader("X-Session-ID")
		}
		if sessionID != "" {
			session, err := m.sessionService.ValidateSession(c.Request.Context(), sessionID)
			if err == nil {
				c.Set("user_id", session.UserID)
				c.Set("auth_method", "session")
				c.Set("session_data", session)

				c.Request.Header.Set("X-User-ID", session.UserID)
				c.Request.Header.Set("X-Auth-Method", "session")

				m.publishValidation(session.UserID, "session", true, "", utils.GetTraceID(c.Request.Context()))

				c.Next()
				return
			}
		}

		username, password, ok := c.Request.BasicAuth()
		if ok && username != "" && password != "" {
			creds, err := m.repo.GetUserByUsername(c.Request.Context(), username)
			if err == nil && utils.CheckPasswordHash(password, creds.PasswordHash) {
				c.Set("user_id", creds.ID)
				c.Set("auth_method", "basic")

				c.Request.Header.Set("X-User-ID", creds.ID)
				c.Request.Header.Set("X-Auth-Method", "basic")

				m.publishValidation(creds.ID, "basic", true, "", utils.GetTraceID(c.Request.Context()))

				c.Next()
				return
			}
		}

		m.publishValidation("", "universal", false, "all auth methods failed", utils.GetTraceID(c.Request.Context()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		c.Abort()
	}
}

// JWTAuth проверяет JWT в заголовке Authorization или куках
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := ""
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			if t, err := c.Cookie("access_token"); err == nil {
				token = t
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token is required"})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ExtractClaims(token)
		if err != nil {
			m.publishValidation("", "jwt", false, "invalid or expired token", utils.GetTraceID(c.Request.Context()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("auth_method", "jwt")
		c.Set("jwt_claims", claims)
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Auth-Method", "jwt")
		m.publishValidation(claims.UserID, "jwt", true, "", utils.GetTraceID(c.Request.Context()))

		c.Next()
	}
}

// OAuthAuth проверяет OAuth JWT в куках
func (m *AuthMiddleware) OAuthAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		oauthToken, err := c.Cookie("oauth_access_token")
		if err != nil || oauthToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "OAuth token is required"})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ExtractClaims(oauthToken)
		if err != nil {
			m.publishValidation("", "oauth", false, "invalid or expired oauth token", utils.GetTraceID(c.Request.Context()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OAuth token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("auth_method", "oauth")
		c.Set("jwt_claims", claims)
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Auth-Method", "oauth")
		m.publishValidation(claims.UserID, "oauth", true, "", utils.GetTraceID(c.Request.Context()))

		c.Next()
	}
}

// BasicAuth проверяет заголовок Authorization: Basic
func (m *AuthMiddleware) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok || username == "" || password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Basic Auth required"})
			c.Abort()
			return
		}

		creds, err := m.repo.GetUserByUsername(c.Request.Context(), username)
		if err != nil || !utils.CheckPasswordHash(password, creds.PasswordHash) {
			m.publishValidation("", "basic", false, "invalid credentials", utils.GetTraceID(c.Request.Context()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		c.Set("user_id", creds.ID)
		c.Set("auth_method", "basic")
		c.Request.Header.Set("X-User-ID", creds.ID)
		c.Request.Header.Set("X-Auth-Method", "basic")
		m.publishValidation(creds.ID, "basic", true, "", utils.GetTraceID(c.Request.Context()))

		c.Next()
	}
}

func (m *AuthMiddleware) publishValidation(userID, method string, isValid bool, reason string, traceID string) {
	if m.eventRepo == nil {
		return
	}

	m.eventRepo.PublishTokenValidationEvent(models.TokenValidationEvent{
		UserID:       userID,
		AuthMethod:   method,
		IsValid:      isValid,
		TraceID:      traceID,
		ErrorMessage: reason,
		Timestamp:    time.Now(),
	})
}

// APIKeyAuth проверяет API-ключ в заголовке X-API-Key
func (m *AuthMiddleware) APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		keyVal := c.GetHeader("X-API-Key")

		if keyVal == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key is required"})
			c.Abort()
			return
		}

		apiKey, err := m.apiKeyService.ValidateAPIKey(keyVal)
		if err != nil {
			m.publishValidation("", "api_key", false, "invalid api key", utils.GetTraceID(c.Request.Context()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Set("user_id", apiKey.UserID)
		c.Set("auth_method", "api_key")
		c.Set("api_key_data", apiKey)
		c.Request.Header.Set("X-User-ID", apiKey.UserID)
		c.Request.Header.Set("X-Auth-Method", "api_key")
		m.publishValidation(apiKey.UserID, "api_key", true, "", utils.GetTraceID(c.Request.Context()))

		c.Next()
	}
}

// SessionAuth проверяет ID сессии в куках или заголовке
func (m *AuthMiddleware) SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			sessionID = c.GetHeader("X-Session-ID")
		}

		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is required"})
			c.Abort()
			return
		}

		session, err := m.sessionService.ValidateSession(c.Request.Context(), sessionID)
		if err != nil {
			m.publishValidation("", "session", false, "invalid or expired session", utils.GetTraceID(c.Request.Context()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Set("auth_method", "session")
		c.Set("session_data", session)
		c.Request.Header.Set("X-User-ID", session.UserID)
		c.Request.Header.Set("X-Auth-Method", "session")
		m.publishValidation(session.UserID, "session", true, "", utils.GetTraceID(c.Request.Context()))

		c.Next()
	}
}
