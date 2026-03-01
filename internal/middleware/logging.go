package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"

	"github.com/gin-gonic/gin"
)

// Logging РІРѕР·РІСЂР°С‰Р°РµС‚ СЃС‚Р°РЅРґР°СЂС‚РЅС‹Р№ Р»РѕРіРіРµСЂ Gin
func Logging() gin.HandlerFunc {
	return gin.Logger()
}

// APILogger Р»РѕРіРёСЂСѓРµС‚ API Р·Р°РїСЂРѕСЃС‹ Рё РѕС‚РїСЂР°РІР»СЏРµС‚ СЃРѕР±С‹С‚РёСЏ РІ СЂРµРїРѕР·РёС‚РѕСЂРёР№ СЃРѕР±С‹С‚РёР№
func APILogger(eventRepo repository.EventRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if eventRepo == nil || c.Request.Method == "OPTIONS" {
			return
		}

		status := c.Writer.Status()
		if status == http.StatusMovedPermanently || status == http.StatusPermanentRedirect || status == http.StatusTemporaryRedirect {
			return
		}

		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/auth") {
			return
		}

		userIDRaw, _ := c.Get("user_id")
		userID, _ := userIDRaw.(string)
		if userID == "" {
			userID = c.GetHeader("X-User-ID")
		}

		authMethodRaw, _ := c.Get("auth_method")
		authMethod, _ := authMethodRaw.(string)
		if authMethod == "" {
			authMethod = c.GetHeader("X-Auth-Method")
			if authMethod == "" {
				authMethod = "none"
			}
		}

		event := models.APIGatewayEvent{
			UserID:     userID,
			AuthMethod: authMethod,
			Action:     c.Request.Method,
			Resource:   c.Request.URL.Path,
			Timestamp:  time.Now(),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		go eventRepo.PublishAPIGatewayEvent(event)
	}
}
