package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/utils"

	"github.com/gin-gonic/gin"
)

// HandleCreateSession СЃРѕР·РґР°РµС‚ РЅРѕРІСѓСЋ СЃРµСЃСЃРёСЋ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
func (h *Handler) HandleCreateSession(c *gin.Context) {
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

	ctx := context.Background()
	session, err := h.sessionService.CreateSession(ctx, creds.ID, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	authEvent := models.AuthEvent{
		UserID:    session.UserID,
		EventType: "session_created",
		Timestamp: time.Now(),
	}
	if h.eventRepo != nil {
		go h.eventRepo.PublishAuthEvent(authEvent)
	}

	c.Set("user_id", session.UserID)
	c.Set("auth_method", "session")

	c.SetCookie("session_id", session.SessionID, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, session)
}

// HandleValidateSession РїСЂРѕРІРµСЂСЏРµС‚ РІР°Р»РёРґРЅРѕСЃС‚СЊ С‚РµРєСѓС‰РµР№ СЃРµСЃСЃРёРё
func (h *Handler) HandleValidateSession(c *gin.Context) {
	sessionData, exists := c.Get("session_data")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	session := sessionData.(*models.Session)
	c.JSON(http.StatusOK, gin.H{"valid": true, "user_id": session.UserID})
}

// HandleRevokeSession РѕС‚Р·С‹РІР°РµС‚ С‚РµРєСѓС‰СѓСЋ СЃРµСЃСЃРёСЋ
func (h *Handler) HandleRevokeSession(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if len(sessionID) == 0 {
		if sid, err := c.Cookie("session_id"); err == nil {
			sessionID = sid
		}
	}

	if len(sessionID) != 0 {
		ctx := context.Background()
		if session, err := h.sessionService.ValidateSession(ctx, sessionID); err == nil {
			c.Set("user_id", session.UserID)
		}
		h.sessionService.RevokeSession(ctx, sessionID)
	}

	c.Set("auth_method", "session")
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}
