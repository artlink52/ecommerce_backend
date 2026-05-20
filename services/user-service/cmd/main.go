package main

import (
	"github.com/artlink52/ecommerce_backend/pkg/logger"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
}
