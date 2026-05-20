package grpc

import (
	"context"

	pb "github.com/artlink52/ecommerce_backend/proto/gen/auth"
	"google.golang.org/grpc"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService AuthService
}

func RegisterAuthHandler(gRPC *grpc.Server, authService AuthService) {
	pb.RegisterAuthServiceServer(gRPC, &AuthHandler{authService: authService})
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {}
