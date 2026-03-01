package router

import (
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes РЅР°СЃС‚СЂР°РёРІР°РµС‚ РјР°СЂС€СЂСѓС‚С‹ РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏРјРё
func SetupUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	r.Use(middleware.APILogger(userHandler.GetEventRepo()))

	api := r.Group("/api/users")
	{
		api.GET("/", userHandler.GetUsers)
		api.GET("/:id", userHandler.GetUserByID)
		api.DELETE("/:id", userHandler.DeleteUser)
	}
}
