package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleBasicAuthValidate РїСЂРѕРІРµСЂСЏРµС‚ Р°РІС‚РѕСЂРёР·Р°С†РёСЋ С‡РµСЂРµР· Basic Auth РёР»Рё Bearer
func (h *Handler) HandleBasicAuthValidate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authMethod, _ := c.Get("auth_method")

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user_id":       userID,
		"method":        authMethod,
	})
}
