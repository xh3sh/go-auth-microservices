package main

import (
	"log"

	"github.com/xh3sh/go-auth-microservices/internal/config"
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/router"
)

func main() {
	cfg := config.Load()

	h := handler.NewHandler(nil, nil, nil, nil, nil, nil, cfg)
	r := router.NewFrontendRouter(h)

	log.Printf("Frontend starting on %s:%s", cfg.FrontendHost, cfg.FrontendPort)
	if err := r.Run(cfg.FrontendHost + ":" + cfg.FrontendPort); err != nil {
		log.Fatalf("Failed to start frontend: %v", err)
	}
}
