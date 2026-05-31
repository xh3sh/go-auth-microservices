package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/config"
	"github.com/xh3sh/go-auth-microservices/internal/handler"
	"github.com/xh3sh/go-auth-microservices/internal/router"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	cfg := config.Load()
	tp := utils.InitTracer()
	h := handler.NewHandler(nil, nil, nil, nil, nil, nil, cfg)
	r := router.NewFrontendRouter(h)
	r.Use(otelgin.Middleware("frontend-service"))

	log.Printf("Frontend starting on %s:%s", cfg.FrontendHost, cfg.FrontendPort)

	addr := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.FrontendPort)
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

	log.Println("Shutting down Frontend...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Tracer provider shutdown error: %v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Frontend stopped")
}
