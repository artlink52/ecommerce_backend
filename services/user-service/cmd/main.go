package main

import (
	"context"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/config"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/repository/pgx"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	ctx := context.Background()

	log.Info("starting service")

	gRPCServer := grpc.NewServer()

	pool, err := pgx.NewPool(ctx, cfg.Postgres)
}
