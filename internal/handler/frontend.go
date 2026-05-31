package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleHome рендерит главную страницу приложения
func (h *Handler) HandleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}
