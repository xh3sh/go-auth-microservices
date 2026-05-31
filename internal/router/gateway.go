package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/xh3sh/go-auth-microservices/internal/middleware"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/gin-gonic/gin"
)

// NewGatewayRouter создает основной роутер API Gateway с функцией проксирования
func NewGatewayRouter(authAddr, userAddr, logAddr, frontendAddr string, authMiddleware *middleware.AuthMiddleware, eventRepo repository.EventRepository) *gin.Engine {
	r := gin.Default()

	r.Use(otelgin.Middleware("api-service"))

	if eventRepo != nil {
		r.Use(middleware.APILogger(eventRepo))
	}

	r.Use(func(c *gin.Context) {
		otel.GetTextMapPropagator().Inject(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		c.Next()
	})
	authURL, _ := url.Parse("http://" + authAddr)
	authProxy := httputil.NewSingleHostReverseProxy(authURL)

	userURL, _ := url.Parse("http://" + userAddr)
	userProxy := httputil.NewSingleHostReverseProxy(userURL)

	logURL, _ := url.Parse("http://" + logAddr)
	logProxy := httputil.NewSingleHostReverseProxy(logURL)

	frontendURL, _ := url.Parse("http://" + frontendAddr)
	frontendProxy := httputil.NewSingleHostReverseProxy(frontendURL)

	authGroup := r.Group("/auth")
	{
		authGroup.Any("/*any", func(c *gin.Context) {
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	api := r.Group("/api")
	api.Use(authMiddleware.UniversalAuth())
	{
		api.Any("/*any", func(c *gin.Context) {
			userProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	audit := r.Group("/audit")
	audit.Use(authMiddleware.UniversalAuth())
	{
		audit.Any("/*any", func(c *gin.Context) {
			logProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	secure := r.Group("/secure")
	secure.Use(authMiddleware.APIKeyAuth())
	{
		secure.GET("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "sensitive info via API Key"})
		})
	}

	r.Any("/static/*any", func(c *gin.Context) {
		frontendProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		frontendProxy.ServeHTTP(c.Writer, c.Request)
	})

	return r
}
