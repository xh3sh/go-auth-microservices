package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/auth"
	"github.com/xh3sh/go-auth-microservices/internal/cache"
	"github.com/xh3sh/go-auth-microservices/internal/config"
	"github.com/xh3sh/go-auth-microservices/internal/handler"
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

	log.Println("Connected to Redis successfully")

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
		log.Println("Auth Service connected to RabbitMQ successfully")
	}

	repo := repository.NewRedisRepository(redisClient.GetClient())

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration, cfg.JWTRefreshExpiration, repo)
	apiKeyService := auth.NewAPIKeyService(repo)
	oauthService := auth.NewOAuthService(repo)
	sessionService := auth.NewSessionService(repo)

	h := handler.NewHandler(
		jwtService,
		apiKeyService,
		oauthService,
		sessionService,
		repo,
		eventRepo,
		cfg,
	)

	authMid := middleware.NewAuthMiddleware(jwtService, apiKeyService, sessionService, repo, eventRepo)

	r := router.NewAuthRouter(h, authMid)

	addr := cfg.AuthServiceHost + ":" + cfg.AuthServicePort
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("Auth Service starting on %s", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Println("Shutting down Auth Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Auth Service stopped")
}
