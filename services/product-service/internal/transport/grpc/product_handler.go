package grpc

import (
	"context"
	"errors"

	pb "github.com/artlink52/ecommerce_backend/proto/gen/product"
	domainerror "github.com/artlink52/ecommerce_backend/services/product-service/internal/domain/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService interface {
	CreateProduct(
		ctx context.Context,
		name string,
		description string,
		price int64,
		stock int32,
	) (int64, error)
}

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	productService ProductService
}

func RegisterProductHandler(grpc *grpc.Server, productService ProductService) {
	pb.RegisterProductServiceServer(grpc, &ProductHandler{productService: productService})
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

	productID, err := h.productService.CreateProduct(ctx, req.GetName(), req.GetDescription(), req.GetPrice(), req.GetStock())
	if err != nil {
		if errors.Is(err, domainerror.ErrProductExists) {
			return nil, status.Error(codes.AlreadyExists, "product already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create product")

	}
	return &pb.CreateProductResponse{ProductId: productID}, nil
}

func validateRequest(req *pb.CreateProductRequest) error {
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetPrice() <= 0 {
		return status.Error(codes.InvalidArgument, "price must be positive")
	}
	if req.GetStock() < 0 {
		return status.Error(codes.InvalidArgument, "stock cannot be negative")
	}
	return nil
}
