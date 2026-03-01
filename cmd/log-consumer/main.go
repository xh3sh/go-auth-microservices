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
	"github.com/xh3sh/go-auth-microservices/internal/mq"
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
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer eventRepo.Close()
	log.Println("Log Consumer connected to RabbitMQ successfully")

	repo := repository.NewRedisRepository(redisClient.GetClient())

	consumer := mq.NewLogConsumer(repo, eventRepo)
	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	consumer.Start(ctx)

	h := handler.NewHandler(nil, nil, nil, nil, repo, eventRepo, cfg)
	r := router.NewLogRouter(h)

	addr := cfg.LogConsumerHost + ":" + cfg.LogConsumerPort
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("Log Consumer API starting on %s", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down Log Consumer...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Log Consumer stopped")
}
