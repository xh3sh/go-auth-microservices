package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/cache"
	"github.com/xh3sh/go-auth-microservices/internal/config"
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"github.com/xh3sh/go-auth-microservices/internal/router"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	redisClient, err := cache.NewRedisClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("User Service connected to Redis successfully")

	tp := utils.InitTracer()
	repo := repository.NewRedisRepository(redisClient.GetClient())

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
		log.Println("User Service connected to RabbitMQ successfully")
	}

	r := gin.Default()
	r.Use(otelgin.Middleware("user-service"))

	userHandler := handler.NewUserHandler(repo, eventRepo)
	router.SetupUserRoutes(r, userHandler)

	addr := "0.0.0.0:" + cfg.UserServicePort
	log.Printf("User Service starting on %s", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Println("Shutting down User Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Tracer provider shutdown error: %v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("User Service stopped")
}
