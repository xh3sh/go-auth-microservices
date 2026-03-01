package router

import (
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewAuthRouter СЃРѕР·РґР°РµС‚ СЂРѕСѓС‚РµСЂ РґР»СЏ СЃРµСЂРІРёСЃР° Р°СѓС‚РµРЅС‚РёС„РёРєР°С†РёРё
func NewAuthRouter(h *handler.Handler, authMid *middleware.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.APILogger(h.GetEventRepo()))

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/registration", h.HandleRegister)

		jwtGroup := authGroup.Group("/jwt")
		{
			jwtGroup.POST("/login", h.HandleJWTLogin)
			jwtGroup.POST("/refresh", h.HandleJWTRefresh)
			jwtGroup.POST("/validate", authMid.JWTAuth(), h.HandleJWTValidate)
			jwtGroup.POST("/logout", h.HandleJWTLogout)
		}

		apiKeyGroup := authGroup.Group("/apikey")
		{
			apiKeyGroup.POST("/generate", h.HandleGenerateAPIKey)
			apiKeyGroup.POST("/validate", authMid.APIKeyAuth(), h.HandleValidateAPIKey)
			apiKeyGroup.POST("/revoke", h.HandleRevokeAPIKey)
		}

		oauthGroup := authGroup.Group("/oauth")
		{
			oauthGroup.GET("/authorize", h.HandleOAuthAuthorize)
			oauthGroup.GET("/google/callback", h.HandleGoogleCallback)
			oauthGroup.POST("/token", h.HandleOAuthToken)
			oauthGroup.POST("/validate", authMid.OAuthAuth(), h.HandleOAuthValidate)
			oauthGroup.POST("/logout", h.HandleOAuthLogout)
		}

		sessionGroup := authGroup.Group("/session")
		{
			sessionGroup.POST("/create", h.HandleCreateSession)
			sessionGroup.POST("/validate", authMid.SessionAuth(), h.HandleValidateSession)
			sessionGroup.POST("/revoke", h.HandleRevokeSession)
		}

		basicGroup := authGroup.Group("/basic")
		{
			basicGroup.POST("/validate", authMid.BasicAuth(), h.HandleBasicAuthValidate)
		}
	}

	return r
}
