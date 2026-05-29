package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, name, description string, price int64, stock int32) (int64, error)
}

type ProductService struct {
	log  *slog.Logger
	repo ProductRepository
}

func NewProductService(log *slog.Logger, repo ProductRepository) *ProductService {
	return &ProductService{log: log, repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, name, description string, price int64, stock int32) (int64, error) {
	const op = "product.CreateProduct"

	log := s.log.With(slog.String("op", op), slog.String("name", name))
	log.Info("creating product")

	productID, err := s.repo.CreateProduct(ctx, name, description, price, stock)
	if err != nil {
		log.Error("failed to create product: ", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("product created")

	return productID, nil
}
