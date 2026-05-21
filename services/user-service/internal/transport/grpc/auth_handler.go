package grpc

import (
	"context"
	"net/mail"
	"unicode/utf8"

	pb "github.com/artlink52/ecommerce_backend/proto/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	minPasswordLength = 4
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

func (h *AuthHandler) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := h.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &pb.RegisterResponse{
		UserId: userID,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {}

func validateLogin(req *pb.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	if utf8.RuneCountInString(req.GetPassword()) < minPasswordLength {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	return nil
}

func validateRegister(req *pb.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	if utf8.RuneCountInString(req.GetPassword()) < minPasswordLength {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	return nil
}
