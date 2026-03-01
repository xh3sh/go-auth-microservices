package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/utils"

	"github.com/gin-gonic/gin"
)

// HandleGenerateAPIKey РіРµРЅРµСЂРёСЂСѓРµС‚ РЅРѕРІС‹Р№ API РєР»СЋС‡ РґР»СЏ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
func (h *Handler) HandleGenerateAPIKey(c *gin.Context) {
	var req struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
		Name     string `json:"name" form:"name"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	creds, err := h.repo.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, creds.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	apiKey, err := h.apiKeyService.GenerateAPIKey(creds.ID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	c.Set("user_id", creds.ID)
	c.Set("auth_method", "api_key")

	cookieMaxAge := int(constants.APIKeyTTL.Seconds())
	c.SetCookie("api_key", apiKey.Key, cookieMaxAge, "/", "", false, true)

	c.JSON(http.StatusOK, apiKey)
}

// HandleValidateAPIKey РїСЂРѕРІРµСЂСЏРµС‚ РІР°Р»РёРґРЅРѕСЃС‚СЊ API РєР»СЋС‡Р°
func (h *Handler) HandleValidateAPIKey(c *gin.Context) {
	apiKey, exists := c.Get("api_key_data")
	if !exists {
		var req struct {
			Key string `json:"key" form:"key"`
		}
		if err := c.ShouldBind(&req); err == nil && req.Key != "" {
			validKey, err := h.apiKeyService.ValidateAPIKey(req.Key)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"valid": true, "user_id": validKey.UserID})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "user_id": apiKey.(*models.APIKeyResponse).UserID})
}

// HandleRevokeAPIKey РѕС‚Р·С‹РІР°РµС‚ (Р°РЅРЅСѓР»РёСЂСѓРµС‚) API РєР»СЋС‡
func (h *Handler) HandleRevokeAPIKey(c *gin.Context) {
	var req struct {
		Key string `json:"key" form:"key"`
	}
	if err := c.ShouldBind(&req); err != nil {
		if k, err := c.Cookie("api_key"); err == nil {
			req.Key = k
		}
	}

	if req.Key != "" {
		if apiKey, err := h.apiKeyService.ValidateAPIKey(req.Key); err == nil {
			c.Set("user_id", apiKey.UserID)
		}

		parts := strings.Split(req.Key, ".")
		if len(parts) == 2 {
			keyID := parts[0]
			err := h.apiKeyService.RevokeAPIKey(keyID)
			if err == nil {
				ctx := context.Background()
				h.repo.DeleteAPIKey(ctx, keyID)
			}
		}
	}

	c.Set("auth_method", "api_key")
	c.SetCookie("api_key", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked and cookies cleared"})
}
