package handler

import (
	"context"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleJWTLogin РѕР±СЂР°Р±Р°С‚С‹РІР°РµС‚ РІС…РѕРґ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ Рё РІС‹РґР°РµС‚ РїР°СЂСѓ JWT С‚РѕРєРµРЅРѕРІ
func (h *Handler) HandleJWTLogin(c *gin.Context) {
	var req models.LoginRequest
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

	tokenPair, err := h.jwtService.GenerateTokenPair(creds.ID, creds.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	authEvent := models.AuthEvent{
		UserID:    creds.ID,
		EventType: "jwt_login",
		Timestamp: time.Now(),
	}
	if h.eventRepo != nil {
		go h.eventRepo.PublishAuthEvent(authEvent)
	}

	c.Set("user_id", creds.ID)
	c.Set("auth_method", "jwt")

	c.SetCookie("access_token", tokenPair.AccessToken, int(h.config.JWTExpiration), "/", "", false, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken, int(h.config.JWTRefreshExpiration), "/", "", false, true)
	c.SetCookie("jwt_logged_in", "true", int(h.config.JWTRefreshExpiration), "/", "", false, false)

	c.JSON(http.StatusOK, tokenPair)
}

// HandleJWTRefresh РѕР±РЅРѕРІР»СЏРµС‚ access С‚РѕРєРµРЅ, РёСЃРїРѕР»СЊР·СѓСЏ РІР°Р»РёРґРЅС‹Р№ refresh С‚РѕРєРµРЅ
func (h *Handler) HandleJWTRefresh(c *gin.Context) {
	var req models.RefreshTokenRequest
	refreshToken := ""

	if err := c.ShouldBindJSON(&req); err == nil {
		refreshToken = req.RefreshToken
	}

	if refreshToken == "" {
		if cookie, err := c.Cookie("refresh_token"); err == nil {
			refreshToken = cookie
		}
	}

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	newAccessToken, err := h.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.SetCookie("access_token", newAccessToken, int(h.config.JWTExpiration), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
	})
}

// HandleJWTValidate РІРѕР·РІСЂР°С‰Р°РµС‚ РґР°РЅРЅС‹Рµ РёР· С‚РѕРєРµРЅР°, РµСЃР»Рё РѕРЅ РІР°Р»РёРґРµРЅ
func (h *Handler) HandleJWTValidate(c *gin.Context) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, claims)
}

// HandleJWTLogout Р°РЅРЅСѓР»РёСЂСѓРµС‚ С‚РµРєСѓС‰РёР№ JWT С‚РѕРєРµРЅ
func (h *Handler) HandleJWTLogout(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		token = c.GetHeader("Authorization")
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
		}
	}

	if token != "" {
		ctx := context.Background()
		h.repo.SetBlacklist(ctx, token, time.Duration(h.config.JWTExpiration)*time.Second)

		if claims, err := h.jwtService.ExtractClaims(token); err == nil {
			c.Set("user_id", claims.UserID)
		}
	}

	c.Set("auth_method", "jwt")

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.SetCookie("jwt_logged_in", "", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
