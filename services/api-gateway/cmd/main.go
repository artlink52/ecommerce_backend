package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
	"github.com/artlink52/ecommerce_backend/services/api-gateway/internal/clients/user"
	"github.com/artlink52/ecommerce_backend/services/api-gateway/internal/config"
	httphandler "github.com/artlink52/ecommerce_backend/services/api-gateway/internal/transport/http"
	"github.com/artlink52/ecommerce_backend/services/api-gateway/internal/transport/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting api-gateway")

	userClient, err := user.New(cfg.UserService.Addr)
	if err != nil {
		log.Error("failed to connect to user-service", logger.Err(err))
		os.Exit(1)
	}
	defer userClient.Close()

	authHandler := httphandler.NewAuthHandler(userClient)

	r := chi.NewRouter()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.Infra.JWTSecret))
		// защищённые роуты будут здесь
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  time.Minute,
	}

	go func() {
		log.Info("listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to serve", logger.Err(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("shutting down gracefully")
	_ = srv.Close()
	log.Info("shut down gracefully")
}
