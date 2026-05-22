package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/config"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/migrator"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/repository"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/repository/pgx"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/services"
	grpchandler "github.com/artlink52/ecommerce_backend/services/user-service/internal/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	ctx := context.Background()

	log.Info("starting service")

	migrator.MustRun(cfg.Postgres)

	pool, err := pgx.NewPool(ctx, cfg.Postgres)
	if err != nil {
		log.Error("failed to create connection pool")
		os.Exit(1)
	}
	userRepository := repository.NewUserRepository(pool)
	authService := services.NewAuthService(log, cfg.TokenTTL, cfg.Infra.JWTSecret, userRepository)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(timeoutInterceptor(cfg.GRPC.Timeout)))

	grpchandler.RegisterAuthHandler(grpcServer, authService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Error("failed to listen:" + err.Error())
		os.Exit(1)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("failed to serve")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Info("shutting down gracefully")
	pool.Close()
	grpcServer.GracefulStop()
	log.Info("shut down gracefully")

}

func timeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return handler(ctx, req)
	}
}
