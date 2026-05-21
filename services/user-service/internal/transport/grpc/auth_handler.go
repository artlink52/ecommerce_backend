package grpc

import (
	"context"
	"errors"
	"net/mail"
	"unicode/utf8"

	pb "github.com/artlink52/ecommerce_backend/proto/gen/auth"
	domainerrors "github.com/artlink52/ecommerce_backend/services/user-service/internal/domain/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const minPasswordLength = 8

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
	if err := validateCredentials(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}

	userID, err := h.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &pb.RegisterResponse{UserId: userID}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validateCredentials(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}

	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domainerrors.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &pb.LoginResponse{Token: token}, nil
}

func validateCredentials(email, password string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if utf8.RuneCountInString(password) < minPasswordLength {
		return status.Errorf(codes.InvalidArgument, "password must be at least %d characters", minPasswordLength)
	}
	return nil
}
