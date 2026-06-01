package router

import (
	"github.com/xh3sh/go-auth-microservices/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	api := r.Group("/api/users")
	{
		api.GET("/", userHandler.GetUsers)
		api.GET("/:id", userHandler.GetUserByID)
		api.DELETE("/:id", userHandler.DeleteUser)
	}
}
