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

// AuthMiddleware РїСЂРµРґРѕСЃС‚Р°РІР»СЏРµС‚ РјРµС‚РѕРґС‹ РґР»СЏ Р°СѓС‚РµРЅС‚РёС„РёРєР°С†РёРё Р·Р°РїСЂРѕСЃРѕРІ
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

// UniversalAuth РїСЂРѕРІРµСЂСЏРµС‚ РІСЃРµ РґРѕСЃС‚СѓРїРЅС‹Рµ СЃРїРѕСЃРѕР±С‹ Р°РІС‚РѕСЂРёР·Р°С†РёРё
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

				m.publishValidation(claims.UserID, "jwt", true, "")

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

				m.publishValidation(claims.UserID, "oauth", true, "")

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

				m.publishValidation(apiKey.UserID, "api_key", true, "")

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

				m.publishValidation(session.UserID, "session", true, "")

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

				m.publishValidation(creds.ID, "basic", true, "")

				c.Next()
				return
			}
		}

		m.publishValidation("", "universal", false, "all auth methods failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		c.Abort()
	}
}

// JWTAuth РїСЂРѕРІРµСЂСЏРµС‚ JWT РІ Р·Р°РіРѕР»РѕРІРєРµ Authorization РёР»Рё РєСѓРєР°С…
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
			m.publishValidation("", "jwt", false, "invalid or expired token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("auth_method", "jwt")
		c.Set("jwt_claims", claims)
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Auth-Method", "jwt")
		m.publishValidation(claims.UserID, "jwt", true, "")

		c.Next()
	}
}

// OAuthAuth РїСЂРѕРІРµСЂСЏРµС‚ OAuth JWT РІ РєСѓРєР°С…
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
			m.publishValidation("", "oauth", false, "invalid or expired oauth token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OAuth token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("auth_method", "oauth")
		c.Set("jwt_claims", claims)
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Auth-Method", "oauth")
		m.publishValidation(claims.UserID, "oauth", true, "")

		c.Next()
	}
}

// BasicAuth РїСЂРѕРІРµСЂСЏРµС‚ Р·Р°РіРѕР»РѕРІРѕРє Authorization: Basic
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
			m.publishValidation("", "basic", false, "invalid credentials")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		c.Set("user_id", creds.ID)
		c.Set("auth_method", "basic")
		c.Request.Header.Set("X-User-ID", creds.ID)
		c.Request.Header.Set("X-Auth-Method", "basic")
		m.publishValidation(creds.ID, "basic", true, "")

		c.Next()
	}
}

func (m *AuthMiddleware) publishValidation(userID, method string, isValid bool, reason string) {
	if m.eventRepo == nil {
		return
	}

	m.eventRepo.PublishTokenValidationEvent(models.TokenValidationEvent{
		UserID:       userID,
		AuthMethod:   method,
		IsValid:      isValid,
		ErrorMessage: reason,
		Timestamp:    time.Now(),
	})
}

// APIKeyAuth РїСЂРѕРІРµСЂСЏРµС‚ API-РєР»СЋС‡ РІ Р·Р°РіРѕР»РѕРІРєРµ X-API-Key
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
			m.publishValidation("", "api_key", false, "invalid api key")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Set("user_id", apiKey.UserID)
		c.Set("auth_method", "api_key")
		c.Set("api_key_data", apiKey)
		c.Request.Header.Set("X-User-ID", apiKey.UserID)
		c.Request.Header.Set("X-Auth-Method", "api_key")
		m.publishValidation(apiKey.UserID, "api_key", true, "")

		c.Next()
	}
}

// SessionAuth РїСЂРѕРІРµСЂСЏРµС‚ ID СЃРµСЃСЃРёРё РІ РєСѓРєР°С… РёР»Рё Р·Р°РіРѕР»РѕРІРєРµ
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
			m.publishValidation("", "session", false, "invalid or expired session")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Set("auth_method", "session")
		c.Set("session_data", session)
		c.Request.Header.Set("X-User-ID", session.UserID)
		c.Request.Header.Set("X-Auth-Method", "session")
		m.publishValidation(session.UserID, "session", true, "")

		c.Next()
	}
}
