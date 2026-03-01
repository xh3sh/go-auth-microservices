package main

import (
	"log"

	"github.com/xh3sh/go-auth-microservices/internal/cache"
	"github.com/xh3sh/go-auth-microservices/internal/config"
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"github.com/xh3sh/go-auth-microservices/internal/router"

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

	userHandler := handler.NewUserHandler(repo, eventRepo)
	router.SetupUserRoutes(r, userHandler)

	addr := cfg.UserServiceHost + ":" + cfg.UserServicePort
	log.Printf("User Service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start User Service: %v", err)
	}
}
