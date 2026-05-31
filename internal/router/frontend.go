package router

import (
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewFrontendRouter создает роутер для фронтенд-сервиса
func NewFrontendRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.SetHTMLTemplate(handler.NewTemplates())
	r.Use(middleware.Logging())
	r.Static("/static", "./static")

	r.GET("/", h.HandleHome)

	return r
}
