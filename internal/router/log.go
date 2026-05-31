package router

import (
	"github.com/xh3sh/go-auth-microservices/internal/handler"

	"github.com/gin-gonic/gin"
)

// NewLogRouter создает роутер для сервиса логов
func NewLogRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.SetHTMLTemplate(handler.NewTemplates())

	auditGroup := r.Group("/audit/entries")
	{
		auditGroup.GET("", h.GetLogs)
		auditGroup.GET("/filter", h.FilterLogs)
	}

	return r
}
