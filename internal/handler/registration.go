package handler

import (
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleRegister СЂРµРіРёСЃС‚СЂРёСЂСѓРµС‚ РЅРѕРІРѕРіРѕ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ РІ СЃРёСЃС‚РµРјРµ
func (h *Handler) HandleRegister(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	existingUser, err := h.repo.GetUserByUsername(c.Request.Context(), req.Username)
	if err == nil && existingUser.Username != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	id, err := utils.GenerateRandomString(8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate user ID"})
		return
	}

	creds := models.UserCredentials{
		ID:           id,
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}

	if err := h.repo.CreateUserCredentials(c.Request.Context(), creds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, models.User{
		ID:       creds.ID,
		Username: creds.Username,
	})
}
