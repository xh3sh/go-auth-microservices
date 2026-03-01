package router

import (
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewFrontendRouter СЃРѕР·РґР°РµС‚ СЂРѕСѓС‚РµСЂ РґР»СЏ С„СЂРѕРЅС‚РµРЅРґ-СЃРµСЂРІРёСЃР°
func NewFrontendRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.SetHTMLTemplate(handler.NewTemplates())
	r.Use(middleware.Logging())
	r.Static("/static", "./static")

	r.GET("/", h.HandleHome)

	return r
}
