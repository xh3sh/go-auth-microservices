package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLogs РІРѕР·РІСЂР°С‰Р°РµС‚ СЃРїРёСЃРѕРє РІСЃРµС… Р»РѕРіРѕРІ СЃРёСЃС‚РµРјС‹
func (h *Handler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.repo.GetLogs(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// FilterLogs С„РёР»СЊС‚СЂСѓРµС‚ Р»РѕРіРё РїРѕ Р·Р°РґР°РЅРЅС‹Рј РїР°СЂР°РјРµС‚СЂР°Рј
func (h *Handler) FilterLogs(c *gin.Context) {
	userID := c.Query("user_id")
	service := c.Query("service")
	logType := c.Query("type")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.repo.FilterLogs(c.Request.Context(), userID, service, logType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
