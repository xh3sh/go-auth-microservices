package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/auth"
	"github.com/xh3sh/go-auth-microservices/internal/cache"
	"github.com/xh3sh/go-auth-microservices/internal/config"
	"github.com/xh3sh/go-auth-microservices/internal/middleware"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"github.com/xh3sh/go-auth-microservices/internal/router"
)

func main() {
	cfg := config.Load()

	redisClient, err := cache.NewRedisClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("Gateway connected to Redis successfully")

	eventRepo, err := repository.NewRabbitMQRepository(
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)
	if err != nil {
		log.Printf("Warning: Failed to connect to RabbitMQ: %v. Events will not be published.", err)
	} else {
		defer eventRepo.Close()
		log.Println("Gateway connected to RabbitMQ successfully")
	}

	repo := repository.NewRedisRepository(redisClient.GetClient())

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration, cfg.JWTRefreshExpiration, repo)
	apiKeyService := auth.NewAPIKeyService(repo)
	sessionService := auth.NewSessionService(repo)

	authMiddleware := middleware.NewAuthMiddleware(
		jwtService,
		apiKeyService,
		sessionService,
		repo,
		eventRepo,
	)

	authAddr := fmt.Sprintf("%s:%s", cfg.AuthServiceHost, cfg.AuthServicePort)
	userAddr := fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort)
	logAddr := fmt.Sprintf("%s:%s", cfg.LogConsumerHost, cfg.LogConsumerPort)
	frontendAddr := fmt.Sprintf("%s:%s", cfg.FrontendHost, cfg.FrontendPort)
	r := router.NewGatewayRouter(authAddr, userAddr, logAddr, frontendAddr, authMiddleware, eventRepo)

	addr := fmt.Sprintf("%s:%s", cfg.APIGatewayHost, cfg.APIGatewayPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("API Gateway starting on %s (Proxying /auth to %s)", addr, authAddr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Println("Shutting down API Gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("API Gateway stopped")
}
