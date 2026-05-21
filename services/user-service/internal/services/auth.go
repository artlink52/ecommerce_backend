package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type UserRepository interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error)
}

type AuthService struct {
	log            *slog.Logger
	tokenTTL       time.Duration
	userRepository UserRepository
}

func NewAuthService(
	log *slog.Logger,
	tokenTTL time.Duration,
	userRepository UserRepository,
) *AuthService {
	return &AuthService{
		log:            log,
		tokenTTL:       tokenTTL,
		userRepository: userRepository,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	email, password string,
) (int64, error) {
	const op = "auth.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registering user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.userRepository.SaveUser(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			log.Error("user already exists", logger.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user registered")
	return userID, nil
}
